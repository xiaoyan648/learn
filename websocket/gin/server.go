package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func main() {
	r := gin.Default()
	r.GET("/echo", func(ctx *gin.Context) {
		echo(ctx.Writer, ctx.Request)
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		res := Logic(message)
		err = c.WriteMessage(mt, res)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func Logic(message []byte) []byte {
	// do something
	log.Printf("recv: %s", message)
	var result []byte
	result = append(result, message...)
	result = append(result, []byte("from server")...)
	return result
}
