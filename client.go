package main

import (
	"fmt"
	"time"
)

func client(windows int, data chan []byte, ackNumber chan int) {
	nextFrame := 0
	failed := 0

	end :=0
	for {
		select {
		case message := <-data:
			if nextFrame == int(message[len(message)-1]) {
				if len(message) != 0 {
					end = 0
					serverColor.Println("Receiver:	Frame", nextFrame%windows, " Received:", string(message[0:len(message)-1]), time.Now())

						time.Sleep(time.Duration(propagation) * time.Millisecond)
						ackNumber <- nextFrame
						nextFrame++


				}

			}
		default:
			failed++
			if failed > 5 {
				ackNumber <- nextFrame - 1
				fmt.Println("Send Ack")
				failed = 0
				end++
			}
			if end > 2 {
				wg.Done()
			}
			timeoutColor.Printf("Reciever:\n\tNo Message\n")
			time.Sleep(1 * time.Second)
		}

	}

}
