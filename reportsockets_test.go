package reportsockets

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestPublish(t *testing.T) {

	// Create the exchange
	exchange := New()
	fmt.Println("Step 1")

	// Create the webserver
	ts := httptest.NewServer(exchange.Handler())
	defer ts.Close()
	fmt.Println("Step 2")

	// Create the webSocket Client
	url := "ws" + ts.URL[4:] + "/"
	clientWs, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Step 3")

	// Send a message
	text := "hello world"
	msg := []byte(text)
	exchange.Publish(msg)
	fmt.Println("Step 4")

	// Read a message
	var recivedMsg = make([]byte, 512)
	var n int
	n, err = clientWs.Read(recivedMsg)
	if err != nil {
		t.Error(err)
		return
	}
	clientWs.Close()
	// clientWs.Close()
	exchange.Stop()

	fmt.Println("Step 5")
	recivedMsg = recivedMsg[:n]
	if !bytes.Equal(recivedMsg, msg) {
		t.Errorf("Incorrect message: %v\nExpected: %v", recivedMsg, msg)
	}
	fmt.Println("Step 6")
}
/*

func TestClientMessageHandler(t *testing.T) {
	var (
		text              = "hello world"
		messageReciveChan = make(chan bool)
	)

	// Create the exchange
	exchange := New()

	// Define the handler funciont
	handler := func(msg []byte, ws *websocket.Conn, exchange *Exchange) {
		fmt.Println("And the handler was call!")
		if bytes.Equal(msg, []byte(text)) {
			t.Errorf("Incorrect message: %v\nExpected: %v", string(msg), text)
			return
		}
		messageReciveChan <- true // Test Helper
	}

	// Set up the handler function
	exchange.ClientMessageHandler = handler

	// Create the webserver
	ts := httptest.NewServer(exchange.Handler())
	defer ts.Close()

	// Create the webSocket Client
	url := "ws" + ts.URL[4:] + "/"
	clientWs, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		t.Error(err)
		return
	}

	// Send a message from the client
	msg := []byte(text)
	_, err = clientWs.Write(msg)
	if err != nil {
		t.Error(err)
		return
	}

	<-messageReciveChan
	exchange.Stop()
}

*/
