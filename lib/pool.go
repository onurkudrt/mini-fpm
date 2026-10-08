package lib

import (
	"context"
	"log"
	"os"
	"time"
	"fmt"
	"os/exec"
	"sync/atomic"
)

func NewWorkerPool(config *Config, ctx context.Context) *WorkerPool { 
	pool := &WorkerPool {
		workers: make([]*Worker, config.WorkerCount),
	}

	for i := 0; i < config.WorkerCount; i++ {
		addr := fmt.Sprintf("127.0.0.1:%d", config.StartPort+i);
		worker := &Worker {
			ID: i,
			Address: addr,
		}

		pool.workers[i] = worker
		go supervisor(ctx, worker, config)
	}

	return pool
}

func supervisor(ctx context.Context, worker *Worker, config *Config) {
	for {
		select {
			case <-ctx.Done():
				return
			default:
		}

		cmd := exec.CommandContext(ctx, config.PHPPath, "-b", worker.Address)

		maxRequests := config.MaxRequests
		
		if maxRequests <= 0 {
			maxRequests = 1000 // default val if not set
		}

		// cmd.Env = append(os.Environ(), "PHP_FCGI_CHILDREN=1", "PHP_FCGI_MAX_REQUESTS=1000")
		cmd.Env = append(os.Environ(), fmt.Sprintf("PHP_FCGI_MAX_REQUESTS=%d", maxRequests))
		
		worker.cmd = cmd

		log.Printf("[Worker %d] Starting php-cgi on %s\n", worker.ID, worker.Address)
		
		err := cmd.Run()

		if err != nil && ctx.Err() == nil {
			log.Printf("[Worker %d] Process exited: %v. Restarting...", worker.ID, err)
		}

		time.Sleep(500* time.Millisecond)
	}
}

func (pool *WorkerPool) NextWorker() *Worker {
	idx := atomic.AddUint64(&pool.cursor, 1) % uint64(len(pool.workers))
	return pool.workers[idx]
}