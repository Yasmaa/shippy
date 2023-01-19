package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shippy/accountservice/common/messaging"
	"github.com/shippy/accountservice/dbclient"
	"github.com/shippy/accountservice/service"
	"github.com/streadway/amqp"
)

var appName = "accountservice"

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT") 
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {
	fmt.Printf("Starting the service  %v\n", appName)

	initializeBoltClient()
	initializeMessaging()
	handleSigterm(func() {
		service.MessagingClient.Close()
	})
	service.StartWebServer("9091")
}

func initializeMessaging() {

	conn := "amqp://" + rabbit_user + ":" +rabbit_password + "@" + rabbit_host + ":" + rabbit_port +"/"

	service.MessagingClient = &messaging.MessagingClient{}
	service.MessagingClient.ConnectToBroker(conn)
	service.MessagingClient.Subscribe("config_event_bus", "topic", appName, HandleRefreshEvent)
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}

func HandleRefreshEvent(d amqp.Delivery) {
	body := d.Body
	fmt.Println(body)

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

