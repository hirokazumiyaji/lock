package lock

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func Serve(sock string, port int) error {
	l, err := listen(sock, port)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	closed := false
	go func() {
		<-sig
		cancel()
		closed = true
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil && closed {
			break
		}
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		wg.Add(1)
		go accept(ctx, &wg, conn)
	}

	if !closed {
		l.Close()
	}

	wg.Wait()

	return nil
}

func listen(sock string, port int) (net.Listener, error) {
	if sock == "" {
		return net.Listen("tcp", fmt.Sprintf(":%d", port))
	} else {
		return net.Listen("unix", sock)
	}
}

func accept(ctx context.Context, wg *sync.WaitGroup, c net.Conn) {
	closed := false

	go func() {
		<-ctx.Done()
		closed = true
		c.Close()
	}()

	defer func() {
		if !closed {
			c.Close()
		}
		wg.Done()
	}()

	r := bufio.NewReader(c)

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil && closed {
			break
		}
		if err != nil {
			log.Printf("%v", err)
			break
		}

		cmd := string(line)
		parts := strings.Split(cmd, " ")
		var res []byte
		if len(parts) == 2 {
			switch parts[0] {
			case "lock":
				if lock(parts[1]) {
					res = []byte("true\n")
				} else {
					res = []byte("false\n")
				}
			case "unlock":
				if unlock(parts[1]) {
					res = []byte("true\n")
				} else {
					res = []byte("false\n")
				}
			default:
				log.Println("%v", err)
				res = []byte("cmd not found\n")
			}
		} else {
			res = []byte("cmd format error\n")
		}

		_, err = c.Write(res)
		if err != nil {
			log.Printf("%v", err)
			break
		}
	}
}
