package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

func main() {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cert, err := certificate.FromPemFile("./config/pushcert.pem", "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = os.Getenv("DEVICE_TOKEN")
	notification.Topic = os.Getenv("TOPIC")
	notification.Payload = []byte(`{"aps":{"alert":"test","sound":"default"}}`) // See Payload section below

	client := apns2.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
