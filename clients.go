package main

import "net"

type Clients struct {
	redis    string
	Sentinel *net.TCPAddr
}
