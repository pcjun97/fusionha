package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

func playlistItemsList(service *youtube.Service, playlistID string, nextPageToken string) *youtube.PlaylistItemListResponse {
	part := []string{"snippet", "id"}
	call := service.PlaylistItems.List(part)

	call = call.PlaylistId(playlistID)
	call = call.MaxResults(50)
	call = call.PageToken(nextPageToken)

	res, err := call.Do()
	if err != nil {
		log.Fatalf("Error querying playlist items: %v", err)
	}

	return res
}

func playlistItemsInsert(service *youtube.Service, playlistID string, videoID string) *youtube.PlaylistItem {
	part := []string{"snippet"}

	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}

	call := service.PlaylistItems.Insert(part, playlistItem)

	item, err := call.Do()
	if err != nil {
		log.Fatalf("Error inserting video %v: %v", videoID, err)
	}

	return item
}

func playlistItemsDelete(service *youtube.Service, playlistItemID string) {
	call := service.PlaylistItems.Delete(playlistItemID)

	err := call.Do()
	if err != nil {
		log.Fatalf("Error deleting video %v: %v", playlistItemID, err)
	}
}

func main() {
	flag.Parse()

	client := getClient(youtube.YoutubeScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	sourcePlaylists := []string{"PLKfpbIXXKvWr3jcCU5H2phUjzMmOe2ffd", "PLj6NQzHFCvkGy5Xarg_gyh9QFzaoA-Q_3", "PL1BDC0510CF0F2E3D", "PLKfpbIXXKvWo-DjIFDnd40zSp6mmJ_5AC"}
	targetPlaylist := "PLhWVA6KF0qKzB182yToZ_1_q6-BKLWa_t"

	videoIDs := make(map[string]bool)
	nextPageToken := ""

	for _, playlistID := range sourcePlaylists {
		for {
			res := playlistItemsList(service, playlistID, nextPageToken)
			for _, item := range res.Items {
				videoIDs[item.Snippet.ResourceId.VideoId] = true
			}
			nextPageToken = res.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	for {
		res := playlistItemsList(service, targetPlaylist, nextPageToken)
		for _, item := range res.Items {
			delete(videoIDs, item.Snippet.ResourceId.VideoId)
		}
		nextPageToken = res.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	for videoID := range videoIDs {
		item := playlistItemsInsert(service, targetPlaylist, videoID)

		fmt.Printf("%v : %v\n", item.Snippet.ResourceId.VideoId, item.Snippet.Title)
		time.Sleep(5 * time.Second)
	}

}
