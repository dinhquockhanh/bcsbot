package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	username     = flag.String("username", "dqkhanh", "username that registered")
	host         = flag.String("host", "popcat.lnquy.com", "the host popcat server")
	DefaultPoint = flag.Int("diff", 555, "diff number, if long time you don't receive the point, plz try decrease the diff number, ex: -diff=100")
	MaxCnn       = 5

	printer = message.NewPrinter(language.English)
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("exited.")
		}
	}()
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime)

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})
	wss := make([]*Websocker, MaxCnn)
	point := make(chan int, MaxCnn)
	fmt.Println(`
  ____   _____  _____   ____   ____ _______ 
 |  _ \ / ____|/ ____| |  _ \ / __ \__   __|
 | |_) | |    | (___   | |_) | |  | | | |   
 |  _ <| |     \___ \  |  _ <| |  | | | |   
 | |_) | |____ ____) | | |_) | |__| | | |   
 |____/ \_____|_____/  |____/ \____/  |_|   
                                            
                                            
`)
	c := 0
	for i := 0; i < MaxCnn; i++ {
		ws := NewWebsocker(*host, *username, fmt.Sprintf("hand %d", i+1))
		if err := ws.Connect(); err != nil {
			log.Printf("%s is broken", ws.name)
			continue
		}
		c++
		wss[i] = ws
		go receiveMsg(done, point, ws)
		go submit(done, ws)
		go showPoints(point, done)
	}

	log.Printf("%s get point with %d big hands (max = %d)", *username, c, MaxCnn)

	go shutdown(signChan, wss)

	<-signChan
	log.Println("Shutting down")

	close(signChan)
	<-done
}

func receiveMsg(done chan struct{}, point chan<- int, ws *Websocker) {
	msg := &ReceiveMsg{}
	defer close(done)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				ws.Connect()
				continue
			}
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Println("read:", err)
			return

		}
		if err := json.Unmarshal(message, msg); err != nil {
			log.Println("decode msg:", err)
		}

		if msg.MessageCode == Changed {
			point <- msg.Data.Point
		}
	}
}

func shutdown(interrupt chan os.Signal, wss []*Websocker) {
	for {
		select {
		case <-interrupt:
			log.Println("stopping pop cat...")
			for _, ws := range wss {
				if err := ws.Close(); err != nil {
					log.Println("close connection:", err)
					return
				}
			}

		}
	}
}

func submit(done chan struct{}, ws *Websocker) {
	tickerPing := time.NewTicker(15 * time.Second)
	defer tickerPing.Stop()
	tickerSend := time.NewTicker(time.Second)
	defer tickerSend.Stop()

	for {
		select {
		case <-done:
			return
		case <-tickerPing.C:
			if err := ws.ping(); err != nil {
				log.Println("ping:", err)
			}
		case <-tickerSend.C:
			if err := ws.submitPoint(*DefaultPoint); err != nil {
				if err.Error() == "websocket: close sent" {
					ws.Connect()
				}
				log.Println("submit point:", err)
			}
		}

	}
}

func showPoints(point <-chan int, done chan struct{}) {
	cp := 0
	t := time.Now()
	for {
		select {
		case <-done:
			return
		case p := <-point:
			if cp != p {
				diff := p - cp
				if cp != 0 {
					fmt.Printf("\rYour Points: %s (incre %s points in %.fs)", printer.Sprintf("%d", p), printer.Sprintf("%d", diff), time.Since(t).Seconds())
				}
				cp = p
				t = time.Now()
			}
		}
	}
}
