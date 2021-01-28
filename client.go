package main

import (
	"bytes"
	"encoding/gob"
	"github.com/emirpasic/gods/lists/arraylist"
	"time"
)

var (
	numberOfFrame int
)

// client starts sending messages to server.
func client(windows int, data chan []byte, ackNumber chan int, done chan bool) {

	// Declare a base var and next frame
	base := 0
	nextFrame := 0

	// Calculate time which need for timeout by windows size

	// Timer
	timer := time.NewTimer(4 * time.Second)
	timer.Stop()

	// Convert message to frames of Byte
	frameList := createArrayOfFrames(fs)
	numberOfFrame = frameList.Size()
	var buf bytes.Buffer
	count := 0
	for count < frameList.Size() {
		// Send Until the windows is not full
		if nextFrame < base+windows && timeout != true {
			// Send frame on link
			sendFrame(frameList, &buf, nextFrame, data)
			// Show sent frame logs
			clientSendLog(windows, buf.Bytes(), time.Now())
			// Simulate sending time
			clientSendTimeSimulation()
			// Reset the buffer to don't effect on next Frame
			buf.Reset()

			// End of loop
			nextFrame++
			count++

		}
		// Check ack status
		checkAck(ackNumber, &base, &nextFrame, timer, windows)

		// Check if timeout flag is on, if so, Start windows Sending Mechanism
		if timeout {
			timeoutColor.Printf("Client:\t\n")
			counter := base
			for counter < nextFrame {
				// Send Frame on link
				sendFrame(frameList, &buf, counter, data)
				// Show Sent Frame Logs
				clientSendLog(windows, buf.Bytes(), time.Now())
				// Simulate Sending Time
				clientSendTimeSimulation()
				buf.Reset()

				// End of loop
				counter++

				// Check for ack of frame that was sent in timeout mechanism
				checkAck(ackNumber, &base, &nextFrame, timer, windows)
			}
			// Turn of timout flag
			timeout = false
		}

	}

	// End of client process
	wg.Done()
	// Calculate times that was taken to send frame
	End  = time.Now().Sub(start)
	// Wait for Server goroutine to get complete.
	wg.Wait()
	// Closing Links
	close(data)
	close(ackNumber)
	// Make main function aware
	done <- true

}
func checkAck(ackNumber chan int, base *int, nextFrame *int, timer *time.Timer, windows int) {
	// Time that should wait to get Ack
	waitTimeForFrame := 2*float64(propagation) + tf
	// Create a non-blocking statement for ackNumber channel to receive next frame ack.
	select {
	case number := <-ackNumber:
		*base = number + 1
		// test whether next ack is equal to base of windows, if was, stop timer which is used to consider a timer to enable timeout mechanism.
		if *base == *nextFrame {
			timer.Stop()
			timeout = false
		} else {
			// Consider a time based on  the ack if frames that maybe be received
			timer.Reset(time.Duration(windows) * time.Duration(waitTimeForFrame) * time.Millisecond)
			go func() {
				<-timer.C
				timeout = true
			}()
		}
	default:
		time.Sleep(1 * time.Second)

	}
}

// createFrame creates frame from message.
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

// SendFrame gets windows elements (frame list) and send them on link
func sendFrame(frameList *arraylist.List, buf *bytes.Buffer, nextFrame int, data chan []byte) {
	enc := gob.NewEncoder(buf)
	key, _ := frameList.Get(nextFrame)
	err := enc.Encode(key)
	if err != nil {
		panic(err)
	}
	//Send data on channel
	data <- buf.Bytes()
}


// clientSendLog outputs some log about frame which was sent
func clientSendLog(windows int, buffer []byte, t time.Time) {
	_, _ = clientColor.Printf("Transmitter: \tFrame %d Sent:%s %s \n", int(buffer[len(buffer)-1])%windows, buffer[1:len(buffer)-1], t)
}
// clientSendTimeSimulation simulate frame transfer time.
func clientSendTimeSimulation() {
	time.Sleep(time.Duration(propagation) * time.Millisecond)
	time.Sleep(time.Duration(tf) * time.Millisecond)
}

// This Function gets Frame size as input and returns a list of frames that are made of message which is going to be sent.
func createArrayOfFrames(fs int) *arraylist.List {

	frameList := arraylist.New()
	nextFrame := 0
	messageInByte := []byte(messages)
	buffer := make([]byte, fs)
	for len(messageInByte) != 0 {
		buffer = createFrame(messageInByte, fs)
		// Add Number
		buffer[len(buffer)-1] = byte(nextFrame)

		// Omit last value that was sent
		if len(messageInByte) > fs {
			messageInByte = messageInByte[fs-2:]
		} else {
			messageInByte = nil
		}
		// Add Frames To list
		frameList.Add(buffer)
		nextFrame++

	}
	return frameList

}
