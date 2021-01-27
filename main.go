package main

import (
	"GBN-CN1/queue"
	"fmt"
	"github.com/fatih/color"
	"log"
	"sync"
	"time"
)

var (
	messages     = "Mehran Ali Mohammad Mehrshad Hossein Iman Aida Mahdie"
	propagation  int
	bandwidth    int
	fs           int
	tf           = 0.0
	wg           sync.WaitGroup
	prompt       = color.New(color.FgGreen)
	serverColor  = color.New(color.FgHiMagenta)
	clientColor  = color.New(color.FgHiYellow)
	ackColor     = color.New(color.FgHiGreen)
	timeoutColor = color.New(color.FgHiRed)
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

	timeout := 0

	waiting := 0
	for len(messageInByte) != 0 || waiting != 0 {
		if len(messageInByte) != 0 {
			if nextFrame < base+windows {
				// Create Frame with Size fs-1 of Data
				buffer = createFrame(messageInByte, fs)
				// Add Number
				buffer[len(buffer)-1] = byte(nextFrame)
				w.Enqueue(buffer)
				waiting++
				// Omit last value that was sent
				if len(messageInByte) > fs {
					messageInByte = messageInByte[fs-2:]
				} else {
					messageInByte = nil
				}
				//Send data on channel
				data <- buffer
				clientSendLog(windows, buffer, time.Now())
				clientSendTimeSimulation()
				nextFrame++

			}
		}
		select {
		case <-ack:
			base = <-ackNumber
			waiting--
			timeout = 0
			w.Dequeue()
			if base == nextFrame {
				timer.Stop()
			} else {
				timer.Reset(1 * time.Second)
				go func() {
					<-timer.C
					fmt.Println("Timout")
					if !w.IsEmpty(){
						sendWindows(&w, data)
					}
				}()
			}
		default:
			if timeout == 3 {
				sendWindows(&w, data)
				timeout = 0
			}
			fmt.Println("Timeout")
			timeout++
			print(timeout)
			time.Sleep(1 * time.Second)

		}

	}
	close(data)
	close(ack)
	done <- true
}


func sendWindows(w *queue.Windows, data chan []byte) {
	var buffer []byte
	fmt.Println("Sending Windows ")
	for !w.IsEmpty() {
		buffer = w.Dequeue()
		fmt.Println(string(buffer))
		data <- buffer
	}
	fmt.Println("End of Sending Windows")

}
func client(windows int, data chan []byte, ack chan bool, ackNumber chan int) {
	failed := 0
	nextFrame := 0
	for failed < 5{
		select {
		case message := <-data:
			failed = 0
			if nextFrame == int(message[len(message)-1]) {
				nextFrame++
				if len(message) != 0 {

					if len(message) > 2 {
						serverColor.Println("Receiver:	Frame", int(message[len(message)-1])%windows, " Received:", string(message[0:len(message)-1]), time.Now())
					} else {
						serverColor.Println("Receiver:	Frame", int(message[1])%windows, "Received", string(message[0]), time.Now())
					}
						fmt.Println("Ack Sent")
						ack <- true
						time.Sleep(time.Duration(propagation) * time.Millisecond)
						ackNumber <- int(message[len(message)-1]) + 1

				}
			}

		default:
			timeoutColor.Printf("Reciever:\n\tNo Message\n")
			failed++
			time.Sleep(1 * time.Second)
		}

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

func clientSendLog(windows int, buffer []byte, t time.Time) {
	_, _ = clientColor.Printf("Transmitter: \tFrame %d Sent:%s %s \n", int(buffer[len(buffer)-1])%windows, buffer[0:len(buffer)-1], t)
}
func clientSendTimeSimulation() {
	time.Sleep(time.Duration(propagation) * time.Millisecond)
	time.Sleep(time.Duration(tf) * time.Millisecond)
}

func main() {
	wg.Add(2)
	// Get Propagation Time & Bandwidth
	prompt.Println("Enter Propagation Time and Bandwidth: ")
	_, err := fmt.Scanf("%d %d\n", &propagation, &bandwidth)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	// Get Frame Size
	prompt.Println("Enter Frame Size")
	_, err = fmt.Scanf("%d", &fs)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	// Windows size
	windows := 4
	// Calculate Transmission Time
	tf = float64(fs / bandwidth)

	// Create a Channel to Transfer Data With Server and Client
	data := make(chan []byte, windows)

	// Create a Channel to Transfer Ack With Server and Client
	ack := make(chan bool, windows)

	// Create a Channel Check The End of Comm
	done := make(chan bool, 1)

	// Create a Channel to client send its next ack number
	ackNumber := make(chan int, windows)

	start := time.Now()
	// Start Communication
	go server(windows, data, ack, ackNumber, done)
	go client(windows, data, ack, ackNumber)

	// Wait Till Process End
	<-done

	fmt.Println("Total Time: ", time.Now().Sub(start))
}
