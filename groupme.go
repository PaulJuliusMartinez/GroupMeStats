package main

import (
	"encoding/json"
	//	"fmt"
	"net/http"
	"time"
)

type GroupMessagesResponseWrapper struct {
	Response GroupMessagesResponse  `json:"response"`
	Meta     map[string]interface{} `json:"meta"`
}

type GroupMessagesResponse struct {
	Count    int              `json:"count"`
	Messages []GroupMeMessage `json:"messages"`
}

type GroupMeMessage struct {
	Attachments []GroupMeMessageAttachment `json:"attachments"`
	AvatarURL   string                     `json:"avatar_url"`
	CreatedAt   int                        `json:"created_at"`
	FavoritedBy []string                   `json:"favorited_by"`
	GroupID     string                     `json:"group_id"`
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	SenderID    string                     `json:"sender_id"`
	SenderType  string                     `json:"sender_type"`
	SourceGUID  string                     `json:"source_guid"`
	System      bool                       `json:"system"`
	Text        string                     `json:"text"`
	UserID      string                     `json:"user_id"`
}

type GroupMeMessageAttachment struct {
	Type    string   `json:"type"`
	UserIDs []string `json:"user_ids"`
	Loci    [][]int  `json:"loci"`
	URL     string   `json:"url"`
}

func fetchOldMessages(gID, apiToken string, dataChan chan GroupMeData) {
	lastID := ""
	for {
		messages := fetchMessages(gID, apiToken, lastID)

		if len(messages) == 0 {
			break
		}

		gmd := <-dataChan
		for _, message := range messages {
			lastID = message.ID

			processMessage(gmd, message)

		}
		dataChan <- gmd

		lastID = messages[len(messages)-1].ID
	}
}

func processMessage(gmd GroupMeData, m GroupMeMessage) {
	gmd.trackPostCount(m)
	gmd.trackLikeCounts(m)
	gmd.trackPostTimeOfDay(m)
}

func fetchMessages(gID, apiToken, lastID string) []GroupMeMessage {
	url := "https://api.groupme.com/v3/groups/" + gID + "/messages?token=" + apiToken + "&limit=100"
	if lastID != "" {
		url += "&before_id=" + lastID
	}
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data GroupMessagesResponseWrapper
	decoder.Decode(&data)
	return data.Response.Messages
}

func (gmd *GroupMeData) trackPostCount(m GroupMeMessage) {
	poster := m.SenderID
	gmd.NumPosts[poster] += 1
}

func (gmd *GroupMeData) trackLikeCounts(m GroupMeMessage) {
	poster := m.SenderID
	if gmd.LikeMatrix[poster] == nil {
		gmd.LikeMatrix[poster] = make(map[string]int)
	}
	for _, liker := range m.FavoritedBy {
		gmd.LikeMatrix[poster][liker] += 1
	}
}

func (gmd *GroupMeData) trackPostTimeOfDay(m GroupMeMessage) {
	poster := m.SenderID
	posted := time.Unix(int64(m.CreatedAt), 0).UTC()
	bucket := (posted.Hour() + posted.Minute()) / (24 * 60 / TIME_GRANULARITY)
	if gmd.TimeOfDayPostMatrix[poster] == nil {
		gmd.TimeOfDayPostMatrix[poster] = make([]int, TIME_GRANULARITY)
	}
	gmd.TimeOfDayPostMatrix[poster][bucket] += 1
}
