package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	SENTINEL_COMMAND string
)

type Sentinel struct {
	redis    string
	Sentinel *net.TCPAddr
}

func (s *Sentinel) read(conn *net.TCPConn) (string, error) {
	err := s.write(conn)

	buffer := make([]byte, 256)

	_, err = conn.Read(buffer)

	parts := strings.Split(string(buffer), "\r\n")

	if err != nil || len(parts) < 5 {
		return "", errors.New("failed to get sentinel")
	}

	return fmt.Sprintf("%s:%s", parts[2], parts[4]), err
}

func (s *Sentinel) write(conn *net.TCPConn) error {
	_, err := conn.Write([]byte(fmt.Sprintf("sentinel get-master-addr-by-name %s\n", "mymaster")))
	return err
}

func (s *Sentinel) checkMaster(ip string) {
	if len(s.redis) == 0 && s.redis != ip {
		s.redis = ip
	}
}

func (s *Sentinel) getMaster(conn *net.TCPConn) {
	for {
		if ip, err := s.read(conn); err == nil {
			s.checkMaster(ip)
		} else {
			go s.connect()
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func (s *Sentinel) connect() {
	if conn, err := net.DialTCP("tcp", nil, s.Sentinel); err == nil {
		s.getMaster(conn)
	} else {
		go s.connect()
	}
}

func (s *Sentinel) init() {
	adr, err := net.ResolveTCPAddr("tcp", ":26379")
	if err != nil {
		log.Fatal("Failed to resolve sentinel address", err)
	}

	s.Sentinel = adr
}

func (s *Sentinel) start() {
	s.init()
	s.connect()
}
