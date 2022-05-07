package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	client "github.com/ory/hydra-client-go"
	"hydra-login-consent-gin/hydra"
	"net/http"
)

func Login(c *gin.Context) {
	challenge := c.Query("login_challenge")
	result, _, _ := hydra.Hydra.AdminApi.GetLoginRequest(context.Background()).
		LoginChallenge(challenge).Execute()
	/*
		login remember me
	*/
	if result.Skip {
		form := &HandleLoginForm{
			Challenge: challenge,
			Email:     result.Subject,
		}
		consentUrl, _ := acceptLoginRequest(form)
		c.Redirect(http.StatusFound, consentUrl)
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"challenge": challenge,
	})
	return
}

// HandleLoginForm Hydra subject use Email
type HandleLoginForm struct {
	Challenge   string `form:"challenge"`
	Email       string `form:"email"`
	Password    string `form:"password"`
	Remember    bool   `form:"remember"`
	RememberFor int64  `form:"remember_for"`
}

func HandleLogin(c *gin.Context) {

	var form HandleLoginForm
	_ = c.ShouldBind(&form)

	redirectUrl := ""
	if form.Email != "foo@bar.com" || form.Password != "foobar" {
		redirectUrl, _ = rejectLoginRequest(&form)
	} else {
		consentUrl, _ := acceptLoginRequest(&form)
		redirectUrl = consentUrl
	}

	c.Redirect(http.StatusFound, redirectUrl)
	return
}

// acceptLoginRequest return consentUrl
func acceptLoginRequest(form *HandleLoginForm) (consentUrl string, err error) {

	body := client.NewAcceptLoginRequest(form.Email)
	if form.Remember {
		body.SetRemember(true)
		body.SetRememberFor(form.RememberFor * 86400) // second
	}

	result, _, err := hydra.Hydra.AdminApi.AcceptLoginRequest(context.Background()).
		LoginChallenge(form.Challenge).AcceptLoginRequest(*body).Execute()

	return result.GetRedirectTo(), err
}

// rejectLoginRequest return redirectUrl
func rejectLoginRequest(form *HandleLoginForm) (redirectUrl string, err error) {

	body := client.NewRejectRequest()
	body.SetError("request_denied")
	body.SetErrorDebug("Email or password is invalid.")
	body.SetErrorDescription("Email or password is invalid. Please check your email or password.")
	body.SetErrorHint("Please check your email or password.")
	body.SetStatusCode(http.StatusUnauthorized)

	result, _, err := hydra.Hydra.AdminApi.RejectLoginRequest(context.Background()).
		LoginChallenge(form.Challenge).RejectRequest(*body).Execute()

	return result.GetRedirectTo(), err
}
