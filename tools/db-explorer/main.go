package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/olekukonko/tablewriter"
)

var databasePath string = "/tmp/esp32-iot"
var database *db.DB
var metrics *db.Col
var err error

type MetricRecord struct {
	DeviceId   string `json:"deviceId"`
	Metric     string `json:"metric"`
	Value      string `json:"value"`
	Timestamp  string `json:"timestamp"`
	InsertedAt string `json:"insertedAt"`
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"DeviceId", "Metric", "Value", "Timestamp", "Inserted At"})

	metrics.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		var message MetricRecord
		json.Unmarshal(docContent, &message)

		ts, _ := strconv.ParseInt(message.InsertedAt, 10, 64)
		item := []string{
			message.DeviceId,
			message.Metric,
			message.Value,
			message.Timestamp,
			time.Unix(0, ts).String(),
		}
		table.Append(item)

		return true
	})

	table.Render()
}
