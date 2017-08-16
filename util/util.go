package util

import (
	"strings"
	"fmt"
	"flag"
	"log"
	"strconv"
)

var Args OSArgs

func init() {
	log.Println("Executin util init")
	Args = ParseArgs()
}

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
	Measurement string `json:"Measurement"`
	NodeID      string `json:"NodeID"`
	TimeStamp   string `json:"TimeStamp"`
	CO          float64 `json:"CO"`
	Temp        float64 `json:"temp"`
	Hum         float64 `json:"hum"`
	Pre         float64 `json:"pre"`
	O3          float64 `json:"O3"`
	SO2         float64 `json:"SO2"`
	NO2         float64 `json:"NO2"`
	PM1         float64 `json:"PM1"`
	PM2_5       float64 `json:"PM2.5"`
	PM10        float64 `json:"PM10"`
	BatLevel    float64 `json:"BatLevel"`
}

func ParsePayload(data string) JsonPayload {
	var dataMap map[string]interface{} = make(map[string]interface{})
	tabSplitData := strings.Split(data, "\t")
	for _, value := range tabSplitData {
		fields := strings.Split(value, "=")
		if fields != nil && len(fields) == 2 {
			dataMap[fields[0]] = fields[1]
		}
	}
	fmt.Println(dataMap)
	return payload(dataMap)
}

func payload(dataMap map[string]interface{}) JsonPayload {
	return JsonPayload{
		State: State{
			Reported: Reported{
				Data: Payload{
					Measurement: "sensor",
					NodeID:      strings.TrimSpace(dataMap["NodeID"].(string)),
					TimeStamp:   strings.TrimSpace(dataMap["TimeStamp"].(string)),
					CO:          ToFloat64(dataMap["CO"]),
					Temp:        ToFloat64(dataMap["temp"]),
					Hum:         ToFloat64(dataMap["hum"]),
					Pre:         ToFloat64(dataMap["pre"]),
					O3:          ToFloat64(dataMap["O3"]),
					SO2:         ToFloat64(dataMap["SO2"]),
					NO2:         ToFloat64(dataMap["NO2"]),
					PM1:         ToFloat64(dataMap["PM1"]),
					PM2_5:       ToFloat64(dataMap["PM2.5"]),
					PM10:        ToFloat64(dataMap["PM10"]),
					BatLevel:    ToFloat64(dataMap["BatLevel"]),
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
	certFile := flag.String("certFile", "/opt/certificates/IoI-certificate.pem.crt", "Certificate file")
	privateKeyFile := flag.String("privateKeyFile", "/opt/certificates/IoI-private.pem.key", "Private key file")
	host := flag.String("host", "a2kbv0xdj5homp.iot.us-east-1.amazonaws.com", "AWS IoT host name")
	shadowUpdate := flag.String("shadowUpdate", "$aws/things/IoI/shadow/update", "Shadow update topic name")
	clientID := flag.String("clientID", "tcp-server-client", "client name")
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

func ToFloat64(any interface{}) float64 {
	num, _ := strconv.ParseFloat(any.(string), 64)
	return num
}
