package main

import (
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.AddHook(filename.NewHook())
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
}

func main() {
	connect()
	defer shutdown()

	doLogin()

	ch := make(chan danmu)
	go handleRead(ch)
	vch := make(chan string, 5)
	go show(ch, vch)
	playVoice(vch)
}
