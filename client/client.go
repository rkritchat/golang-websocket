package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"time"
)

type socketReq struct {
	ChannelName string `json:"channelName"`
}

func main() {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://127.0.0.1:9001/socket")
	if err != nil {
		panic(err)
	}

	b, _ := json.Marshal(&socketReq{
		ChannelName: "rkritchat",
	})
	err = wsutil.WriteClientMessage(conn, ws.OpText, b)
	if err != nil {
		panic(err)
	}

	for {
		msg, _, err := wsutil.ReadServerData(conn)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(msg))

		time.Sleep(2 * time.Second)
	}
}
