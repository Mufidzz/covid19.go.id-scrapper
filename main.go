package main

import (
	educationEngine "./Engine/Education"
	hoaxEngine "./Engine/Hoax"
	newsEngine "./Engine/News"
	protocolEngine "./Engine/Protocol"
	"time"
)

func main() {
	go newsEngine.Start(time.Second * 90)
	go hoaxEngine.Start(time.Second * 90)
	go educationEngine.Start(time.Second * 90)
	protocolEngine.Start(time.Second * 90)
}
