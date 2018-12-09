package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

const douyuserver = "openbarrage.douyutv.com:8601"

var conn net.Conn

func connect() {
	var err error
	conn, err = net.Dial("tcp", douyuserver)
	if err != nil {
		log.Errorf("Error connecting: %s", err)
		os.Exit(1)
	}
	log.Infof("Connecting to %s", douyuserver)
}

func shutdown() {
	log.Info("Close Connect")
	conn.Close()
}

func send(data string) {

	msglen := uint32(len(data) + 8)

	b := new(bytes.Buffer)

	binary.Write(b, binary.LittleEndian, msglen)
	binary.Write(b, binary.LittleEndian, msglen)

	sendType := int16(689)
	binary.Write(b, binary.LittleEndian, sendType)
	zeroPad := uint8(0)
	binary.Write(b, binary.LittleEndian, zeroPad)
	binary.Write(b, binary.LittleEndian, zeroPad)
	binary.Write(b, binary.LittleEndian, []byte(data))

	_, err := conn.Write(b.Bytes())

	if err != nil {
		log.Warnf("send err %s", err)
		return
	}
}

func recv() []byte {
	buf := make([]byte, 1024*10)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.Errorf("Error to read message because of %s", err)
		return buf
	}

	if reqLen < 18 {
		log.Errorf("Not Valid Package len %d", reqLen)
		return buf
	}

	return buf[12:reqLen]
}
