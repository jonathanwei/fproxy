package main

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

func proxyTCPConn(from net.Conn, toAddr string) {
	defer from.Close()

	to, err := net.Dial("tcp", toAddr)
	if err != nil {
		glog.Warningf("Tried to connect to %v, got %v.", toAddr, err)
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
		glog.Infof("Closing connection due to %v.", err)
	}
}

func proxyTCP(from string, to string) {
	var wg sync.WaitGroup

	l, err := net.Listen("tcp", from)
	if err != nil {
		glog.Warningf("Tried to listen on %v, got %v.", from, err)
		return
	}
	defer l.Close()

	glog.Infof("Listening on %v.", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Warningf("Tried to accept on %v, got %v.", from, err)
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			proxyTCPConn(conn, to)
		}()
	}

	wg.Wait()
}

func runTCPProxy(routes []*pb.TCPProxyRoute) {
	var wg sync.WaitGroup
	for _, route := range routes {
		route := route
		wg.Add(1)
		go func() {
			defer wg.Done()
			proxyTCP(route.Listen, route.Dial)
		}()
	}
	wg.Wait()
}
