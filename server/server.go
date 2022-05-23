package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"net/http"
)

type socketReq struct {
	ChannelName string `json:"channelName"`
}

func main() {
	rdb := initRedis()

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			http.Error(w, "cannot update http", http.StatusInternalServerError)
			return
		}
		go consumeMsg(conn, rdb, w)
	})
	fmt.Println("start receive on port 9001")
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		panic(err)
	}
}

func consumeMsg(conn net.Conn, rdb *redis.Client, w http.ResponseWriter) {
	defer conn.Close()
	msg, op, err := wsutil.ReadClientData(conn)
	if err != nil {
		http.Error(w, "cannot update http", http.StatusInternalServerError)
		return
	}
	var req socketReq
	err = json.Unmarshal(msg, &req)
	if err != nil {
		http.Error(w, "cannot update http", http.StatusInternalServerError)
		return
	}

	fmt.Printf("request channel: %v\n", req.ChannelName)
	sub := rdb.Subscribe(context.Background(), req.ChannelName)
	defer sub.Close()
	for {
		m, err := sub.ReceiveMessage(context.Background())
		if err != nil {
			http.Error(w, "cannot update http", http.StatusInternalServerError)
			return
		}

		err = wsutil.WriteServerMessage(conn, op, []byte(m.Payload))
		if err != nil {
			fmt.Println("client is close")
			http.Error(w, "cannot write response message", http.StatusInternalServerError)
			break
		}
	}
}

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	pong := rdb.Ping(context.Background())
	if pong.Val() != "PONG" {
		panic("cannot connect redis")
	}
	return rdb
}
