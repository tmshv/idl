package config

import (
	"time"

	"github.com/alecthomas/kong"
)

type Config struct {
	Input   string
	Dir     string
	Workers uint
	Reload  bool
	Resize  [2]int
	Timeout time.Duration
	Fields  struct {
		URL  string
		File string
	}
}

var CLI struct {
	Input     string `help:"Path to CSV file." type:"path"`
	Dir       string `help:"Path to directory to save images." type:"path"`
	Workers   uint   `help:"Number of workers." default:"5"`
	Timeout   string `help:"Number of workers." default:"10s"`
	Reload    bool   `help:"Force download image if it exists."`
	Resize    string `help:"Resize image before saving."`
	URLField  string `help:"Name of field of URL in CSV file."`
	FileField string `help:"Name of field of the name of file in CSV file."`

	// flag.StringVar(&source, "i", "", "path to file with urls")
	// flag.StringVar(&dest, "o", ".", "path to output folder")
	// flag.IntVar(&workers, "workers", 5, "number of workers")

	// flag.IntVar(&skip, "skip", 0, "pagination skip")
	// flag.IntVar(&limit, "limit", 0, "pagination limit")

	// flag.IntVar(&sample, "sample", 0, "download sample of urls")
	// flag.BoolVar(&reload, "reload", false, "skip loaded file or not")
	// flag.DurationVar(&timeout, "timeout", 10*time.Second, "http request timeout")
	// flag.StringVar(&urlField, "url-field", "url", "name of field of url in csv file")
	// flag.StringVar(&fileField, "file-field", "file", "name of field of file in csv file")
}

func Get() (Config, error) {
	kong.Parse(&CLI)
	resize, err := ParseResize(CLI.Resize, [2]int{0, 0}, [2]int{10000, 10000})
	if err != nil {
		return Config{}, err
	}
	timeout, err := ParseTimeout(CLI.Timeout)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Input:   CLI.Input,
		Dir:     CLI.Dir,
		Workers: CLI.Workers,
		Timeout: timeout,
		Reload:  CLI.Reload,
		Resize:  resize,
		Fields: struct {
			URL  string
			File string
		}{
			URL:  CLI.URLField,
			File: CLI.FileField,
		},
	}, nil
}
