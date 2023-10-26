package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscribers    int    `json:"subscribers"`
	ChannelName    string `json:"channelName"`
	MinutesWatched int    `json:"minutesWatched"`
	Views          int    `json:"views"`
}

func getChannelStats(k string, channelId string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := context.Background()
		yts, err := youtube.NewService(ctx, option.WithAPIKey(k))
		if err != nil {
			fmt.Println("failed to create service")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		call := yts.Channels.List([]string{"snippet, contentDetails, statistics"})
		response, err := call.Id(channelId).Do()
		if err != nil {
			fmt.Println("failed to create service")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(response.Items[0].Snippet.Title)
		yt := YoutubeStats{
			Subscribers: int(response.Items[0].Statistics.SubscriberCount),
			ChannelName: response.Items[0].Snippet.Title,
			Views:       int(response.Items[0].Statistics.ViewCount),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(yt); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
