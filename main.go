package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	users        arrayFlags
	host         = flag.String("host", "popcat.lnquy.com", "the host popcat server")
	defaultPoint = flag.Int("diff", 398, "the diff number, if long time you don't receive the point, plz try decrease the diff number, ex: -diff=100")
	maxCnn       = flag.Int("max", 4, "the max connections, if max = 5, you don't have connection for browser...")
	rps          = flag.Int("rps", 1, "the request per second")
	token        = flag.String("token", "", "the team's password")

	printer = message.NewPrinter(language.English)
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("exited.")
		}
	}()

	flag.Var(&users, "u", "list users name, ex: -u=user1 -u=user2")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime)

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	wss := make([]*Websocker, *maxCnn)
	done := make(chan struct{})
	fmt.Println(`
  ____   _____  _____   ____   ____ _______ 
 |  _ \ / ____|/ ____| |  _ \ / __ \__   __|
 | |_) | |    | (___   | |_) | |  | | | |   
 |  _ <| |     \___ \  |  _ <| |  | | | |   
 | |_) | |____ ____) | | |_) | |__| | | |   
 |____/ \_____|_____/  |____/ \____/  |_|   
                                            
                                            
`)

	for _, userId := range users {
		name := userId
		go func() {
			for i := 0; i < *maxCnn; i++ {
				n := i + 1
				go Pop(name, n, &wss, done)
			}
		}()
	}

	go shutdown(signChan, &wss)

	<-signChan
	log.Println("Shutting down")
	// TODO: handle close connection
}
func Pop(username string, hand int, wss *[]*Websocker, done chan struct{}) {
	point := make(chan int, *maxCnn)
	ws := &Websocker{}
	for {
		ws = NewWebsocker(*host, username, fmt.Sprintf("hand %d", hand), *token)
		log.Printf("waiting hand %d of %s login", hand, username)
		if err := ws.Connect(); err == nil {
			*wss = append(*wss, ws)

			break

		}
	}
	log.Printf("hand %d of %s ready to popcat", hand, username)
	go receiveMsg(done, point, ws)
	go submit(done, ws)
	go showPoints(point, done, username)
}
func receiveMsg(done chan struct{}, point chan<- int, ws *Websocker) {
	msg := &ReceiveMsg{}
	defer close(done)

	for {
		_, rawMsg, err := ws.ReadMessage()
		if err != nil {
			for {
				if err := ws.Connect(); err == nil {
					break
				}
				time.Sleep(300 * time.Millisecond)
			}
			continue
		}
		if err := json.Unmarshal(rawMsg, msg); err != nil {
			log.Println("decode msg:", err)
		}

		if msg.MessageCode == Changed {
			point <- msg.Data.Point
		}
	}
}

func shutdown(interrupt chan os.Signal, wss *[]*Websocker) {
	for {
		select {
		case <-interrupt:
			log.Println("stopping pop cat...")
			for _, ws := range *wss {
				if ws == nil {
					continue
				}
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

	tickerSend := time.NewTicker(time.Duration(*rps) * time.Second)
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
			if *defaultPoint == 0 {
				continue
			}
			if err := ws.submitPoint(*defaultPoint); err != nil {
				if err.Error() == "websocket: close sent" {
					ws.Connect()
				}
				//log.Println("submit point:", err)
			}
		}
	}
}

func showPoints(point <-chan int, done chan struct{}, username string) {
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
					fmt.Printf(
						"\r%s Points: %s (incre %s points in %.fs)",
						username,
						printer.Sprintf("%d", p),
						printer.Sprintf("%d", diff),
						time.Since(t).Seconds(),
					)
				}
				cp = p
				t = time.Now()
			}
		}
	}
}
