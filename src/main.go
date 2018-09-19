package main

import (
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

func main() {

	cert, err := certificate.FromPemFile("./config/pushcert.pem", "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	notification.Topic = "com.sideshow.Apns2"
	notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

	client := apns2.NewClient(cert).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
