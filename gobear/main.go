package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // Setup main logger
    logFile, err := os.OpenFile("/home/stormrider/mnt/Dropbox/var/log/gobear/gobear.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Setup performance metrics logger
    perfLogFile, err := os.OpenFile("/home/stormrider/mnt/Dropbox/var/log/gobear/gobear.perform.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening performance log file: %v", err)
    }
    defer perfLogFile.Close()
    perfLog := log.New(perfLogFile, "PERF: ", log.LstdFlags|log.Lshortfile)

    log.Println("Starting up...")

    // Setup Prometheus metrics
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "my_counter",
        Help: "An example counter metric",
    })
    prometheus.MustRegister(counter)

    // Setup HTTP server for Prometheus
    http.Handle("/metrics", promhttp.Handler())
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // Setup signal handling
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    // Context to handle cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        sig := <-sigs
        log.Printf("Received signal: %v\n", sig)
        cancel() // Trigger cancellation
    }()

    // Main program logic
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                counter.Inc()
                perfLog.Printf("Counter incremented, current value: %v", counter)
                time.Sleep(1 * time.Second)
            }
        }
    }()

    // Wait for signal
    <-ctx.Done()
    log.Println("Shutting down")
}
