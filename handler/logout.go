package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"hydra-login-consent-gin/hydra"
	"net/http"
)

func Logout(c *gin.Context) {
	challenge := c.Query("logout_challenge")

	result, _, _ := hydra.Hydra.AdminApi.GetLogoutRequest(context.Background()).
		LogoutChallenge(challenge).Execute()
	_ = result.Subject
	_ = result.Client

	c.HTML(http.StatusOK, "logout.html", gin.H{
		"challenge": challenge,
	})
	return
}

type HandleLogoutForm struct {
	Challenge string `form:"challenge"`
	IsLogout  string `form:"is_logout"`
}

func HandleLogout(c *gin.Context) {
	var form HandleLogoutForm
	_ = c.ShouldBind(&form)

	if form.IsLogout == "Yes" {
		redirectUrl, _ := acceptLogoutRequest(&form)
		c.Redirect(http.StatusFound, redirectUrl)
		return
	}

	_ = rejectLogoutRequest(&form)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "User reject logout.", "data": nil})
	return
}

func acceptLogoutRequest(form *HandleLogoutForm) (consentUrl string, err error) {

	result, _, err := hydra.Hydra.AdminApi.AcceptLogoutRequest(context.Background()).
		LogoutChallenge(form.Challenge).Execute()

	return result.GetRedirectTo(), err
}

func rejectLogoutRequest(form *HandleLogoutForm) error {

	_, err := hydra.Hydra.AdminApi.RejectLogoutRequest(context.Background()).
		LogoutChallenge(form.Challenge).Execute()

	return err
}

func LogoutSuccessful(c *gin.Context) {

	c.HTML(http.StatusOK, "error.html", gin.H{
		"info":            "Logout successful.",
		"infoDescription": "You are logout successful from OIDC.",
	})
	return
}

func HandleError(c *gin.Context) {

	c.HTML(http.StatusOK, "error.html", gin.H{
		"info":            c.Query("error"),
		"infoDescription": c.Query("error_description"),
	})
	return
}
