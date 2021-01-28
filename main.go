package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"sync"
	"time"
)

var (
	messages = `Take a trip inside my head..Leave your sticks and stones
				Broken bones, I'm left for dead
				But still I carry on
				When I am down
				I carry on
				When it's cold in this wild, wild world
				Everyone's trying to dig your grave
				I carry on
				When you're told you don't the fit the mold
				Now everybody's got a say
				I carry on
				When the madness all around us starts to take it's toll
				I carry on
				It's a long, dark, winding road we're on
				Oh, I carry on`
	propagation  int
	bandwidth    int
	fs           int
	tf           = 0.0

	// Color Objects
	prompt       = color.New(color.FgGreen)
	serverColor  = color.New(color.FgHiMagenta)
	clientColor  = color.New(color.FgHiYellow)
	ackColor     = color.New(color.FgHiGreen)
	timeoutColor = color.New(color.FgHiRed)

	timeout      = false
	wg           sync.WaitGroup
	End          time.Duration
	start        time.Time
)

func main() {
	wg.Add(2)
	// Get Propagation Time & Bandwidth
	_, _ = prompt.Println("Enter Propagation Time and Bandwidth: ")
	_, err := fmt.Scanf("%d %d\n", &propagation, &bandwidth)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	// Get Frame Size and Windows Size
	var windows int
	prompt.Println("Enter Frame  and Windows Size ")
	_, err = fmt.Scanf("%d %d", &fs, &windows)
	if err != nil {
		log.Fatalln("Fatal Error")
	}

	// Calculate Transmission Time
	tf = float64(fs / bandwidth)

	// Create a Channel to Transfer Data With Server and Client
	data := make(chan []byte, windows)

	// Create a Channel Check The End of Comm
	done := make(chan bool, 1)

	// Create a Channel to server send its next ack number
	ackNumber := make(chan int, windows)

	start = time.Now()
	// Start Communication
	go client(windows, data, ackNumber, done)
	go server(windows, data, ackNumber)

	// Wait Till Process End
	<-done

	// End of simulation and Print Total Time and Average for every frame
	_, _ = prompt.Println("Total Time: ", time.Now().Sub(start))
	fmt.Printf("Avarage Time For Every Frame: %f",
		float64(End.Milliseconds()/int64(numberOfFrame)))

}
