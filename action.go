package main

import (
	"flag"
	"fmt"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

var roomid = 0

func init() {
	targetrid := flag.Int("room", 92000, "roomid")
	flag.Parse()
	if targetrid != nil {
		roomid = *targetrid
	}
}

func doLogin() {
	loginMsg := make(message)

	loginMsg["type"] = "loginreq"
	loginMsg["roomid"] = roomid

	log.Infof("join room %d", roomid)

	body := encode(loginMsg)

	send(body)
}

func doJoinGroup() {
	loginMsg := make(message)

	loginMsg["type"] = "joingroup"
	loginMsg["rid"] = roomid
	loginMsg["gid"] = -9999

	body := encode(loginMsg)

	send(body)

	log.Infof("Join Room Group %d", roomid)

	ticker := time.NewTicker(time.Second * 45)
	go func() {
		for range ticker.C {
			doPing()
		}
	}()
}

func doPing() {
	msg := make(message)

	msg["type"] = "mrkl"
	body := encode(msg)

	send(body)
}

func show(ch chan danmu, vch chan string) {
	for {
		var dm danmu
		dm = <-ch
		switch dm.Type {
		case TypeTextMsg:
			re, _ := regexp.Compile("\\[.*\\]")
			line := re.ReplaceAllString(dm.Line, "")
			//filter fure emoji
			if len(line) > 0 {
				showline := fmt.Sprintf("%s说:%s", dm.User, line)
				log.Info(showline)
				if len(vch) < 2 {
					vch <- showline
				} else {
					//log.Infof("skip voice,delay quene %d", len(vch))
				}
			} else {
				log.Infof("%s发了一些表情 %s", dm.User, dm.Line)
			}
		case TypeJoin:
			showline := fmt.Sprintf("欢迎%s进入直播间", dm.User)
			log.Info(showline)
			if len(vch) < 3 {
				vch <- showline
			} else {
				log.Infof("skip voice,delay quene %d", len(vch))
			}
		}
	}
}

func playVoice(vch chan string) {
	for {
		var line string
		line = <-vch
		trans(line)
	}
}

func handleRead(ch chan danmu) {
	for {
		inbody := recv()
		msg := decode(inbody)
		if len(inbody) < 3000 {
			//log.Infof("get msg : %v", msg)
		} else {
			if msg["type"] != nil {
				log.Infof("long message type %s", msg["type"])
			} else {
				log.Infof("uknown msg %v recvlen %d %s", msg, len(inbody), inbody)
			}
		}
		switch msg["type"] {
		case "loginres":
			log.Infof("get msg : %v", msg)
			doJoinGroup()
		case "chatmsg":
			dm := danmu{
				User: msg["nn"].(string),
				Line: msg["txt"].(string),
				Type: TypeTextMsg,
			}
			ch <- dm
		case "dgb":
			log.Infof("get present from %s", msg["nn"].(string))
		case "uenter":
			log.Infof("get new user %s enter", msg["nn"].(string))
			dm := danmu{
				User: msg["nn"].(string),
				Type: TypeJoin,
			}
			ch <- dm
		case "al":
			log.Info("zhubo level room")
		case "ab":
			log.Info("zhubo back to room")
		}
	}
}
