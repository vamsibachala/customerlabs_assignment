package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestData struct {
	Ev     string `json:"ev"`
	Et     string `json:"et"`
	Id     string `json:"id"`
	Uid    string `json:"uid"`
	Mid    string `json:"mid"`
	T      string `json:"t"`
	P      string `json:"p"`
	L      string `json:"l"`
	Sc     string `json:"sc"`
	Atrk1  string `json:"atrk1"`
	Atrv1  string `json:"atrv1"`
	Atrt1  string `json:"atrt1"`
	Atrk2  string `json:"atrk2"`
	Atrv2  string `json:"atrv2"`
	Atrt2  string `json:"atrt2"`
	Uatrk1 string `json:"uatrk1"`
	Uatrv1 string `json:"uatrv1"`
	Uatrt1 string `json:"uatrt1"`
	Uatrk2 string `json:"uatrk2"`
	Uatrv2 string `json:"uatrv2"`
	Uatrt2 string `json:"uatrt2"`
	Uatrk3 string `json:"uatrk3"`
	Uatrv3 string `json:"uatrv3"`
	Uatrt3 string `json:"uatrt3"`
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ConvertedData struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppId           string               `json:"app_id"`
	UserId          string               `json:"user_id"`
	MessageId       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageUrl         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	Traits          map[string]Attribute `json:"traits"`
}

func handleRequest(requests chan<- RequestData, w http.ResponseWriter, r *http.Request) {
	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	requests <- data
}

func worker(requests <-chan RequestData) {
	for data := range requests {
		converted := ConvertedData{
			Event:           data.Ev,
			EventType:       data.Et,
			AppId:           data.Id,
			UserId:          data.Uid,
			MessageId:       data.Mid,
			PageTitle:       data.T,
			PageUrl:         data.P,
			BrowserLanguage: data.L,
			ScreenSize:      data.Sc,
			Attributes: map[string]Attribute{
				data.Atrk1: {Value: data.Atrv1, Type: data.Atrt1},
				data.Atrk2: {Value: data.Atrv2, Type: data.Atrt2},
			},
			Traits: map[string]Attribute{
				data.Uatrk1: {Value: data.Uatrv1, Type: data.Uatrt1},
				data.Uatrk2: {Value: data.Uatrv2, Type: data.Uatrt2},
				data.Uatrk3: {Value: data.Uatrv3, Type: data.Uatrt3},
			},
		}
		jsonData, err := json.Marshal(converted)
		if err != nil {
			fmt.Printf("Error converting data to JSON: %v\n", err)
			continue
		}
		resp, err := http.Post("https://webhook.site/", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error sending POST request: %v\n", err)
			continue
		}
		defer resp.Body.Close()
		fmt.Printf("Response status: %s\n", resp.Status)
	}
}

func main() {
	requests := make(chan RequestData)
	go worker(requests)
	http.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(requests, w, r)
	})
	http.ListenAndServe(":8080", nil)
}
