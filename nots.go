package main

import (
	"fmt"
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

func send_android(alert *string, badge int) {

}
