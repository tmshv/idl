package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Record struct {
	URL  string
	File string
}

type CSV struct {
	FieldURL  string
	FieldFile string

	urlIdx  int
	fileIdx int
}

func (x *CSV) updateIndex(header []string) error {
	x.urlIdx = -1
	x.fileIdx = -1
	for i, field := range header {
		if field == x.FieldURL {
			x.urlIdx = i
		}
		if field == x.FieldFile {
			x.fileIdx = i
		}
	}
	if x.fileIdx == -1 {
		return fmt.Errorf("file field %s not found", x.FieldFile)
	}
	if x.urlIdx == -1 {
		return fmt.Errorf("url field %s not found", x.FieldURL)
	}
	return nil
}

func (x *CSV) makeRecord(row []string) Record {
	url := row[x.urlIdx]
	file := row[x.fileIdx]
	return Record{
		URL:  url,
		File: file,
	}
}

func (x *CSV) Read(ctx context.Context, path string, ch chan<- Record) error {
	defer close(ch)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	err = x.updateIndex(header)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			row, err := reader.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if len(row) == 0 {
				continue
			}

			rec := x.makeRecord(row)
			ch <- rec
		}
	}
}

func New(url, file string) *CSV {
	return &CSV{
		FieldURL:  url,
		FieldFile: file,
	}
}
