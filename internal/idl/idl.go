package idl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/tmshv/idl/internal/config"
	"github.com/tmshv/idl/internal/csv"
	"github.com/tmshv/idl/internal/dl"
	"github.com/tmshv/idl/internal/preprocess"
	"github.com/tmshv/idl/internal/preprocess/compose"
	"github.com/tmshv/idl/internal/preprocess/empty"
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

func (idl *IDL) count(ctx context.Context, cfg config.Config) int64 {
	ch := make(chan csv.Record)
	errCh := make(chan struct{})
	go func() {
		reader := csv.New(cfg.Fields.URL, cfg.Fields.File)
		err := reader.Read(ctx, cfg.Input, ch)
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			errCh <- struct{}{}
		}
	}()

	var count int64
	for {
		select {
		case <-ctx.Done():
			return 0
		case <-errCh:
			return 0
		case _, ok := <-ch:
			if !ok {
				return count
			}
			count += 1
		}
	}
}

func (idl *IDL) Run(ctx context.Context, cfg config.Config) error {
	total := idl.count(ctx, cfg)

	err := os.MkdirAll(cfg.Dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to make dir: %v\n", err)
	}

	recordCh := make(chan csv.Record)
	go (func() {
		reader := csv.New(cfg.Fields.URL, cfg.Fields.File)
		err := reader.Read(ctx, cfg.Input, recordCh)
		if errors.Is(err, context.Canceled) {
			return
		}
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

		prep := idl.makePreprocessors(&cfg)

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
	bar := progressbar.Default(total)
	for i := range imgCh {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := os.WriteFile(i.file, i.data, 0644)
			if err != nil {
				fmt.Printf("Failed to save file: %v\n", err)
				return
			}
			err = bar.Add(1)
			if err != nil {
				fmt.Printf("Failed update progressbar: %v\n", err)
				return
			}
		}()
	}
	wg.Wait()

	return nil
}

func (idl *IDL) makePreprocessors(cfg *config.Config) preprocess.Preprocessor {
	if cfg.Resize == [2]int{0, 0} {
		fmt.Printf("Image resizing is disabled")
		return empty.New()
	}
	return compose.New(
		resize.New(cfg.Resize[0], cfg.Resize[1]),
	)
}

func New(cpu int) *IDL {
	return &IDL{
		cpu: cpu,
	}
}
