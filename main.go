package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func scanPortWorker(host string, ports_channel <-chan int, wg *sync.WaitGroup) {

	defer wg.Done()

	for port := range ports_channel {

		address := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", address, 1*time.Second)

		if err != nil {
			// fmt.Printf("[Error] %s", err)
			// WARNING: DO NOT USE 'return' HERE , IT'LL KILL THE WORKER
		} else {
			fmt.Printf("[OPEN] Port : %d\n", port)
			conn.Close()
		}

		// wg.Done()

	}

}

func main() {

	// 1. Creaete a Channel of size say 100
	// 2. Create a worker of size 100
	// 3. Keep adding the ports to the Channel as they get empty

	var wg sync.WaitGroup

	host := "localhost"

	ports_channel := make(chan int, 100)

	for i := 0; i < 100; i++ { // worker pool of 100 workers
		wg.Add(1)
		go scanPortWorker(host, ports_channel, &wg)
	}

	for port := 1; port <= 65535; port++ {
		// wg.Add(1)
		ports_channel <- port
		// wg.Done()
	}

	// even though the ports_channel is empty , the workers are waiting for the 65536th port to be pushed into it ,they'll only stop when the ports_channel is explicitly closed

	close(ports_channel)

	wg.Wait()

}

