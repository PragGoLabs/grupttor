# Handlers

# Table of Contents
- [WrapHandler](#wraphandler)

## WrapHandler
```go
    	handlers.NewWrapHandler(
    		func(interrupter *grupttor.Grupttor) {}, // interrupt handler
    		func(interrupter *grupttor.Grupttor) {}, // stop handler
    	)
```