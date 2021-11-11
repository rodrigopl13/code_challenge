package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jobsity-code-challenge/broker"
	"jobsity-code-challenge/controllers"
	"jobsity-code-challenge/middleware"
	"jobsity-code-challenge/repo"
	"jobsity-code-challenge/token"
	"jobsity-code-challenge/use_cases"
)

func Setup(broker *broker.Broker, db *repo.RepoDB, tokenService *token.TokenService, stooqURLString string) *gin.Engine {
	useCases := use_cases.New(db, broker, stooqURLString)
	ws := controllers.NewWebsocketServer(useCases, broker)
	go ws.Run()
	controller := controllers.New(useCases, ws, tokenService)
	router := gin.New()
	router.LoadHTMLGlob("public/*.html")
	v1 := router.Group("/v1")
	v1.Static("/assets", "./public/assets")
	v1.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	v1.GET("/signup.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})
	v1.GET("/main.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})
	v1.POST("/login", controller.Login)
	v1.POST("/signup", controller.SignUp)
	auth := v1.Group("/")
	auth.Use(middleware.AuthToken(tokenService.ValidateToken))
	auth.GET("/ws", controller.ProcessMsg)
	auth.GET("/message", controller.Message)
	return router
}
