package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"
	"gofpm/lib"
	"gofpm/lib/helper"
)

func main() {
	configPath, err := helper.GetConfigFilePath()

	if err != nil {
		log.Fatalf("Failed to find config file: %v", err)
		return
	}

	config, err := lib.LoadConfig(configPath)

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := lib.NewWorkerPool(config, ctx)

	time.Sleep(1*time.Second)

	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}

	defer listener.Close()

	log.Println("Go FastCGI Master listening on 127.0.0.1:9000")

	for {
		clientConn, err := listener.Accept()

		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go func(client net.Conn) {
			defer client.Close()
			worker := pool.NextWorker()
			backendConn, err := net.DialTimeout("tcp", worker.Address, 2*time.Second)

			if err != nil {
				log.Printf("Failed to connect to worker %d (%s): %v", worker.ID, worker.Address, err)
				return
			}

			defer backendConn.Close()

			errc := make(chan error, 2)

			go func() {
				_, err := io.Copy(backendConn, client)
				errc <- err
			}()

			go func() {
				_, err := io.Copy(client, backendConn)
				errc <- err
			}()

			<-errc
		}(clientConn)
	}
}