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
	groupTmpl.Execute(w, nil)
}

func showUserStats(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	user := req.URL.Query().Get("id")
	if user == "" {
		groupTmpl.Execute(w, nil)
	} else {
		userTmpl.Execute(w, user)
	}
}

func processNewMessage(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	decoder := json.NewDecoder(req.Body)
	var message GroupMeMessage
	decoder.Decode(&message)

	gmd := <-dataChan
	processMessage(gmd, message)
	fmt.Println("Received message: ", message.Text)
	// TODO(paulmtz): Keep track of last N messages and recalculate likes
	dataChan <- gmd
}

func refreshTemplates() {
	adminTmpl = template.Must(template.ParseFiles("./web/tmpl/admin.go.html"))
	groupTmpl = template.Must(template.ParseFiles("./web/tmpl/group.go.html"))
	userTmpl = template.Must(template.ParseFiles("./web/tmpl/user.go.html"))
}

func flushAndReloadData() {
}
