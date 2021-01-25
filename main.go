package main

import (
	"GBN-CN1/queue"
	"fmt"
	"log"
	"time"
)

var (
	messages    = "The world! my name is mehran!? How are you "
	flagNumber  = 0
	propagation int
	bandwidth   int
	fs          int
)

func server(windows int, data chan []byte, ack chan bool, ackNumber chan int, done chan bool) {
	// Declare a buffer
	buffer := make([]byte, fs)

	// Convert Message To Byte
	messageInByte := []byte(messages)

	// Define a Windows That is like a Queue
	var w queue.Windows
	w.New(windows)

	// Declare a base var and next frame
	base := 0
	nextFrame := 0

	// Timer
	timer := time.NewTimer(4 * time.Second)
	timer.Stop()

	for len(messageInByte) != 0 {

		if nextFrame < base+windows {
			// Create Frame with Size fs-1 of Data
			buffer = createFrame(messageInByte, fs-1)
			// Add Number
			buffer[len(buffer)-1] = byte(nextFrame % 4)
			w.Enqueue(buffer)
			// Omit last value that was sent
			if len(messageInByte) > fs {
				messageInByte = messageInByte[fs-2:]
			} else {
				messageInByte = nil
			}
			//Send data on channel
			data <- buffer
			time.Sleep(2 * time.Second)

			nextFrame++
		}
		select {
		case <-ack:
			base = <-ackNumber + 1
			if base == nextFrame {
				w.Dequeue()
				timer.Stop()
			} else {
				w.Dequeue()
				timer.Reset(4 * time.Second)
				go func() {
					<-timer.C
					sendWindows(&w, data)
				}()
			}
		default:
			fmt.Println("No ack")
			time.Sleep(2 * time.Second)

		}

	}
	close(data)
	close(ack)
	done <- true
}

func client(windows int, data chan []byte, ack chan bool, ackNumber chan int) {

	for {
		select {
		case message := <-data:
			if len(message) != 0 {
				fmt.Println("Receiver: ")
				if len(message) > 2 {
					fmt.Println("	Frame", int(message[len(message)-1]), " Received:", string(message[0:len(message)-1]), time.Now())
				} else {
					fmt.Println("	Frame", int(message[1]), "Received", string(message[0]), time.Now())
				}
				time.Sleep(time.Duration(propagation) * time.Microsecond)
				ack <- true
				ackNumber <- int(message[len(message)-1]) + 1

			}

		default:
			fmt.Println("No Message")
			time.Sleep(2 * time.Second)
		}

	}

}

func sendWindows(w *queue.Windows, data chan []byte) {
	var buffer []byte
	for !w.IsEmpty() {
		buffer = w.Dequeue()
		data <- buffer
	}

}

func createFrame(message []byte, fs int) []byte {
	buffer := make([]byte, fs)
	if len(message) < fs {
		for i, b := range message {
			buffer[i] = b
		}
	} else {
		for i, b := range message[0 : fs-1] {
			buffer[i] = b
		}
	}
	return buffer
}
func main() {

	// Get Propagation Time & Bandwidth
	fmt.Println("Enter Propagation Time and Bandwidth: ")
	_, err := fmt.Scanf("%d %d", &propagation, &bandwidth)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	// Get Frame Size
	fmt.Println("Enter Frame Size")
	_, err = fmt.Scanf("%d", &fs)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	// Calculate Transmission Time
	//tf = float64(fs / bandwidth)

	// Create a Channel to Transfer Data With Server and Client
	data := make(chan []byte, 4)

	// Create a Channel to Transfer Ack With Server and Client
	ack := make(chan bool, 1)

	// Create a Channel Check The End of Comm
	done := make(chan bool, 1)

	// Create a Channel to client send its next ack number
	ackNumber := make(chan int, 4)

	// Windows size
	windows := 4
	// Start Communication
	go server(windows, data, ack, ackNumber, done)
	go client(windows, data, ack, ackNumber)

	// Wait Till Process End
	<-done
}
