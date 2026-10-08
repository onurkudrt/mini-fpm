package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	defer file.Close()

	cfg := &Config{
		ListenAddr:  "127.0.0.1:9000",
		StartPort:   9001,
		WorkerCount: 4,
		MaxRequests: 1000,
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())


		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}


		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
			case "php_path":
				cfg.PHPPath = value
			case "listen_address":
				cfg.ListenAddr = value
			case "start_port":
				if port, err := strconv.Atoi(value); err == nil {
					cfg.StartPort = port
				}
			case "worker_count":
				if count, err := strconv.Atoi(value); err == nil {
					cfg.WorkerCount = count
				}
			case "max_requests":
				if max, err := strconv.Atoi(value); err == nil {
					cfg.MaxRequests = max
				}	
		}
	}

	if cfg.PHPPath == "" {
		return nil, fmt.Errorf("php_path is required in config")
	}

	return cfg, nil
}