// Fan - Out , Fan - In Concurrency Pattern

package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func scanPortWorker(host string, ports_channel <-chan int, results_channel chan int, wg *sync.WaitGroup) {

	defer wg.Done()

	for port := range ports_channel {

		address := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

		if err != nil {
			continue
		}

		results_channel <- port

		conn.Close()

	}

}

func main() {

	var wg sync.WaitGroup

	// create 2 channels , 1 for fanning out , 1 for fanning in

	ports_channel := make(chan int, 500)
	results_channel := make(chan int, 500)

	host := "localhost"

	for workers := 0; workers < 500; workers++ {
		wg.Add(1)
		go scanPortWorker(host, ports_channel, results_channel, &wg)
	}

	for port := 1; port <= 65535; port++ {
		ports_channel <- port
	}

	// we use another single go routine to get all the ports from the results_channel and store them in a slice

	var open_ports []int

	// we need to verify whether this fanning in is complete or not since the wg WaitGroup only checks on the workers
	// we can do this by creating another channel but an empty struct in it
	// instead we could have a bool and put in a true and then take it out when the fanning - in is complete but then , it takes 1 - byte for storing that boolean
	// instead of we use an empty struct it only takes 0 - bytes

	done := make(chan struct{})

	go func() {

		for open_port := range results_channel {
			open_ports = append(open_ports, open_port)
		}

		close(done)

	}()

	// NOTE: The below order is rlly rlly important

	close(ports_channel) // stops the channel from taking new inputs

	wg.Wait() // waits for the workers to finish their job

	close(results_channel) // stops the results channel from taking in new inputs

	<-done // wait for that faning - in function to do it's job

	fmt.Print(open_ports)

}

