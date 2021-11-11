package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const maxLimit = 50

func (c Controller) Message(g *gin.Context) {
	limitString := g.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		log.Default().Println(fmt.Sprintf("limit query param wrong type %v", err))
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 || limit > maxLimit {
		s := fmt.Sprintf("limit query param %d is invalid, out of range", limit)
		log.Default().Println(s)
		g.JSON(http.StatusBadRequest, gin.H{"error": s})
		return
	}
	msgs, err := c.usecase.GetMessages(g.Request.Context(), limit)
	if err != nil {
		log.Default().Println(fmt.Sprintf("error in use case: %v", err))
		g.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong in server"})
		return
	}
	g.JSON(http.StatusOK, msgs)
}
