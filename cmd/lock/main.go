package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hirokazumiyaji/lock/lock"
)

var (
	version string
)

func main() {
	var (
		port  int
		sock  string
		cap   int
		vFlag bool
		err   error
	)

	flag.IntVar(&port, "port", 6800, "listen port")
	flag.IntVar(&port, "p", 6800, "listen port")
	flag.StringVar(&sock, "sock", "", "unix domain socket")
	flag.StringVar(&sock, "s", "", "unix domain socket")
	flag.IntVar(&cap, "cap", 1000, "capacity")
	flag.IntVar(&cap, "c", 1000, "capacity")
	flag.BoolVar(&vFlag, "version", false, "version")
	flag.BoolVar(&vFlag, "v", false, "version")

	flag.Parse()

	if vFlag {
		fmt.Printf("lock version: %s\n", version)
		return
	}

	lock.Initialize(cap)
	if err = lock.Serve(sock, port); err != nil {
		log.Fatalf("lock serve error: %v", err)
	}
}
