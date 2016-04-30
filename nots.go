package main

import (
	"fmt"
	"github.com/alexjlockwood/gcm"
	apns "github.com/anachronistic/apns"
)

func send_ios(token string, alertText *string, badge int) {
	payload := apns.NewPayload()
	payload.Alert = alertText
	payload.Badge = badge

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "certs/key.pem", "certs/secret.pkey")
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	fmt.Println("  Alert:", alert)
	fmt.Println("Success:", resp.Success)
	fmt.Println("  Error:", resp.Error)

}

func send_android(message string, token string) {
	// Create the message to be sent.
	data := map[string]interface{}{"message": message}
	regIDs := []string{token}
	msg := gcm.NewMessage(data, regIDs...)

	// Create a Sender to send the message.
	sender := &gcm.Sender{ApiKey: "AIzaSyDIjqQTv8AZRhDmSAOTDoLW1tEvSwmiLPg"}

	// Send the message and receive the response after at most two retries.
	test, err := sender.Send(msg, 2)
	fmt.Println(test)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
}
