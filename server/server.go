package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/dicedb/dice/config"
	"github.com/dicedb/dice/core"
)

func RunSyncTCPServer() {
	log.Println("Starting a synchronous TCP server on", config.Host, ":", config.Port)

	var con_clients int = 0

	// listening to the configured host:port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Println("err", err)
		return
	}

	// infinite for loop whose job is i am infinitely waiting for new clients to connect
	for {
		c, err := lsnr.Accept() // this is a blocking call
		if err != nil {
			log.Println("err", err)
			return
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
			respond(cmd, c)
		}

	}
}

func readCommand(c io.ReadWriter) (*core.RedisCmd, error) { // Client → TCP → bytes → buffer → string → server
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:]) // this is a blocking call this is where we read data from client
	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c io.ReadWriter) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
}
