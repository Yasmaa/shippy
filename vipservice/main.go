package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"

	"github.com/shippy/vipservice/common/messaging"
	"github.com/shippy/vipservice/service"
	"github.com/streadway/amqp"
)

var appName = "vipservice"

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

var messagingClient messaging.IMessagingClient

func main() {
	fmt.Println("Starting " + appName + "...")

	initializeMessaging()

	// Makes sure connection is closed when service exits.
	handleSigterm(func() {
		if messagingClient != nil {
			messagingClient.Close()
		}
	})
	service.StartWebServer("9092")
}

func onMessage(delivery amqp.Delivery) {
	fmt.Printf("Got a message: %v\n", string(delivery.Body))
	vipNotification := Reply{Msg: "recieved"}
	data, _ := json.Marshal(vipNotification)
	err := messagingClient.PublishOnQueue(data, "vip_queue_reply")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initializeMessaging() {

	messagingClient = &messaging.MessagingClient{}
	conn := "amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/"

	messagingClient.ConnectToBroker(conn)

	// Call the subscribe method with queue name and callback function
	err := messagingClient.SubscribeToQueue("vip_queue", appName, onMessage)
	failOnError(err, "Could not start subscribe to vip_queue")

	err = messagingClient.Subscribe("config_event_bus", "topic", appName, HandleRefreshEvent)
	failOnError(err, "Could not start subscribe to "+"config_event_bus"+" topic")
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func HandleRefreshEvent(d amqp.Delivery) {
	vipNotification := Reply{Msg: "recieved"}
	data, _ := json.Marshal(vipNotification)
	err := messagingClient.PublishOnQueue(data, "vip_queue_reply")
	if err != nil {
		fmt.Println(err.Error())
	}

}

type UpdateToken struct {
	Type               string `json:"type"`
	Timestamp          int    `json:"timestamp"`
	OriginService      string `json:"originService"`
	DestinationService string `json:"destinationService"`
	Id                 string `json:"id"`
}
type Reply struct {
	Msg string `json:"msg"`
}
