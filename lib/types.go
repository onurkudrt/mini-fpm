package lib

import (
	"sync"
	"os/exec"
)

type Worker struct {
	ID int
	Address string
	cmd *exec.Cmd
}

type WorkerPool struct {
	cursor uint64
	workers []*Worker
	mu sync.RWMutex
}

type Config struct {
	PHPPath     string
	ListenAddr  string
	StartPort   int
	WorkerCount int
	MaxRequests   int
}