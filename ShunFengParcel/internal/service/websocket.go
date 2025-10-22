package service

import (
	"ShunFengParcel/inits"
	"fmt"
	"log"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/websocket"
)

type Node struct {
	Conn *websocket.Conn
	Data chan []byte
}

var ClientMap map[string]*Node = make(map[string]*Node, 0)
var Wg sync.WaitGroup

type Chats struct {
	UserId string `json:"userId,omitempty"`
	SendId string `json:"sendId,omitempty"`
}

func Chat(context http.Context) error {
	var a Chats
	id := a.UserId
	sendId := a.SendId

	var msg = &Chats{
		UserId: id,
		SendId: sendId,
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _ := upgrader.Upgrade(context.Response(), context.Request(), nil)

	node := &Node{
		Conn: conn,
		Data: make(chan []byte, 100),
	}

	ClientMap[id+"_"+sendId] = node

	Wg.Add(4)
	go sendMassage(node, msg)
	go readMassage(node, msg)
	Wg.Wait()
	return nil
}

func readMassage(node *Node, msg *Chats) {
	defer Wg.Done()
	for {
		select {
		case data := <-node.Data:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			rabbitmq := inits.NewRabbitMQSimple("websocket" +
				"message")
			rabbitmq.PublishSimple(string(data))
			fmt.Println("发送成功！")
			if err != nil {
				log.Println("发送失败", err)
				return
			}
		}
	}
}

func sendMassage(node *Node, msg *Chats) {
	defer Wg.Done()
	for {
		_, message, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println("接收失败")
			return
		}
		node, ok := ClientMap[msg.SendId+"_"+msg.UserId]

		if ok {

			rabbitmq := inits.NewRabbitMQSimple("websocket" +
				"message")
			rabbitmq.ConsumeSimple()
			node.Data <- message
		}

	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源的跨域请求，生产环境应按需调整
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(context http.Context) error {

	conn, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	return nil
}
