package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/tmshv/idl/internal/config"
	"github.com/tmshv/idl/internal/idl"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Options: %+v\n", cfg)

	cpu := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n", cpu)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	idl := idl.New(cpu)
	err = idl.Run(ctx, cfg)
	if err != nil {
		panic(err)
	}
}
