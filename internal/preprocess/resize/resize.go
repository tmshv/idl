package resize

import (
	"github.com/h2non/bimg"
)

type resize struct {
	Width  int
	Height int
}

func (e *resize) Run(input []byte) ([]byte, error) {
	i := bimg.NewImage(input)
	newImage, err := i.ResizeAndCrop(e.Width, e.Height)
	if err != nil {
		return nil, err
	}

	return newImage, nil
}

func New(width, height int) *resize {
	return &resize{
		Width:  width,
		Height: height,
	}
}
