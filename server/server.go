package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/dicedb/dice/config"
)

func RunSyncTCPServer() {
	log.Println("Starting a synchronous TCP server on", config.Host, ":", config.Port)

	var con_clients int = 0

	// listening to the configured host:port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	// infinite for loop whose job is i am infinitely waiting for new clients to connect
	for {
		c, err := lsnr.Accept() // this is a blocking call
		if err != nil {
			panic(err)
		}

		con_clients += 1 // now i have these many clients connected
		log.Println("New client connected with address:", c.RemoteAddr().String(), "concurrent clients", con_clients)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("New client connected with address:", c.RemoteAddr().String(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("Error reading command:", err)
			}
			log.Println("Received command:", cmd)
			if err = respond(cmd, c); err != nil {
				log.Print("Error write:", err)
			}
		}

	}
}

func readCommand(c net.Conn) (string, error) { // Client → TCP → bytes → buffer → string → server
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:]) // this is a blocking call this is where we read data from client
	if err != nil {
		return "", err
	}
	return string(buf[0:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}
