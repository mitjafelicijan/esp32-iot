package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	coap "github.com/dustin/go-coap"
)

type IncomingMessage struct {
	DeviceId string  `xml:"deviceId" json:"deviceId"`
	Message  bool    `xml:"message" json:"message"`
	Value    float64 `xml:"value" json:"value"`
	Data     []struct {
		Metric    string  `xml:"metric" json:"metric"`
		Value     float64 `xml:"value" json:"value"`
		Timestamp int16   `xml:"timestamp" json:"timestamp"`
	}
}

func handleMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	var message IncomingMessage
	json.Unmarshal(m.Payload, &message)

	fmt.Println(message)

	for _, item := range message.Data {
		fmt.Println(item.Metric, "-->", item.Value)
	}

	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("Message recieved"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)
		return res
	}
	return nil
}

func main() {
	log.Println("Server starting")

	mux := coap.NewServeMux()
	mux.Handle("/message", coap.FuncHandler(handleMessage))

	log.Fatal(coap.ListenAndServe("udp", ":5683", mux))
}
