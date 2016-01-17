package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func proxyTCPConn(from net.Conn, toAddr string) {
	defer from.Close()

	to, err := net.Dial("tcp", toAddr)
	if err != nil {
		log.Printf("Tried to connect to %v, got %v.", toAddr, err)
		return
	}
	defer to.Close()

	ch := make(chan error, 2)
	fn := func(to net.Conn, from net.Conn) {
		_, err := io.Copy(to, from)
		ch <- err
	}

	go fn(from, to)
	go fn(to, from)

	err = <-ch
	if err != nil {
		log.Printf("Closing connection due to %v.", err)
	}
}

func proxyTCP(from string, to string) {
	var wg sync.WaitGroup
	defer wg.Wait()

	l, err := net.Listen("tcp", from)
	if err != nil {
		log.Printf("Tried to listen on %v, got %v.", from, err)
		return
	}
	defer l.Close()

	log.Printf("Listening on %v.", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Tried to accept on %v, got %v.", from, err)
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			proxyTCPConn(conn, to)
		}()
	}
}

func runTCPProxy() {
	var config struct {
		Routes []struct {
			Listen string `json:"listen"`
			Dial   string `json:"dial"`
		} `json:"routes"`
	}

	configPath := filepath.Join(os.Getenv("HOME"), ".config/fproxy/tcp.json")
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	jsonBlob, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("JSON config file couldn't be read: %v", err)
	}

	err = json.Unmarshal(jsonBlob, &config)
	if err != nil {
		log.Fatalf("JSON config was malformed: %v", err)
	}

	var wg sync.WaitGroup
	for _, route := range config.Routes {
		route := route
		wg.Add(1)
		go func() {
			defer wg.Done()
			proxyTCP(route.Listen, route.Dial)
		}()
	}
	wg.Wait()
}
