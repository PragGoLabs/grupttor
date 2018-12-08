# Handlers

# Table of Contents
- [WrapHandler](#wraphandler)
- [HTTPServerHandler](#httpserverhandler)
- [AmqpHandler](#amqphandler)

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

## AmqpHandler
Allows you to safely interrupt amqp consumer and after that exit the application.

```go
    package main
    
    import (
        "fmt"
        "github.com/PragGoLabs/grupttor"
        "github.com/PragGoLabs/grupttor/handlers"
        "github.com/PragGoLabs/grupttor/hooks"
        "github.com/streadway/amqp"
        "os"
        "syscall"
        "time"
    )
    
    var consumerTag = "randomTag"
    var queueName = "testQueue"
    
    func main() {
        session, err := amqp.Dial("amqp://guest:guest@localhost:5672//")
        if err != nil {
            panic(err)
        }
        channel, err := session.Channel()
        if err != nil {
            panic(err)
        }
    
        if err := channel.Qos(1, 0, false); err != nil {
            panic(err)
        }
    
        // interruptor initialization
        interruptor := grupttor.NewGrupttor(
            handlers.NewAmqpHandler(channel, consumerTag),
            []grupttor.Hook{},
        )
        interruptor.AddHook(hooks.NewSystemInterruptHook([]os.Signal{syscall.SIGKILL, syscall.SIGTERM}))
    
        // start consuming the queue
        deliveries, err := channel.Consume(queueName, consumerTag, false, false, false, false, nil) // todo: add consumer tag
        if err != nil {
            panic(err)
        }
    
        consumerFunc := consume()
        go consumerFunc(deliveries, interruptor)
    
        // start and wait interruptor
        interruptor.StartAndWait()
    }
    
    func consume() func(deliveries <-chan amqp.Delivery, interruptor *grupttor.Grupttor) {
        return func(deliveries <-chan amqp.Delivery, interruptor *grupttor.Grupttor) {
            defer func() {
                // we will wait until the all messages are read from buffer
                if interruptor.IsInterrupting() {
                    interruptor.Stop()
                }
            }()
    
            for d := range deliveries {
                // just print the message
                fmt.Println(string(d.Body))
                time.Sleep(5*time.Second)
    
                d.Ack(false)
            }
        }
    }

```