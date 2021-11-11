package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c Controller) ProcessMsg(g *gin.Context) {
	conn, err := wsupgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		log.Default().Println(fmt.Sprintf("websocket conn chat error: %v", err))
		return
	}

	client := newClient(conn, c.wsServer)
	go client.writePump()
	go client.readPump()

	c.wsServer.register <- client
}
