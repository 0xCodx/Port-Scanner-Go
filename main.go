package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
func main() {
	ports := make(chan int, 100)
	results := make(chan int)
	var oPorts []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	go func() {
		for i := 1; i <= 65535; i++ {
			ports <- i
		}
	}()
	for i := 0; i < 1000; i++ {
		port := <-results
		if port != 0 {
			oPorts = append(oPorts, port)
		}
	}
	defer close(ports)
	defer close(results)
	sort.Ints(oPorts)
	for _, port := range oPorts {
		fmt.Printf("%d open\n", port)
	}
}
