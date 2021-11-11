package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"jobsity-code-challenge/entities"
)

func (c Controller) SignUp(g *gin.Context) {
	var data entities.User
	if err := g.ShouldBind(&data); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.usecase.SignUp(g.Request.Context(), data); err != nil {
		log.Default().Println(fmt.Sprintf("use case signup error: %v", err))
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.HTML(http.StatusOK, "index.html", nil)
}
