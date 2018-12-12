# Grupttor
Grupttor is two phase graceful application interruptor(that's why grupttor). It supports handle interrupt, notify 
application parts when you need to wait for finish and then stop the whole application, or do whatever else you want 
to do.

Grupttor was released in initial phase and you can look at the [Roadmap](#roadmap) for next steps and features

[![GoDoc](https://godoc.org/github.com/praggolabs/grupttor?status.svg)](https://godoc.org/github.com/praggolabs/grupttor)
[![GoReport](https://goreportcard.com/badge/praggolabs/grupttor)](https://goreportcard.com/report/praggolabs/grupttor)
[![Version](https://img.shields.io/badge/version-1.0.1-blue.svg)](https://github.com/praggolabs/grupttor/releases/latest)


# Table of Contents
- [How it works](#how-it-works)
- [Installing](#installing)
- [Getting Started](#getting-started)
  * [Configure grupttor](#configure-grupttor)
  * [Add hooks](#add-hooks)
  * [Run grupttor](#run-grupttor)
- [Example](#example)
- [Hooks](#hooks)
- [Handlers](#handlers)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

# How it works
Grupttor has several phases:
1) when you initialize, it switch to init state and wait for configuration(hook add, etc.)
2) you run startAndWait() method in new goroutine, you have to, because wihout it you will not 
be able to run your own code :-), then it switch to WAITING state and start wait channel for 
interrupt signal
3) when it receive interrupt signal, it switch state to INTERRUPTING and run interruptHandler
which you specify when you init grupttor
4) the interrupt handle do your bussiness and after that you run Stop method on the grupttor
to init stop phase - everything its in your hand, so you can do anything you want
5) when it receive stop signal, it just switch state to stop and run your stop handler, when you 
exit application or send kill signal whatever you need 

# Installing
Just run:
`go get -u github.com/PragGoLabs/grupttor`

## Configure grupttor
```go
    // init interruptor
    interruptor := grupttor.NewGrupttor(
    	handlers.NewWrapHandler(
            func(interrupter *grupttor.Grupttor) {
                // implement interrupt handler
                interrupter.Stop()
            },
    
            func(interrupter *grupttor.Grupttor) {
                // and stop handler
            },
        ),
        // add hooks
        []grupttor.Hook{},
    )
```

## Add hooks
```go
    // if you need add later phase
    // or add hooks manually
    interruptor.AddHook(hooks.NewSystemInterruptHook([]os.Signal{syscall.SIGKILL, syscall.SIGTERM}))
```

## Run grupttor
```go
    // and start, it will wait for signals
    go interruptor.StartAndWait()
```

# Hooks
Grupttor supports hook system, which you allow to specify logic which will run the interrupt phase.

By default there are 2 hooks:
1. SystemInterruptHook
    - it handles system interruption calls, like sigterm, sigkill, sigabrt, etc.
    - you can specify the signals you want to handle
    ```go
        interruptor.AddHook(
            hooks.NewSystemInterruptHook(
                []os.Signal{
                    syscall.SIGKILL, 
                    syscall.SIGTERM,
                }
            )
        )
    ```

2. TimedInterruptHook
    - its one shot timed interrupt action
    - you specify how long it will wait to shot interrupt task
    - if you need restart you app every hour for example and you run it in docker with restart always config
    ```go
       interruptor.AddHook(hooks.NewTimedInterruptHook(
           10*time.Second,
       ))
    ```
    
**You can implement your own hooks, just use the Hook interface and pass out.**

# Handlers
Handlers are used for handling the signals, interrupt, stop.

Now is there only WrapHandler(), which will wrap the two handler function - for keeping 
simple interface, when you need simple handler. 

I'll do a lot of handlers in future, for http, rabbitmq, etc. 
For details look at there: [Handlers](https://github.com/praggolabs/grupttor/tree/master/handlers/)

# Example
```go
    package main
    
    import (
        "fmt"
        "github.com/PragGoLabs/grupttor"
        "github.com/PragGoLabs/grupttor/hooks"
        "os"
        "syscall"
        "time"
    )

    func main() {
        interruptor := grupttor.NewGrupttor(
        	handlers.NewWrapHandler(
                func(interrupter *grupttor.Grupttor) {
                    fmt.Println("Received interrupt signal")
        
                    time.Sleep(4*time.Second)
        
                    interrupter.Stop()
                },
                func(interrupter *grupttor.Grupttor) {
                    // exit application
                    os.Exit(0)
                },
            ),
            []grupttor.Hook{},
        )
    
        interruptor.AddHook(hooks.NewSystemInterruptHook([]os.Signal{syscall.SIGKILL, syscall.SIGTERM}))
    
        go interruptor.StartAndWait()
    
        for {
            time.Sleep(2*time.Second)
    
            fmt.Println("Running loop and waiting for kill or sigterm")
        }
    }
```

# Roadmap
1. Phase deploy the initial version
2. Add support for passing the interrupt handler in concrete context not in init phase
3. Add support for multiple interrupt handlers
4. Add more interrupters

# Contributing

1. Fork it
2. Clone (`git clone https://github.com/your_username/grupttor && cd grupttor`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Do changes, Commit, Push
5. Create new pull request
6. Thanks in advance :-) 

# License

See [LICENSE.txt](https://github.com/praggolabs/grupttor/LICENSE.md)
