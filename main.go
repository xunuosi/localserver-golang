package main

import (
	"bufio"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"os"
	"time"
)

var records map[string]string

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
	case "dns":
		CreateLocalDNSServer()
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

func CreateLocalDNSServer() {
	records = map[string]string{
		"google.com": "216.58.196.142",
		"amazon.com": "176.32.103.205",
	}

	// Bind to a specific address and port
	addr, err := net.ResolveUDPAddr("udp", ":8090")
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
	fmt.Println("Create local DNS server success and executing......")

	// Run indefinitely
	for {
		// Read incoming data
		data := make([]byte, 1024)
		_, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Print the incoming data
		fmt.Printf("Received DNS request from addr:[%s]", addr.String())
		packet := gopacket.NewPacket(data, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		dns, _ := dnsPacket.(*layers.DNS)
		serveDNS(conn, addr, dns)
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error
	var ok bool
	ip, ok = records[string(request.Questions[0].Name)]
	if !ok {
		//Todo: Log no data present for the IP and handle:todo
	}
	a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	fmt.Println(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	u.WriteTo(buf.Bytes(), clientAddr)
}
