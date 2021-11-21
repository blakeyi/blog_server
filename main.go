package main

import (
	"github.com/sirupsen/logrus"
	"blog_server/api"
	"blog_server/log"
	"blog_server/db"
	"net/http"
)

func main() {
	log.Init()
	db.Init()
	for k, v := range api.Router {
		http.HandleFunc(k, v)
	}
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		logrus.Fatal("ListenAndServe: ", err)
	}
}
