package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	cer, err := tls.LoadX509KeyPair("/home/vikram/Downloads/IoI/8e8a17c03a-certificate.pem.crt",
		"/home/vikram/Downloads/IoI/8e8a17c03a-private.pem.key")
	check(err)

	cid := "ClientID"

	connOpts := &MQTT.ClientOptions{
		ClientID:             cid,
		CleanSession:         true,
		AutoReconnect:        true,
		MaxReconnectInterval: 1 * time.Second,
		KeepAlive:            30 * time.Second,
		TLSConfig:            tls.Config{Certificates: []tls.Certificate{cer}},
	}

	host := "a2kbv0xdj5homp.iot.us-east-1.amazonaws.com"
	port := 8883
	path := "/mqtt"
	payload := `{"state":{"reported":{"data":"CO=7.18,temp=26.45,hum=34.35,pre=97274.31,O3=-1.00,SO2=13.44,NO2=0.00	PM1=8.19,PM2.5=9.71,PM10=10.38","location":"41.12,-71.34"}}}`

	brokerURL := fmt.Sprintf("tcps://%s:%d%s", host, port, path)
	connOpts.AddBroker(brokerURL)

	mqttClient := MQTT.NewClient(connOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("[MQTT] Connected")

	token := mqttClient.Publish("$aws/things/IoI/shadow/update", byte(0), false, payload)
	token.Wait()
	log.Println("Published")

	quit := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		mqttClient.Disconnect(250)
		fmt.Println("[MQTT] Disconnected")

		quit <- struct{}{}
	}()
	<-quit
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
