package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	client "github.com/ory/hydra-client-go"
	"hydra-login-consent-gin/hydra"
	"net/http"
	"time"
)

func Consent(c *gin.Context) {
	challenge := c.Query("consent_challenge")
	result, _, _ := hydra.Hydra.AdminApi.GetConsentRequest(context.Background()).
		ConsentChallenge(challenge).Execute()
	/*
		consent remember me
	*/
	if *result.Skip {
		form := &HandleConsentForm{
			Challenge:  challenge,
			GrantScope: result.RequestedScope,
		}
		redirectUrl, _ := acceptConsentRequest(*form)
		c.Redirect(http.StatusFound, redirectUrl)
		return
	}
	c.HTML(http.StatusOK, "consent.html", gin.H{
		"challenge":       challenge,
		"User":            result.Subject,
		"ClientName":      result.Client.ClientName,
		"RequestedScopes": result.RequestedScope,
	})
	return
}

type HandleConsentForm struct {
	Challenge   string   `form:"challenge"`
	GrantScope  []string `form:"grant_scope"`
	Remember    bool     `form:"remember"`
	RememberFor int64    `form:"remember_for"`
	IsConsent   string   `form:"is_consent"`
}

func HandleConsent(c *gin.Context) {
	var form HandleConsentForm
	_ = c.ShouldBind(&form)

	redirectUrl := ""
	if form.IsConsent != "Allow access" || len(form.GrantScope) == 0 {
		redirectUrl, _ = rejectConsentRequest(form)
	} else {
		redirectUrl, _ = acceptConsentRequest(form)
	}

	c.Redirect(http.StatusFound, redirectUrl)
	return
}

func acceptConsentRequest(form HandleConsentForm) (redirectUrl string, err error) {

	body := client.NewAcceptConsentRequest()
	body.SetGrantScope(form.GrantScope)
	if form.Remember {
		body.SetRemember(true)
		body.SetRememberFor(form.RememberFor * 86400) // second
	}
	/*
		OIDC protocol scope claims
	*/
	accessTokenField, idTokenField := fillScopeClaims(form.GrantScope)
	body.SetSession(client.ConsentRequestSession{
		AccessToken: accessTokenField,
		IdToken:     idTokenField,
	})

	result, _, err := hydra.Hydra.AdminApi.AcceptConsentRequest(context.Background()).
		ConsentChallenge(form.Challenge).AcceptConsentRequest(*body).Execute()

	return result.RedirectTo, err
}

func rejectConsentRequest(form HandleConsentForm) (redirectUrl string, err error) {

	body := client.NewRejectRequest()
	body.SetError("request_denied")

	if len(form.GrantScope) == 0 {
		body.SetErrorDebug("User not choose grant scope.")
	} else if form.IsConsent == "Deny access" {
		body.SetErrorDebug("User deny access.")
	} else {
		body.SetErrorDebug("Unknown reason.")
	}

	body.SetErrorDescription("User deny access.")
	body.SetErrorHint("User deny access.")
	body.SetStatusCode(http.StatusForbidden)

	result, _, err := hydra.Hydra.AdminApi.RejectConsentRequest(context.Background()).
		ConsentChallenge(form.Challenge).RejectRequest(*body).Execute()

	return result.GetRedirectTo(), err
}

// fillScopeClaims fill OIDC protocol scope claims
func fillScopeClaims(scope []string) (accessTokenField map[string]interface{}, idTokenField map[string]interface{}) {

	field := make(map[string]interface{})

	address := make(map[string]interface{})
	address["formatted"] = ""
	address["street_address"] = ""
	address["locality"] = ""
	address["region"] = ""
	address["postal_code"] = ""
	address["country"] = ""

	for _, i := range scope {
		switch i {
		case "email":
			field["email"] = "foo@bar.com"
			field["email_verified"] = true
		case "address":
			field["address"] = address
		case "phone":
			field["phone_number"] = "0123456789"
			field["phone_number_verified"] = true
		case "profile":
			//field["sub"] = "" // Hydra will set sub
			field["name"] = "Tom"
			field["given_name"] = ""
			field["family_name"] = ""
			field["middle_name"] = ""
			field["nickname"] = ""
			field["preferred_username"] = ""
			field["profile"] = ""
			field["picture"] = ""
			field["website"] = ""
			field["email"] = "foo@bar.com"
			field["email_verified"] = true
			field["gender"] = "male"
			field["birthdate"] = ""
			field["phone_number"] = "0123456789"
			field["phone_number_verified"] = true
			field["address"] = address
			field["zoneinfo"] = ""
			field["locale"] = ""
			field["updated_at"] = time.Now().Unix()
		}
	}

	return field, field
}
