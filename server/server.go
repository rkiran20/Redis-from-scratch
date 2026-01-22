package server

import (
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
			cmds, err := readCommands(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("New client connected with address:", c.RemoteAddr().String(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("Error reading command:", err)
			}
			log.Println("Received command:", cmds)
			respond(cmds, c)
		}

	}
}

func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return as, nil
}

func readCommands(c io.ReadWriter) (core.RedisCmds, error) { // Client → TCP → bytes → buffer → string → server
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:]) // this is a blocking call this is where we read data from client
	if err != nil {
		return nil, err
	}
	values, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}

	var cmds []*core.RedisCmd = make([]*core.RedisCmd, 0)
	for _, value := range values {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil
}

func respond(cmds core.RedisCmds, c io.ReadWriter) {
	core.EvalAndRespond(cmds, c)
}
