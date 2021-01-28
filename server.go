package main

import (
	"time"
)

// server receives message from client and sends frame's ack back.
func server(windows int, data chan []byte, ackNumber chan int) {
	nextFrame := 0
	failed := 0

	end := 0

	// A while statement to listen on link ---> non-blocking
	for {
		// Implement non-blocking mode on channel.
		select {
		// Received Data
		case message := <-data:
			if nextFrame == int(message[len(message)-1]) {
				if len(message) != 0 {
					end = 0
					_, _ = serverColor.Println("Receiver:	Frame", nextFrame%windows, " Received:", string(message[0:len(message)-1]), time.Now())
					time.Sleep(time.Duration(propagation) * time.Millisecond)
					ackNumber <- nextFrame
					nextFrame++

				}

			}
		// Doesn't received anything
		default:
			// wait for 3 time to receive data
			failed++
			// if after 6 time nothing happened, means that the link was terminated...!
			// end * failed = max number of failure
			if end > 1{
				wg.Done()
			}
			if failed > 3 {
				// make client to send last frame again.
				ackNumber <- nextFrame - 1
				_, _ = timeoutColor.Printf("Server:\t Send Ack")
				failed = 0
				end++
			}
			_, _ = timeoutColor.Printf("Reciever:\n\tNo Message\n")
			time.Sleep(1 * time.Second)
		}

	}

}
