package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

var (
	content  = flag.String("content", "", "Content (can also be set via FLAG_SERVER_CONTENT environment variable)")
	filepath = flag.String("filepath", "~/flag.txt", "Filepath (can also be set via FLAG_SERVER_FILEPATH environment variable)")
	host     = flag.String("host", "0.0.0.0", "Host (can also be set via FLAG_SERVER_HOST environment variable)")
	port     = flag.Int("port", 0, "Port number (can also be set via FLAG_SERVER_PORT environment variable, defaults to random)")
	protocol = flag.String("protocol", "tcp", "Specify what protocol to use, permitted values are tcp, udp, and http. (can also be set via FLAG_SERVER_PROTOCOL environment variable, defaults to \"tcp\")")

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

	if envProtocol := os.Getenv("FLAG_SERVER_PROTOCOL"); envProtocol  != "" {
        protocol = &envProtocol
	}

	log.Printf("Starting new flag server instance")
	log.Printf("Host: %s", *host)
	log.Printf("Port: %d", *port)
	log.Printf("Filepath: %s", *filepath)
	log.Printf("Protocol: %s", *protocol)


	// identify flag content,
	var buffer []byte
    if *content != "" {
		buffer = []byte(*content)
		log.Printf("Using content from command line: %s", *content)
	} else if envContent := os.Getenv("FLAG_SERVER_CONTENT"); envContent != "" {
		buffer = []byte(envContent)
		log.Printf("Using content from FLAG_SERVER_CONTENT env var: %s", envContent)
	} else if *filepath != "" {
		log.Printf("Flag file %s is %d bytes", *filepath, getFileSize(*filepath))
		buffer = readFileIntoBuffer(*filepath)
        log.Printf("Actual flag is:")
        fmt.Println(string(buffer))
	}


    // use appropriate protocol
    switch *protocol {
        case "http":
            serveContentViaHTTP(buffer, *port)
        case "tcp":
            serveContentViaTcp(buffer)
        case "udp":
            serveContentViaUdp(buffer)
        default:
            log.Printf("ERROR, unknown protocol: %s", *protocol)
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

// serve content via HTTP
func serveContentViaHTTP(content []byte, port int) {

    // add an extra newline
    content = append(content, '\n')

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write(content)
    })

    // create a listener on the specified port
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
