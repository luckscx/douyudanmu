package main

import (
	"bytes"
	"fmt"
	"strings"
)

type message map[string]interface{}

func encode(msg message) string {

	var tokens []string
	for k, v := range msg {
		tokens = append(tokens, fmt.Sprintf("%s@=%v/", k, v))
	}

	tokens = append(tokens, string('\u0000'))
	str := strings.Join(tokens, "")

	//log.Infof("encode msg %v", str)

	return str
}

func decode(inmsg []byte) message {
	msg := make(message)

	for _, v := range bytes.Split(inmsg, []byte("/")) {
		tokens := bytes.Split(v, []byte("@="))
		if len(tokens) == 2 {
			key := string(tokens[0])
			val := string(tokens[1])
			val = strings.Replace(val, "@A", "@", -1)
			val = strings.Replace(val, "@S", "/", -1)
			msg[key] = val
		}
	}

	return msg
}
