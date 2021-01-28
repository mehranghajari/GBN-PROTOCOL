package main

import (
	"bytes"
	"encoding/gob"
	"github.com/emirpasic/gods/lists/arraylist"
	"time"
	"fmt"
)

func server(windows int, data chan []byte, ackNumber chan int, done chan bool) {

	// Declare a base var and next frame
	base := 0
	nextFrame := 0

	// Timer
	timer := time.NewTimer(4 * time.Second)
	timer.Stop()

	frameList := createArrayOfFrames(fs)
	var buf bytes.Buffer
	count := 0
	for count < frameList.Size() {
		if nextFrame < base+windows && timeout != true {

			sendFrame(frameList, &buf, nextFrame, data)
			clientSendLog(windows, buf.Bytes(), time.Now())
			clientSendTimeSimulation()
			buf.Reset()
			nextFrame++
			count++

		}
		checkAck(ackNumber, &base, &nextFrame, timer)

		if timeout {
			fmt.Println("Timeout")
			counter := base
			for counter < nextFrame {
				sendFrame(frameList, &buf, counter, data)
				clientSendLog(windows, buf.Bytes(), time.Now())
				clientSendTimeSimulation()
				buf.Reset()
				counter++
				checkAck(ackNumber, &base, &nextFrame, timer)
			}
			timeout = false
		}

	}
	wg.Done()
	wg.Wait()
	close(data)
	close(ackNumber)
	done <- true

}
func checkAck(ackNumber chan int, base *int, nextFrame *int, timer *time.Timer) {
	select {
	case number := <-ackNumber:
		*base = number + 1
		if *base == *nextFrame {
			timer.Stop()
			timeout = false
		} else {
			timer.Reset(2* time.Second)
			go func() {
				<-timer.C
				timeout = true
			}()
		}
	default:
		time.Sleep(1 * time.Second)

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

func clientSendLog(windows int, buffer []byte, t time.Time) {
	_, _ = clientColor.Printf("Transmitter: \tFrame %d Sent:%s %s \n", int(buffer[len(buffer)-1])%windows, buffer[1:len(buffer)-1], t)
}
func clientSendTimeSimulation() {
	time.Sleep(time.Duration(propagation) * time.Millisecond)
	time.Sleep(time.Duration(tf) * time.Millisecond)
}
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
		frameList.Add(buffer)
		nextFrame++

	}
	return frameList

}