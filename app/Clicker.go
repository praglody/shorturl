package app

import (
	"log"
	"time"
)

var addClick = make(chan string)

const (
	_ = iota
	ASave
	AShutDown
)

func init() {
	go Clicker()
}

//异步统计点击数
func Clicker() {
	go func() {
		var clicks = make(map[string]int, 100)
		for {
			c := <-addClick
			if c == string(ASave) {
				go UrlCode{}.SaveClicks(clicks)
				clicks = make(map[string]int, 100)
			} else if c == string(AShutDown) {
				go UrlCode{}.SaveClicks(clicks)
			} else {
				clicks[c]++
				if len(clicks) > 1000 {
					go UrlCode{}.SaveClicks(clicks)
					clicks = make(map[string]int, 100)
				}
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			<-ticker.C
			addClick <- string(ASave)
		}
	}()
}

func StopTheWorld() {
	addClick <- string(AShutDown)
	log.Println("The program is going to shutdown,save clicks,waiting for 5s")
	time.Sleep(5 * time.Second)
}
