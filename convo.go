package main

import "time"

type Convo struct {
	name     string
	messages *[]Message
}

type Message struct {
	sender string
	text   string
	time   time.Time
}

func newConvo(name string) *Convo {
	return &Convo{name: name, messages: &[]Message{}}
}

func addMessage(convo *Convo, message Message) *Convo {
	*convo.messages = append(*convo.messages, message)
	return convo
}
