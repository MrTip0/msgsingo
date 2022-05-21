package main

import (
	"container/list"
	"fmt"
)

func connectionController(update chan update, reciver chan message) {
	sendChannels := list.New()
	sendChannelsMutex := make(chan bool, 1)

	sendChannelsMutex <- true

	go func() {
		for {
			upd := <-update
			<-sendChannelsMutex
			if upd.action == Add {
				sendChannels.PushBack(upd.rx)
			} else if upd.action == Remove {
				el := sendChannels.Front()
				for el != nil && el.Value != upd.rx {
					el = el.Next()
				}
				if el != nil {
					sendChannels.Remove(el)
				}
			}
			sendChannelsMutex <- true
			if upd.action == Add {
				reciver <- message{fmt.Sprintf("%s has joined", upd.name), "server", upd.rx}
			} else if upd.action == Remove {
				reciver <- message{fmt.Sprintf("%s has left", upd.name), "server", nil}
			}
		}
	}()

	for {
		nmess := <-reciver
		<-sendChannelsMutex

		el := sendChannels.Front()

		for el != nil {
			channel := el.Value.(chan message)
			if channel != nmess.rx {
				channel <- nmess
			}
			el = el.Next()
		}

		sendChannelsMutex <- true
	}
}
