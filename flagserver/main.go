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
)


func main() {
	flag.Parse()
	log.Printf("Starting new file server instance")

	log.Printf("File %s is %d bytes", *filepath, getFileSize(*filepath))

	content := readFileIntoBuffer(*filepath)
	log.Printf("File content is:")
	fmt.Println(string(content))

	serveContent(content)
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
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return empty
	}
	defer file.Close()

	fileContent, err := os.ReadFile(filename)
	if err == io.EOF {
		log.Printf("Finished reading file")
	}

	return fileContent
}

// create network listener and serve content
func serveContent(content []byte) {
	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening for connections on %s", listener.Addr().String())

	// process each new connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go processNewConnection(conn, content)
		}
	}
}

// process incoming connections
func processNewConnection(conn net.Conn, content []byte) {
	defer conn.Close()

	// identify source IP
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		log.Printf("Accepted new connection from client on address: %s", addr.IP.String())
	}

	conn.Write(content)
    log.Printf("Finished serving content")

}

