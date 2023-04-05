package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	filepath = flag.String("f", "~/flag.txt", "Filepath (can also be set via FLAG_SERVER_FILEPATH environment variable)")
	host     = flag.String("h", "0.0.0.0", "Host (can also be set via FLAG_SERVER_HOST environment variable)")
	port     = flag.Int("p", 0, "Port number (can also be set via FLAG_SERVER_PORT environment variable (default to random)")
	udp      = flag.Bool("u", false, "Use UDP instead of TCP (can also be set via FLAG_SERVER_UDP environment variable (default to \"false\")")
)

func main() {
	flag.Parse()

	// read environment variables
	if envHost := os.Getenv("FLAG_SERVER_HOST"); envHost != "" {
		host = &envHost
	}
	if envPort := os.Getenv("FLAG_SERVER_PORT"); envPort != "" {
		if portVal, err := strconv.Atoi(envPort); err == nil {
			port = &portVal
		}
	}
	if envFilepath := os.Getenv("FLAG_SERVER_FILEPATH"); envFilepath != "" {
		filepath = &envFilepath
	}
	if envUdp := os.Getenv("FLAG_SERVER_UDP"); envUdp != "" {
		if udpVal, err := strconv.ParseBool(envUdp); err == nil {
			udp = &udpVal
		}
	}

	log.Printf("Starting new flag server instance")
	log.Printf("Host: %s", *host)
	log.Printf("Port: %d", *port)
	log.Printf("Filepath: %s", *filepath)
	log.Printf("Use UDP: %t", *udp)
	log.Printf("Flag file %s is %d bytes", *filepath, getFileSize(*filepath))

	content := readFileIntoBuffer(*filepath)
	log.Printf("Flag is:")
	fmt.Println(string(content))

	if *udp {
		serveContentViaUdp(content)
	} else {
		serveContentViaTcp(content)
	}
}

// get the file size
func getFileSize(filename string) int64 {
	file, err := os.Stat(filename)
	if err != nil {
		log.Printf("Error identifying file: %s", err.Error())
		return -1
	}
	return file.Size()
}

// read file into a buffer
func readFileIntoBuffer(filename string) []byte {
	var empty []byte
	// read the file to serve
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
        os.Exit(1)
		return empty
	}
	if err == io.EOF {
		log.Printf("Finished reading file")
	}
	return fileContent
}

// create network listener and serve content via TCP
func serveContentViaTcp(content []byte) {
	addr := fmt.Sprintf("%s:%d", *host, *port)

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Printf("Listening for TCP connections on %s", tcpListener.Addr().String())

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go processNewTcpConnection(conn, content)
		}
	}
}

// process incoming TCP connections
func processNewTcpConnection(conn net.Conn, content []byte) {
	defer conn.Close()

	// identify source IP
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		log.Printf("Accepted new TCP connection from client on address: %s", addr.IP.String())
	}

	conn.Write(content)
	log.Printf("Finished serving content over TCP")

}

// create network listener and serve content via UDP
func serveContentViaUdp(content []byte) {
	addr := fmt.Sprintf("%s:%d", *host, *port)

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("Listening for UDP connections on %s", udpConn.LocalAddr().String())

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading data from client: %s", err.Error())
		} else {
			log.Printf("Accepted new UDP connection from client on address: %s", addr.IP.String())
			udpConn.WriteToUDP(content, addr)
			log.Printf("Finished serving content over UDP")
		}
	}
}
