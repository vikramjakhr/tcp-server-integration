package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/vikramjakhr/tcp-server/util"
	"encoding/json"
)

const (
	port         = 8883
	shadowUpdate = "$aws/things/IoI/shadow/update"
)

var certificate tls.Certificate = loadCertificate()
var host string = "a2kbv0xdj5homp.iot.us-east-1.amazonaws.com"
var brokerURL string = fmt.Sprintf("tcps://%s:%d%s", host, port, "/mqtt")
var cid string = "ClientID"
var connOpts *MQTT.ClientOptions = clientOptions()
var mqttClient MQTT.Client = client()

func init() {
	log.Println("Executing init of mqtt")
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("[MQTT] Connected")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println(c)
		mqttClient.Disconnect(250)
		fmt.Println("[MQTT] Disconnected")
	}()
}

func loadCertificate() tls.Certificate {
	cer, err := tls.LoadX509KeyPair("/home/vikram/Downloads/IoI/8e8a17c03a-certificate.pem.crt",
		"/home/vikram/Downloads/IoI/8e8a17c03a-private.pem.key")
	check(err)
	return cer
}

func clientOptions() *MQTT.ClientOptions {
	config := &MQTT.ClientOptions{
		ClientID:             cid,
		CleanSession:         true,
		AutoReconnect:        true,
		MaxReconnectInterval: 1 * time.Second,
		KeepAlive:            30 * time.Second,
		TLSConfig:            tls.Config{Certificates: []tls.Certificate{certificate}},
	}
	config.AddBroker(brokerURL)
	return config
}

func client() MQTT.Client {
	return MQTT.NewClient(connOpts)
}

func Publish(payload util.JsonPayload) {
	bytes, _ := json.Marshal(payload);
	log.Println("Publishing payload : ", string(bytes))
	token := mqttClient.Publish(shadowUpdate, byte(0), false, string(bytes))
	token.Wait()
	log.Println("Published")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
