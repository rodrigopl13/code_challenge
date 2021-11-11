package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"jobsity-code-challenge/entities"
)

func (c Controller) Login(g *gin.Context) {
	var data entities.LoginData
	if err := g.ShouldBind(&data); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := c.usecase.Login(g.Request.Context(), data)
	if err != nil {
		log.Default().Println(fmt.Sprintf("use case Login error: %s", err.Error()))
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := c.tokenizer.GenerateToken(*user)
	g.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	g.JSON(http.StatusOK, user)
}
