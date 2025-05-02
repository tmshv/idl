package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/h2non/bimg"
	"github.com/tmshv/idl/internal/config"
	"github.com/tmshv/idl/internal/csv"
	"github.com/tmshv/idl/internal/dl"
	"github.com/tmshv/idl/internal/preprocess"
	"github.com/tmshv/idl/internal/preprocess/resize"
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
		for rec := range recordCh {
			body, err := dl.Download(rec.URL)
			if err != nil {
				fmt.Printf("Failed to download image: %v\n", err)
				continue
			}
			fmt.Printf("Record: %v\n", rec)
			dlCh <- image{file: rec.File, data: body}
		}
	}()

	imgCh := make(chan image, cpu)
	go func() {
		defer close(imgCh)

		preps := []preprocess.Preprocessor{
			resize.New(cfg.Resize[0], cfg.Resize[1]),
		}

		for i := range dlCh {
			for _, prep := range preps {
				i.data, err = prep.Run(i.data)
				if err != nil {
					fmt.Printf("Failed to preprocess image: %v\n", err)
					return
				}
			}
			imgCh <- i
		}
	}()

	for i := range imgCh {
		file := path.Join(cfg.Dir, i.file)
		err := bimg.Write(file, i.data)
		if err != nil {
			fmt.Printf("Failed to write image: %v\n", err)
			continue
		}
	}
}
