package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/HouzuoGuo/tiedot/db"
	coap "github.com/dustin/go-coap"
)

var databasePath string = "/tmp/esp32-iot"
var database *db.DB
var metrics *db.Col
var err error

type IncomingMessage struct {
	DeviceId string  `xml:"deviceId" json:"deviceId"`
	Message  bool    `xml:"message" json:"message"`
	Value    float64 `xml:"value" json:"value"`
	Data     []struct {
		Metric    string `xml:"metric" json:"metric"`
		Value     string `xml:"value" json:"value"`
		Timestamp string `xml:"timestamp" json:"timestamp"`
	}
}

func handleMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {

	log.Println(">>> Incoming message")

	var message IncomingMessage
	json.Unmarshal(m.Payload, &message)

	fmt.Println(message)

	for _, item := range message.Data {
		payload := map[string]interface{}{
			"deviceId":  message.DeviceId,
			"metric":    item.Metric,
			"value":     item.Value,
			"timestamp": item.Timestamp,
		}

		docId, err := metrics.Insert(payload)
		if err != nil {
			log.Println(err)
		}

		log.Println("Inserted metric with id:", docId)
		log.Println("With payload:", payload)
	}

	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("OK"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)
		return res
	}
	return nil
}

func main() {

	database, err = db.OpenDB(databasePath)
	if err != nil {
		panic(err)
	}

	// try creating a colection 'Metrics'
	database.Create("Metrics")

	log.Println("Available collections:")
	for _, name := range database.AllCols() {
		log.Printf(" - %s\n", name)
	}

	metrics = database.Use("Metrics")

	// Process all documents (note that document order is undetermined)
	// metrics.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
	// 	fmt.Println("Document", id, "is", string(docContent))
	// 	return true
	// })

	log.Println("Server starting")

	mux := coap.NewServeMux()
	mux.Handle("/message", coap.FuncHandler(handleMessage))

	log.Fatal(coap.ListenAndServe("udp", ":5683", mux))
}
