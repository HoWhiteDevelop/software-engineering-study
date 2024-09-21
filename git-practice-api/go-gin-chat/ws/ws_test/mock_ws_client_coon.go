package main

import (
	"flag"
	"git-practice-api/go-gin-chat/result"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"strconv"
)

func start() {
	var addr = flag.String("addr", "localhost:8322", "http service address")

	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {

		result.Failture(result.APIcode.DialError, result.APIcode.GetMessage(result.APIcode.DialError), nil, &err)
	}
	defer c.Close()

	p := os.Args

	log.Println("Args", p)

	d := make(map[string]interface{})
	d["status"] = 1

	// string转成int64：
	uid, _ := strconv.ParseFloat(p[1], 64)

	d["data"] = map[string]interface{}{
		"uid":       uid,
		"room_id":   "1",
		"avatar_id": "4",
		"username":  "suiji" + p[1],
	}

	c.WriteJSON(d)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			result.Failture(result.APIcode.ReadError, result.APIcode.GetMessage(result.APIcode.ReadError), nil, &err)
			break
		}
		log.Printf("recv: %s", message)
	}

}

func main() {
	start()
}
