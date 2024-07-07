package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"math/rand"
)

type Payload struct {
	Id          uint
	AccountNo   string `json:"account_no"`
	BankName    string `json:"bank_name"`
	AccountName string `json:"account_name"`
}

func main() {
	http.HandleFunc("/payment", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "method not allowed",
				"data":    nil,
			})
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var payload Payload
		err = json.Unmarshal(body, &payload)
		if err != nil {
			log.Fatal(err)
		}
		payload.Id = uint(rand.Intn(100))

		payloadMinify, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
		}

		// send to broker
		err = NewBroker("amqp://testing:secrettesting@localhost:5672/").Send(payloadMinify)
		if err != nil {
			log.Fatal(err)
		}

		// response to client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "sended to broker..",
		})
	})

	fmt.Println("runnning on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
