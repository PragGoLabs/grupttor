# Handlers

# Table of Contents
- [WrapHandler](#wraphandler)
- [HTTPServerHandler](#httpserverhandler)

## WrapHandler
Simple wrap handler
```go
    	handlers.NewWrapHandler(
    		func(interrupter *grupttor.Grupttor) {}, // interrupt handler
    		func(interrupter *grupttor.Grupttor) {}, // stop handler
    	)
```

## HTTPServerHandler
Http handler for safe interruption of http server, example of usage:
```go
    package main
    
    import (
        "fmt"
        "github.com/PragGoLabs/grupttor"
        "github.com/PragGoLabs/grupttor/handlers"
        "github.com/PragGoLabs/grupttor/hooks"
        "net/http"
        "os"
        "sync"
        "syscall"
    )
    
    func main() {
        wg := &sync.WaitGroup{}
    
        fmt.Println(os.Getpid())
        server := &http.Server{Addr: ":8083"}
    
        interruptor := grupttor.NewGrupttor(
            handlers.NewHTPPServerHandler(server, wg, nil),
            []grupttor.Hook{},
        )
    
        interruptor.AddHook(hooks.NewSystemInterruptHook([]os.Signal{syscall.SIGKILL, syscall.SIGTERM}))
    
        go interruptor.StartAndWait()
        wg.Add(1)
    
        http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
            writer.WriteHeader(200)
            writer.Write([]byte("Everything is fine"))
        })
    
        wg.Add(1)
        server.ListenAndServe()
    
        // wait for both
        wg.Wait()
    }
```