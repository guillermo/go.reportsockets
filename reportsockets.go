/*
Package reportsockets implements a websocket interface where all the clients
get the same data.

It internally works like a publish/subscribe model where each websocket
connection subscribe to a publisher.

It implements the http.Handler interface so you can use it with the standard net/http.

To use it, you first need to declare the channel or exchange.

exchange := reportsockets.New()

To tie it to a url, use the standard http.

http.Handle("/report", exchange.Handler()

To send messages to all the clients:

msg := "Hello world"
exchange.Publish(&msg)

Other thing that you can do is declare a handler for reciving client messages.

func myFunc(msg []byte, ws *websocket.Conn, exchange *reportsockets.Exchange){
  ...
  // You can responds to the specific client:
  ws.Write(myData)
  /e/ Or send a messages to other clients.
  exchagne.Publish(msg)
}

exchange.ClientMessageHandler = myFunc

If you don't define that method, all the client messages will be ignored.
*/
package reportsockets

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

type Exchange struct {
	ClientMessageHandler func(msg []byte, ws *websocket.Conn, exchange *Exchange)
	clients              []*websocket.Conn
	newClientChan        chan *websocket.Conn
	publishChan          chan []byte
	removeClientChan     chan *websocket.Conn
	stopChan             chan bool
}

func New() *Exchange {
	e := new(Exchange)
	e.clients = make([]*websocket.Conn, 0)
	e.newClientChan = make(chan *websocket.Conn)
	e.publishChan = make(chan []byte)
	e.stopChan = make(chan bool)
	e.removeClientChan = make(chan *websocket.Conn)
	go e.loop()
	return e
}
func (e *Exchange) loop() {

forLoop:
	for {
		fmt.Println("Starting for loop")
		select {
		case newClient := <-e.newClientChan:
			fmt.Println("We recive a new client")
			e.clients = append(e.clients, newClient)
		case msg := <-e.publishChan:
			fmt.Println("We recive a new message")
			for i, client := range e.clients {
				_, err := client.Write(msg)
				if err != nil {
					client.Close()
					e.clients = append(e.clients[:i], e.clients[(i+1):]...)
				}
			}
		case <-e.stopChan:
			fmt.Println("We were tell to stop")
			for _, client := range e.clients {
				client.Close()
			}
			break forLoop
		case oldWs := <-e.removeClientChan:
			fmt.Println("We remove an old client")
			for i, ws := range e.clients {
				if oldWs == ws {
					e.clients = append(e.clients[:i], e.clients[(i+1):]...)
				}
			}
		}
	}
	fmt.Println("This is the end")
}

func (e *Exchange) Handler() http.Handler {
	// ws.Close will be called when returning from this func
	handler := func(ws *websocket.Conn) {
		e.newClientChan <- ws
		for {
			data := make([]byte, 1024*10)
			_, err := ws.Read(data)
			if err != nil {
				fmt.Println("We teel to stop")
				// We can not remove the client e.removeClientChan <- ws
				break
			}
			// Do something with erthe message recived by the client
		}
	}
	return websocket.Handler(handler)
}

func (e *Exchange) Publish(msg []byte) {
	e.publishChan <- msg
}

func (e *Exchange) Stop() {
	close(e.stopChan)
}
