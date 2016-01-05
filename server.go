package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var groupTmpl, userTmpl, adminTmpl *template.Template

func setUpServer(port int, dataChan chan GroupMeData) {
	adminTmpl = template.Must(template.ParseFiles("./html/admin.go.html"))
	groupTmpl = template.Must(template.ParseFiles("./html/group.go.html"))
	userTmpl = template.Must(template.ParseFiles("./html/user.go.html"))

	http.HandleFunc("/abf", func(w http.ResponseWriter, req *http.Request) {
		ShowGroupStats(w, req, dataChan)
	})
	http.HandleFunc("/abf/warboy", func(w http.ResponseWriter, req *http.Request) {
		ShowUserStats(w, req, dataChan)
	})
	http.HandleFunc("/abf/max", func(w http.ResponseWriter, req *http.Request) {
		adminTmpl.Execute(w, nil)
	})
	err := http.ListenAndServe(":8080", nil)
	fmt.Println("Server up and listening")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ShowGroupStats(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	groupTmpl.Execute(w, nil)
}

func ShowUserStats(w http.ResponseWriter, req *http.Request, dataChan chan GroupMeData) {
	user := req.URL.Query().Get("id")
	if user == "" {
		groupTmpl.Execute(w, nil)
	} else {
		userTmpl.Execute(w, user)
	}
}
