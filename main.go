package main

import "log"

var (
	sentinel Sentinel
	proxy    Proxy
)

func main() {

	port := ":8080"

	log.Println("starting proxy:", port)

	sentinel = Sentinel{}

	sentinel.start()

	proxy = Proxy{ip: port}

	proxy.start()
}
