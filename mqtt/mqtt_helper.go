package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
	"gitlab.intelligrape.net/tothenew/tcp-server-integration/util"
)

const (
	mqttPort = 8883
)

var certificate tls.Certificate = loadCertificate()
var brokerURL string = fmt.Sprintf("tcps://%s:%d%s", util.Args.Host, mqttPort, "/mqtt")
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
		mqttClient.Disconnect(250)
		fmt.Println("[MQTT] Disconnected")
		os.Exit(0)
	}()
}

func loadCertificate() tls.Certificate {
	log.Println("Loading certificates ", util.Args.CertFile+" and ", util.Args.PrivateKeyFile)
	cer, err := tls.LoadX509KeyPair(util.Args.CertFile, util.Args.PrivateKeyFile)
	check(err)
	return cer
}

func clientOptions() *MQTT.ClientOptions {
	log.Println("Creating MQTT config object")
	config := &MQTT.ClientOptions{
		ClientID:             util.Args.ClientID,
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
	log.Println("Creating New MQQT client")
	return MQTT.NewClient(connOpts)
}

func Publish(payload util.JsonPayload) {
	bytes, _ := json.Marshal(payload);
	payloadAsString := string(bytes)
	log.Println("Publishing payload : ", payloadAsString)
	token := mqttClient.Publish(util.Args.ShadowUpdate, byte(0), false, payloadAsString)
	token.Wait()
	log.Println("Published")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
