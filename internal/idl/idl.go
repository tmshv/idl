package idl

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

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

type IDL struct {
	cpu int
}

func (idl *IDL) Run(ctx context.Context, cfg config.Config) error {
	recordCh := make(chan csv.Record, idl.cpu)
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
		for range cfg.Workers {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for rec := range recordCh {
					file := path.Join(cfg.Dir, rec.File)
					if !cfg.Reload && utils.FileExists(file) {
						fmt.Printf("Skip: %s\n", file)
						continue
					}

					body, err := dl.Download(rec.URL)
					if err != nil {
						fmt.Printf("Failed to download: %v\n", err)
						continue
					}

					dlCh <- image{file: file, data: body}
				}
			}()
		}

		wg.Wait()
	}()

	imgCh := make(chan image, idl.cpu)
	go func() {
		defer close(imgCh)

		prep := compose.New(
			resize.New(cfg.Resize[0], cfg.Resize[1]),
		)

		wg := sync.WaitGroup{}
		for range idl.cpu {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := range dlCh {
					data, err := prep.Run(i.data)
					if err != nil {
						fmt.Printf("Failed to preprocess image: %v\n", err)
						return
					}
					i.data = data
					imgCh <- i
				}
			}()
		}
		wg.Wait()
	}()

	wg := sync.WaitGroup{}
	for i := range imgCh {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := os.WriteFile(i.file, i.data, 0644)
			if err != nil {
				fmt.Printf("Failed to save file: %v\n", err)
				return
			}
			fmt.Printf("OK: %v\n", i.file)
		}()
	}
	wg.Wait()

	return nil
}

func New(cpu int) *IDL {
	return &IDL{
		cpu: cpu,
	}
}
