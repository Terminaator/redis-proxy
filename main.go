package main

import "log"

func main() {

	port := ":8080"

	log.Println("starting proxy:", port)

	proxy := Proxy{ip: port}

	proxy.start()
}
