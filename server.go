package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

var groupTmpl, userTmpl, adminTmpl *template.Template

func setUpServer(port int, dataChan chan GroupMeData) {
	adminTmpl = template.Must(template.ParseFiles("./web/tmpl/admin.go.html"))
	groupTmpl = template.Must(template.ParseFiles("./web/tmpl/group.go.html"))
	userTmpl = template.Must(template.ParseFiles("./web/tmpl/user.go.html"))

	http.HandleFunc("/abf", func(w http.ResponseWriter, req *http.Request) {
		showGroupStats(w, req, dataChan)
	})
	http.HandleFunc("/abf/warboy", func(w http.ResponseWriter, req *http.Request) {
		showUserStats(w, req, dataChan)
	})
	http.HandleFunc("/abf/max", func(w http.ResponseWriter, req *http.Request) {
		adminTmpl.Execute(w, nil)
	})
	http.HandleFunc("/abf/answerback", func(w http.ResponseWriter, req *http.Request) {
		processNewMessage(w, req, dataChan)
	})
	http.HandleFunc("/abf/refreshtmpls", func(w http.ResponseWriter, req *http.Request) {
		refreshTemplates()
	})
	http.HandleFunc("/abf/flushandrefetch", func(w http.ResponseWriter, req *http.Request) {
		flushAndReloadData()
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Default routing\n")
		io.WriteString(w, req.URL.String())
	})

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func showGroupStats(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	gmd := <-dataChan
	group := gmd.GroupID
	api := gmd.APIToken
	dataChan <- gmd

	groupData := fetchGroup(group, api)

	gmd = <-dataChan
	groupDataStr, _ := json.Marshal(groupData)
	postDataStr, _ := json.Marshal(gmd)
	dataChan <- gmd
	groupTmpl.Execute(w, map[string]template.HTML{
		"GroupData": template.HTML(groupDataStr),
		"PostData":  template.HTML(postDataStr),
	})
}

func showUserStats(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	user := req.URL.Query().Get("id")
	if user == "" {
		groupTmpl.Execute(w, nil)
	} else {
		gmd := <-dataChan
		str, _ := json.Marshal(gmd)
		io.WriteString(w, string(str))

		dataChan <- gmd
	}
}

func processNewMessage(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	decoder := json.NewDecoder(req.Body)
	var message GroupMeMessage
	decoder.Decode(&message)

	gmd := <-dataChan
	processMessage(gmd, message)
	fmt.Println("Received message: ", message.Text)
	gid := gmd.GroupID
	api := gmd.APIToken
	recentMessages := fetchMessages(gid, api, "")

	var oldIndex = 0
	for _, newMessage := range recentMessages[1:] {
		// Find same message in old list
		for ; oldIndex < len(gmd.PreviousMessages); oldIndex++ {
			if newMessage.ID == gmd.PreviousMessages[oldIndex].ID {
				break
			}
		}

		if oldIndex == len(gmd.PreviousMessages) {
			break
		} else {
			// Remove old likes, add new likes
			gmd.trackLikeCounts(gmd.PreviousMessages[oldIndex], -1)
			gmd.trackLikeCounts(newMessage, 1)
		}
	}
	gmd.PreviousMessages = recentMessages
	fmt.Println(gmd.LikeMatrix)
	dataChan <- gmd
}

func refreshTemplates() {
	adminTmpl = template.Must(template.ParseFiles("./web/tmpl/admin.go.html"))
	groupTmpl = template.Must(template.ParseFiles("./web/tmpl/group.go.html"))
	userTmpl = template.Must(template.ParseFiles("./web/tmpl/user.go.html"))
}

func flushAndReloadData() {
}
