package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	host     = flag.String("h", "localhost", "Host")
	port     = flag.Int("p", 0, "Port")
	filepath = flag.String("f", "~/flag.txt", "filepath")
	udp      = flag.Bool("u", false, "use UDP instead of TCP")
)

func main() {
	flag.Parse()
	log.Printf("Starting new file server instance")

	log.Printf("File %s is %d bytes", *filepath, getFileSize(*filepath))

	content := readFileIntoBuffer(*filepath)
	log.Printf("File content is:")
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
