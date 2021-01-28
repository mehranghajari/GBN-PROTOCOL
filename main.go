package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"sync"
	"time"
)

var (
	messages     = "MehranAli Mohammad Mehrshad Hossein Iman Aida Mahdie. In the name Of god. my name is mehran ghajari"
	propagation  int
	bandwidth    int
	fs           int
	tf           = 0.0
	prompt       = color.New(color.FgGreen)
	serverColor  = color.New(color.FgHiMagenta)
	clientColor  = color.New(color.FgHiYellow)
	ackColor     = color.New(color.FgHiGreen)
	timeoutColor = color.New(color.FgHiRed)
	timerFlag    = false
	timeout      = false
	wg sync.WaitGroup
)


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
	windows := 127
	// Calculate Transmission Time
	tf = float64(fs / bandwidth)

	// Create a Channel to Transfer Data With Server and Client
	data := make(chan []byte, windows)

	// Create a Channel Check The End of Comm
	done := make(chan bool, 1)

	// Create a Channel to client send its next ack number
	ackNumber := make(chan int, windows)

	start := time.Now()
	// Start Communication
	go server(windows, data, ackNumber, done)
	go client(windows, data, ackNumber)

	// Wait Till Process End
	<-done
	fmt.Println("Total Time: ", time.Now().Sub(start))

}
