package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const timeout = time.Second * 5

func dialPort(ctx context.Context, host, port string) (Status, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(host, port))
	if err != nil {
		if ctx.Err() != nil {
			return StatusError, ctx.Err()
		}

		if netErr, ok := errors.AsType[net.Error](err); ok && netErr.Timeout() {
			return StatusFiltered, nil
		}

		if errors.Is(err, syscall.ECONNREFUSED) {
			return StatusClosed, nil
		}

		return StatusError, err
	}
	defer conn.Close()

	return StatusOpen, nil
}

func worker(ctx context.Context, cancel context.CancelCauseFunc, host string, tasks <-chan string, islocalAddress bool) {
	for {
		select {
		case <-ctx.Done():
			return
		case port, ok := <-tasks:
			if !ok {
				return
			}

			status, err := dialPort(ctx, host, port)
			if err != nil {
				cancel(err)
				return
			}

			if status == StatusClosed {
				continue
			}

			if !islocalAddress && status == StatusFiltered {
				continue
			}

			fmt.Printf("%-10s %-10s\n", port+"/tcp", status)
		}
	}
}

func usage() {
	fmt.Println("\nUsage")
	fmt.Println("portscan [flags] <host> <port(s)>")
	flag.PrintDefaults()
	fmt.Println(`
Examples
portscan localhost 80
portscan -workers 2 localhost 80,443,25565
portscan -workers 6 localhost 1-1024
portscan -workers 6 localhost --
`)
}

var wg sync.WaitGroup

func main() {
	workers := flag.Int("workers", 4, "the amount of concurrent workers")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("host and port(s) are required")
		usage()
		os.Exit(2)
	}

	args := flag.Args()

	addr, err := netip.ParseAddr(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	ports, err := parsePorts(args[1])
	if err != nil {
		fmt.Println(err)
		usage()
		os.Exit(2)
	}

	isLocalAddress := addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast()

	fmt.Printf("Starting portscan at %s\n", time.Now().Format(time.DateTime))
	fmt.Printf("%-10s %-10s\n", "PORT", "STATUS")

	tasks := make(chan string)

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancelCause(signalCtx)
	defer cancel(nil)

	for range *workers {
		wg.Go(func() {
			worker(ctx, cancel, host, tasks, isLocalAddress)
		})
	}

	go func(ctx context.Context) {
		defer close(tasks)
		for _, port := range ports {
			select {
			case <-ctx.Done():
				return
			case tasks <- port:
			}
		}
	}(ctx)

	wg.Wait()

	if err := context.Cause(ctx); err != nil {
		fmt.Printf("Scan failed: %v\n", err)
	}
}
