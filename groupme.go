package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	return
	fetchMessages(gID, apiToken, "")
}

func fetchMessages(gID, apiToken, lastID string) string {
	url := "https://api.groupme.com/v3/groups/" + gID + "/messages?token=" + apiToken + "&limit=100"
	if lastID != "" {
		url += "&before_id=" + lastID
	}
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data GroupMessagesResponseWrapper
	decoder.Decode(&data)
	numMessages := len(data.Response.Messages)
	fmt.Println(numMessages)
	fmt.Println(data.Response.Messages[numMessages-1].Text)
	fmt.Println(data.Response.Messages[numMessages-1].ID)
	return data.Response.Messages[numMessages-1].ID
}
