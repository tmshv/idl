package main

import (
	"context"
	"fmt"
	"runtime"

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

	idl := idl.New(cpu)
	ctx := context.Background()
	err = idl.Run(ctx, cfg)
	if err != nil {
		panic(err)
	}
}
