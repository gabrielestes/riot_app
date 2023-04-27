package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// YouTube API docs: https://developers.google.com/youtube/v3/docs
// This Go app fetches YouTube channel data, such as video info.
// `ChannelResponse` and `VideoResponse` structs parse the JSON responses from the YouTube API.
// `getChannelInfo` and `getVideoInfo`, issue HTTP requests to the API for YT channel data and video info.

const (
	youtubeAPIKey = "YOUR_API_KEY"
)

type Channel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

var tmpl = template.Must(template.New("channelInfo").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>YouTube Channel Information</title>
	</head>
	<body>
		<h1>Channel Information</h1>
		{{range .Items}}
			<h2>{{.Snippet.Title}}</h2>
			<p>{{.Snippet.Description}}</p>
		{{end}}
	</body>
	</html>
`))

type ChannelResponse struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt string  `json:"publishedAt"`
			ChannelID   string  `json:"channelId"`
			Title       string  `json:"title"`
			Description string  `json:"description"`
			Thumbnails  string  `json:"thumbnails"`
			Channel     Channel `json:"channel"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo PageInfo `json:"pageInfo"`
}

type VideoResponse struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelID   string `json:"channelId"`
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo PageInfo `json:"pageInfo"`
}

var videoTmpl = template.Must(template.New("videoInfo").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>YouTube Video Information</title>
	</head>
	<body>
		<h1>Video Information</h1>
		{{range .Items}}
			<h2>{{.Snippet.Title}}</h2>
			<p>{{.Snippet.Description}}</p>
			<p>Published at: {{.Snippet.PublishedAt}}</p>
		{{end}}
	</body>
	</html>
`))

func getChannelInfo(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("id")
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=snippet&id=%s&key=%s", channelID, youtubeAPIKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var channelResponse ChannelResponse
	if err := json.Unmarshal(body, &channelResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, channelResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getVideoInfo(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("id")
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoID, youtubeAPIKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var videoResponse VideoResponse
	if err := json.Unmarshal(body, &videoResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := videoTmpl.Execute(w, videoResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", getChannelInfo)
	http.HandleFunc("/video", getVideoInfo)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
