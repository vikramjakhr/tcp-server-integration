package util

import (
	"strings"
	"fmt"
	"strconv"
	"flag"
)

var Args OSArgs

type JsonPayload struct {
	State State `json:"state"`
}

type State struct {
	Reported Reported `json:"reported"`
}

type Reported struct {
	Data Payload `json:"data"`
}

type Payload struct {
	CO    float64 `json:"CO"`
	Temp  float64 `json:"temp"`
	Hum   float64 `json:"hum"`
	Pre   float64 `json:"pre"`
	O3    float64 `json:"O3"`
	SO2   float64 `json:"SO2"`
	NO2   float64 `json:"NO2"`
	PM1   float64 `json:"PM1"`
	PM2_5 float64 `json:"PM2.5"`
	PM10  float64 `json:"PM10"`
}

func ParsePayload(data string) JsonPayload {
	var dataMap map[string]float64 = make(map[string]float64)
	tabSplitData := strings.Split(data, "\t")
	for _, value := range tabSplitData {
		fields := strings.Split(value, "=")
		if fields != nil && len(fields) == 2 {
			val, _ := strconv.ParseFloat(fields[1], 64)
			dataMap[fields[0]] = val
		}
	}
	fmt.Println(dataMap)
	return payload(dataMap)
}

func payload(dataMap map[string]float64) JsonPayload {
	return JsonPayload{
		State: State{
			Reported: Reported{
				Data: Payload{
					CO:    dataMap["CO"],
					Temp:  dataMap["temp"],
					Hum:   dataMap["hum"],
					Pre:   dataMap["pre"],
					O3:    dataMap["O3"],
					SO2:   dataMap["SO2"],
					NO2:   dataMap["NO2"],
					PM1:   dataMap["PM1"],
					PM2_5: dataMap["PM2.5"],
					PM10:  dataMap["PM10"],
				},
			},
		},
	}
}

type OSArgs struct {
	Port           string
	CertFile       string
	PrivateKeyFile string
	Host           string
	ShadowUpdate   string
	ClientID       string
}

func ParseArgs() OSArgs {
	port := flag.String("port", "9292", "TCP port to listen")
	certFile := flag.String("certFile", "/opt/certificates/certificate.pem.crt", "Certificate file")
	privateKeyFile := flag.String("privateKeyFile", "/opt/certificates/private.pem.key", "Private key file")
	host := flag.String("host", "a2kbv0xdj5homp.iot.us-east-1.amazonaws.com", "AWS IoT host name")
	shadowUpdate := flag.String("shadowUpdate", "$aws/things/IoI/shadow/update", "Shadow update topic name")
	clientID := flag.String("clientID", "Client", "Shadow update topic name")
	flag.Parse()
	return OSArgs{
		Port:           *port,
		CertFile:       *certFile,
		PrivateKeyFile: *privateKeyFile,
		Host:           *host,
		ShadowUpdate:   *shadowUpdate,
		ClientID:       *clientID,
	}
}
