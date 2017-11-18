package main

import (
	"./hand"
	"github.com/gin-gonic/gin"
	"github.com/zhouhui8915/go-socket.io-client"
	"testing"
	"time"
)

func Server() {
	r := gin.Default()
	r.GET("/socket.io/", gin.WrapH(hand.Server()))
	go r.Run("localhost: 1488")
}

func NewClient() *socketio_client.Client {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	uri := "http://localhost:1488"
	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		println("kaka")
		return nil
	}
	return client
}

func TestSomething(t *testing.T) {
	Server()
	time.Sleep(1 * time.Second)
	c := NewClient()
	c.On("error", func() {
		t.Error("kaka")
	})
}
