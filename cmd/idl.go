package main

import (
	"fmt"
	"path"
	"runtime"
	"sync"

	"github.com/h2non/bimg"
	"github.com/tmshv/idl/internal/config"
	"github.com/tmshv/idl/internal/csv"
	"github.com/tmshv/idl/internal/dl"
	"github.com/tmshv/idl/internal/preprocess/compose"
	"github.com/tmshv/idl/internal/preprocess/resize"
	"github.com/tmshv/idl/internal/utils"
)

type image struct {
	file string
	data []byte
}

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cfg)

	cpu := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n", cpu)

	recordCh := make(chan csv.Record, cpu)
	go (func() {
		reader := csv.New(cfg.Fields.URL, cfg.Fields.File)
		err := reader.Read(cfg.Input, recordCh)
		if err != nil {
			fmt.Printf("Failed to read CSV file: %v\n", err)
		}
	})()

	dlCh := make(chan image, cfg.Workers)
	go func() {
		defer close(dlCh)
		dl := dl.New(cfg.Timeout, 1)
		wg := sync.WaitGroup{}
		for rec := range recordCh {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				file := path.Join(cfg.Dir, rec.File)
				if !cfg.Reload && utils.FileExists(file) {
					fmt.Printf("Skip: %s\n", file)
					return
				}

				body, err := dl.Download(rec.URL)
				if err != nil {
					fmt.Printf("Failed to download image: %v\n", err)
					return
				}

				dlCh <- image{file: file, data: body}
			}(&wg)
		}
		wg.Wait()
	}()

	imgCh := make(chan image, cpu)
	go func() {
		defer close(imgCh)

		prep := compose.New(
			resize.New(cfg.Resize[0], cfg.Resize[1]),
		)

		wg := sync.WaitGroup{}
		for range cfg.Workers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := range dlCh {
					i.data, err = prep.Run(i.data)
					if err != nil {
						fmt.Printf("Failed to preprocess image: %v\n", err)
						return
					}
					imgCh <- i
				}
			}()
		}
		wg.Wait()
	}()

	for i := range imgCh {
		err := bimg.Write(i.file, i.data)
		if err != nil {
			fmt.Printf("Failed to write image: %v\n", err)
			continue
		}
		fmt.Printf("OK: %v\n", i.file)
	}
}
