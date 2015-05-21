package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

func RunLockagent() {
	var token string
	deadChan := make(chan error)
	r := mux.NewRouter()

	r.HandleFunc("/password", func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Token") == token {
			http.Error(res, os.Getenv("PASSWD"), 200)
		}
	}).Methods("GET")

	r.HandleFunc("/kill", func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Token") == token {
			deadChan <- fmt.Errorf("Quit")
		}
	}).Methods("POST")

	http.Handle("/", r)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Error("Error listening: %s", err)
	}
	defer listener.Close()

	fmt.Println(listener.Addr().String())
	port := strings.Split(listener.Addr().String(), ":")[1]

	token = uuid.NewV4().String()
	ioutil.WriteFile(os.Getenv("DBFILE"), []byte(fmt.Sprintf("%s::%s", token, port)), 0600)

	go func() {
		deadChan <- http.Serve(listener, nil)
	}()

	t := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-deadChan:
			os.Remove(os.Getenv("DBFILE"))
			return
		case <-t.C:
			token = uuid.NewV4().String()
			ioutil.WriteFile(os.Getenv("DBFILE"), []byte(fmt.Sprintf("%s::%s", token, port)), 0600)
		}
	}
}
