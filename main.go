package main

import (
	"os"
	"strconv"
)

const TIME_GRANULARITY = 12

type GroupMeData struct {
	// posterID -> numPosts
	NumPosts map[string]int

	// posterID -> likerID -> count
	LikeMatrix map[string]map[string]int

	// posterID -> timeOfDay -> count
	TimeOfDayPostMatrix map[string][]int
}

func main() {
	groupMeData := GroupMeData{
		make(map[string]int),
		make(map[string]map[string]int),
		make(map[string][]int),
	}

	dataChan := make(chan GroupMeData, 1)
	dataChan <- groupMeData

	groupID := os.Args[1]
	apiToken := os.Args[2]
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
