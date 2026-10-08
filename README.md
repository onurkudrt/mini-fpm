# mini-fpm

A lightweight, Go-based application server and process supervisor designed to manage `php-cgi` worker pools and proxy incoming web traffic on Windows.

---

## How It Works

`mini-fpm` listens on a single port (`:9000`) for incoming requests from a web server like Nginx or Apache. When a connection arrives, it uses a round-robin selector (`pool.NextWorker()`) to route the connection to an active background `php-cgi` worker. It also supervises the worker processes, handling auto-restarts and enforcing request limits.

```text
 [ Web Client ]
       │
       ▼
   [ Nginx / Apache ]  (Reverse Proxy)
       │
       ▼ (TCP :9000)
 ┌────────────────────────────────────────┐
 │   mini-fpm (Go Application Server)     │
 │   - TCP Listener & Stream Proxy        │
 │   - Round-Robin Worker Selection       │
 │   - Process Supervisor & Auto-Restart  │
 └─────┬──────────────┬──────────────┬────┘
       │              │              │
       ▼              ▼              ▼
(:9001 Worker) (:9002 Worker) (:9003 Worker)
```


## Features

* **TCP Stream Proxying:** Listens on port 9000 and tunnels incoming client connections to available background workers.

* **Process Supervision & Auto-Healing:** Monitors background php-cgi workers and instantly restarts them if they crash or exit.

* **Memory Management & Recycling:** Automatically recycles workers after they hit their request limits (PHP_FCGI_MAX_REQUESTS) to prevent memory bloat.

* **Graceful Shutdown:** Listens for OS signals (SIGINT, SIGTERM) to cleanly terminate all background subprocesses.

## Getting Started

```bash
git clone https://github.com/onurkudrt/mini-fpm.git
cd mini-fpm
go run main.go
```

### Configure Your Web Server

Check the examples/ directory for ready-to-use virtual host configurations:

* **Nginx**: See ``examples/nginx.conf``
* **Apache**: See ``examples/apache.conf``

**NOTE:** The application currently listens on 127.0.0.1:9000 by default.