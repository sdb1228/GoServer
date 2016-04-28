package main

import (
	"bytes"
	"fmt"
	apns "github.com/anachronistic/apns"
	"io"
	"os"
)

var apnKey = file2str("certs/apn.key")
var apnPem = file2str("certs/apn.pem")

func file2str(filename string) string {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open(filename) // Error handling elided for brevity.
	io.Copy(buf, f)           // Error handling elided for brevity.
	f.Close()

	s := string(buf.Bytes())

	return s
}

func send_ios(token string, alertText *string, badge int) {
	payload := apns.NewPayload()
	payload.Alert = alertText
	payload.Badge = 42

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", apnKey, apnPem)
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	fmt.Println("  Alert:", alert)
	fmt.Println("Success:", resp.Success)
	fmt.Println("  Error:", resp.Error)
}

func send_android(alert *string, badge int) {

}
