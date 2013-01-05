package main

import (
    "github.com/howeyc/fsnotify"
    "log"
    "time"
)

func main() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

// Process events
    go func() {
        for {
            select {
            case ev := <-watcher.Event:
                log.Println("event:", ev)
            case err := <-watcher.Error:
                log.Println("error:", err)
            }
        }
    }()

    err = watcher.Watch("/tmp")
    if err != nil {
        log.Fatal(err)
    }

    /* ... do stuff ... */
    time.Sleep(10 * time.Minute)
    watcher.Close()
}
