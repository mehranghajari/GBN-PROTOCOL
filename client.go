package main

import (
	"math/rand"
	"time"
)

func client(windows int, data chan []byte, ack chan bool, ackNumber chan int) {
	failed := 0
	nextFrame := 0
	for failed < 5 {
		select {
		case message := <-data:
			if len(message)== 0 {
				continue
			}
			failed = 0
			if nextFrame == int(message[len(message)-1]) {
				nextFrame++
				if len(message) != 0 {

					if len(message) > 2 {
						serverColor.Println("Receiver:	Frame", int(message[len(message)-1])%windows, " Received:", string(message[0:len(message)-1]), time.Now())
					} else {
						serverColor.Println("Receiver:	Frame", int(message[1])%windows, "Received", string(message[0]), time.Now())
					}
					p := rand.Float64()
					if p > 0.8 {
						ack <- true
						time.Sleep(time.Duration(propagation) * time.Millisecond)
						ackNumber <- int(message[len(message)-1]) +1
					} else {
						println("Bad ACK")
						ack <- true
						time.Sleep(time.Duration(propagation) * time.Millisecond)
						ackNumber <- int(message[len(message)-1]) - 1
					}

				}
			} else {
				ack <- true
				ackNumber <- nextFrame - 1
			}

		default:
			timeoutColor.Printf("Reciever:\n\tNo Message\n")
			failed++
			time.Sleep(1 * time.Second)
		}

	}

}