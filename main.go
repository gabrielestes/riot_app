package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// YouTube API docs: https://developers.google.com/youtube/v3/docs
// This Go app fetches YouTube channel data, such as video info.
// `ChannelResponse` and `VideoResponse` structs parse the JSON responses from the YouTube API.
// `getChannelInfo` and `getVideoInfo`, issue HTTP requests to the API for YT channel data and video info.

const (
	baseURL       = "https://www.googleapis.com/youtube/v3/"
	channelInfo   = "channels?part=snippet,statistics&id="
	videoInfo     = "search?part=snippet&type=video&order=date&maxResults=10&channelId="
	apiKeyParam   = "&key="
)

type ChannelResponse struct {
	Items []struct {
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			PublishedAt string `json:"publishedAt"`
		} `json:"snippet"`
		Statistics struct {
			SubscriberCount string `json:"subscriberCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type VideoResponse struct {
	Items []struct {
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			Title       string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

func main() {
	http.HandleFunc("/channel-info", channelInfoHandler)

	port := "8080"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func channelInfoHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel_id")
	if channelID == "" {
		http.Error(w, "channel_id is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	channelInfo, err := getChannelInfo(channelID, apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(channelInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getChannelInfo(channelID, apiKey string) (*ChannelResponse, error) {
	resp, err := http.Get(baseURL + channelInfo + channelID + apiKeyParam + apiKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var channelResponse ChannelResponse
	err = json.Unmarshal(body, &channelResponse)
	if err != nil {
		return nil, err
	}

	return &channelResponse, nil
}