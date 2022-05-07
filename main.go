package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hydra-login-consent-gin/handler"
	"hydra-login-consent-gin/hydra"
	"net/http"
)

func main() {

	hydra.InitHydra()

	var templatePath string
	flag.StringVar(&templatePath, "t", "./templates/*", "template path")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(gin.Recovery())
	e.LoadHTMLGlob(templatePath)

	router := e.Group("")
	{
		router.GET("/login", handler.Login)
		router.POST("/login", handler.HandleLogin)

		router.GET("/consent", handler.Consent)
		router.POST("/consent", handler.HandleConsent)

		router.GET("/logout", handler.Logout)
		router.POST("/logout", handler.HandleLogout)

		router.GET("/logout-successful", handler.LogoutSuccessful)
		router.GET("/error", handler.HandleError)
	}

	srv := &http.Server{
		Addr:    ":3000",
		Handler: e,
	}
	logrus.Fatal(srv.ListenAndServe())
}
