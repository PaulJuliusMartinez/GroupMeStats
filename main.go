package main

import (
	"os"
	"strconv"
)

const TIME_GRANULARITY = 12

type GroupMeData struct {
	GroupID  string `json:"-"`
	APIToken string `json:"-"`

	// posterID -> numPosts
	NumPosts map[string]int `json:"numPosts"`

	// posterID -> likerID -> count
	LikeMatrix map[string]map[string]int `json:"likes"`

	// posterID -> timeOfDay -> count
	TimeOfDayPostMatrix map[string][]int `json:"postingTimes"`

	// previous posts for tracking likes
	PreviousMessages []GroupMeMessage `json:"-"`
}

func main() {
	groupID := os.Args[1]
	apiToken := os.Args[2]

	groupMeData := GroupMeData{
		groupID,
		apiToken,
		make(map[string]int),
		make(map[string]map[string]int),
		make(map[string][]int),
		nil,
	}

	dataChan := make(chan GroupMeData, 1)
	dataChan <- groupMeData
	port := 8080
	if len(os.Args) > 3 {
		parsedPort, err := strconv.Atoi(os.Args[3])
		if err != nil {
			port = 8080
		} else {
			port = parsedPort
		}
	}

	go fetchOldMessages(groupID, apiToken, dataChan)
	setUpServer(port, dataChan)
}
