package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	conn, _, err := websocket.DefaultDialer.Dial("ws://:6000", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < 5; i++ {
		err := conn.WriteJSON("ping")
		if err != nil {
			log.Fatal(err)
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))

		time.Sleep(time.Second)
	}
}
