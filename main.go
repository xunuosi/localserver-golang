package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// addrs, err := net.LookupHost("www.finalfantasyxiv.com")
	// if err != nil {
	// 	log.Print(err)
	// } else {
	// 	log.Printf("addrs:[%+v]", addrs)
	// }

	// const n = 300
	// var wg sync.WaitGroup
	// wg.Add(n)
	// for i := 0; i < n; i++ {
	// 	go func() {
	// 		_, err := net.LookupHost("www.finalfantasyxiv.com")
	// 		if err != nil {
	// 			fmt.Println("error", err)
	// 		}
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()

	// fmt.Println("complete......")
	arguments := os.Args
	connectType := arguments[1]
	TestLocalServer(connectType)
}

func TestLocalServer(connectType string) {
	switch connectType {
	case "tcp":
		CreateLocalTCPServer()
	case "udp":
		CreateLocalUDPServer()
	}
}

func CreateLocalUDPServer() {
	// Bind to a specific address and port
	addr, err := net.ResolveUDPAddr("udp", ":80")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen on the address
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Create local UDP server success and executing......")

	// Run indefinitely
	for {
		// Read incoming data
		data := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the incoming data
		fmt.Printf("Received: %s from %s\n", string(data[:n]), addr.String())

		// Echo the data back to the client
		_, err = conn.WriteToUDP(data[:n], addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func CreateLocalTCPServer() {
	tcpListener, err := net.Listen("tcp", ":80")

	if err != nil {
		fmt.Printf("Create local TCP server failed:[%s]", err.Error())
		return
	}
	fmt.Printf("Create local TCP server success and executing......")
	defer tcpListener.Close()

	conn, err := tcpListener.Accept()

	if err != nil {
		fmt.Printf("Create local TCP server failed:[%s]", err.Error())
		return
	}
	// defer conn.Close()

	for {
		readString, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("local TCP server read err:[%s]", err.Error())
			return
		}
		fmt.Printf("local TCP server read client string:[%s]", readString)
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		_, err = conn.Write([]byte(myTime))
		if err != nil {
			fmt.Printf("local TCP server send err:[%s]", err.Error())
			return
		}
		fmt.Printf("local TCP server send client success and close connection")
	}

}
