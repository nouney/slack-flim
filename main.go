package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SLACK_TOKEN = ""
)

type SlackResponse struct {
	ResponseType string                     `json:"response_type"`
	Text         string                     `json:"text"`
	Attachments  []SlackResponseAttachments `json:"attachments"`
}

type SlackResponseAttachments struct {
	ImageURL string `json:"image_url"`
	ThumbURL string `json:"thumb_url"`
	Text     string `json:"text"`
}

func BaseUrl(r *http.Request) string {
	var proto string
	if r.TLS != nil {
		proto = `https`
	} else {
		proto = `http`
	}

	host := r.Host
	return strings.Join([]string{proto, `://`, host}, ``)
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(j)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if len(r.PostForm["token"]) > 0 && r.PostForm["token"][0] != SLACK_TOKEN {
		return
	}

	var idx int
	if len(r.PostForm["text"]) > 0 {
		selected := []int{}
		for i, q := range DB {
			if strings.Contains(q.Text, r.PostForm["text"][0]) {
				selected = append(selected, i)
			}
		}
		if len(selected) > 0 {
			idx = selected[rand.Intn(len(selected))]
		} else {
			idx = rand.Intn(len(DB))
		}
	} else {
		idx = rand.Intn(len(DB))
	}

	imgurl := BaseUrl(r) + "/imgs/" + DB[idx].Img
	sr := SlackResponse{
		ResponseType: "ephemeral",
		Attachments: []SlackResponseAttachments{
			{
				ImageURL: imgurl,
				Text:     DB[idx].Text,
			},
		},
	}
	WriteJSON(w, sr)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rand.Seed(time.Now().UnixNano())
	http.Handle("/imgs/", http.StripPrefix("/imgs/", http.FileServer(http.Dir("./imgs/"))))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
