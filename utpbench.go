package main

import (
	"fmt"
	"flag"
	"io"
	"sync"
	"time"
	"log"
	"net"
	"bytes"

	"github.com/anacrolix/utp"
)

var (
	flClientMode = flag.Bool("c", true, "client mode")
	flServerMode = flag.Bool("s", false, "server mode")
	flHost = flag.String("h", "127.0.0.1", "host")
	flPort = flag.Int("p", 6001, "port")
	flLen = flag.Int("l", 1400, "length of data")
	flThreads = flag.Int("t", 1, "threads")
	flDuration = flag.Duration("d", time.Second * 10, "duration")
)

func main() {
	log.Printf("UTP Benchmark Tool by Artem Andreenko (miolini@gmail.com)")
	flag.Parse()
	ts := time.Now()
	wg := sync.WaitGroup{}
	if *flServerMode {
		wg.Add(1)
		go server(&wg, *flHost, *flPort)
	} else {
		wg.Add(*flThreads)
		for i:=0;i<*flThreads;i++ {
			go client(&wg, *flHost, *flPort, *flLen, *flDuration)
		}
	}
	wg.Wait()
	log.Printf("time takes %.2fsec", time.Since(ts).Seconds())
}

func server(wg *sync.WaitGroup, host string, port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %s", r)
		}
		wg.Done()
	}()
	log.Printf("server listen %s:%d", host, port)
	s, err := utp.NewSocket("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
	defer s.Close()
	for {
		conn, err := s.Accept()
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go readConn(conn)
	}
}

func readConn(conn net.Conn) {
	defer conn.Close()
	defer log.Printf("client disconnected")
	log.Printf("new connection")
	buf := make([]byte, 4096)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("err: %s", err)
			}
			break
		}
	}
}

func client(wg *sync.WaitGroup, host string, port, len int, duration time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %s", r)
		}
		log.Printf("disconnected")
		wg.Done()
	}()
	log.Printf("connecting to %s:%d, len %d, duration %s", host, port, len, duration.String())
	conn, err := utp.DialTimeout(fmt.Sprintf("%s:%d", host, port), time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Printf("connected")
	buf := bytes.Repeat([]byte("H"), len)
	ts := time.Now()
	te := ts
	count := 0
	for time.Since(ts) < duration {
		since := time.Since(te)
		if since >= time.Second {
			te = time.Now()
			log.Printf("speed %.4f mbit/sec", float64(count) * 8 / since.Seconds() / 1024 / 1024)
			count = 0
		}
		n, err := conn.Write(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		count += n
	}
}