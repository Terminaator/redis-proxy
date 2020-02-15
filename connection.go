package main

import (
	"bytes"
	"io"
	"log"
)

type Connection struct {
	ip    string
	rwc   io.ReadWriteCloser
	redis Redis
}

func (c *Connection) close() {
	log.Println("connection closed", c.ip)
	c.redis.close()
	c.rwc.Close()
}

func (c *Connection) outPipe(out []byte) {
	c.rwc.Write(out)
}

func (c *Connection) redisPipe(message []byte) {
	log.Println("message from", c.ip, message)
	out := make([]byte, 128)
	c.redis.do(message, out)
	c.outPipe(bytes.Trim(out, "\x00"))
}

func (c *Connection) inPipe() {
	for {
		buffer := make([]byte, 128)

		if _, err := c.rwc.Read(buffer); err == nil {
			go c.redisPipe(bytes.Trim(buffer, "\x00"))
		} else {
			c.close()
			break
		}

	}
}
