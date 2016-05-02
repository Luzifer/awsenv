package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

func runLockagent() {
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
		log.Errorf("Error listening: %s", err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Errorf("Unable to close listener: %s", err)
		}
	}()

	fmt.Println(listener.Addr().String())
	port := strings.Split(listener.Addr().String(), ":")[1]

	token, err = writeTokenFile(port)
	if err != nil {
		os.Exit(1)
	}

	go func() {
		deadChan <- http.Serve(listener, nil)
	}()

	t := time.NewTicker(10 * time.Second)

	if timeoutRaw := os.Getenv("TIMEOUT"); timeoutRaw != "" && timeoutRaw != "0" {
		timeout, err := time.ParseDuration(timeoutRaw)
		if err == nil {
			go func(timeout time.Duration) {
				<-time.After(timeout)
				deadChan <- errors.New("Agent timeout")
			}(timeout)
		}
	}

	for {
		select {
		case <-deadChan:
			err := os.Remove(os.Getenv("DBFILE"))
			if err != nil {
				log.Errorf("Could not delete token file: %s", err)
			}
			return
		case <-t.C:
			token, err = writeTokenFile(port)
			if err != nil {
				os.Exit(1)
			}
		}
	}
}

func writeTokenFile(port string) (string, error) {
	t := uuid.NewV4().String()
	if err := os.MkdirAll(path.Dir(os.Getenv("DBFILE")), 0755); err != nil {
		return "", err
	}
	err := ioutil.WriteFile(os.Getenv("DBFILE"), []byte(fmt.Sprintf("%s::%s", t, port)), 0600)
	if err != nil {
		log.Errorf("Unable to save token file: %s", err)
		return "", err
	}
	return t, nil
}
