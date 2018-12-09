package main

type danmu struct {
	User string
	Line string
	Type int
}

//danmu type enum
const (
	TypeTextMsg = iota
	TypeJoin
)
