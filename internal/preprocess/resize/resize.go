package resize

import (
	"bytes"

	"github.com/disintegration/imaging"
)

type resize struct {
	Width  int
	Height int
}

func (e *resize) Run(input []byte) ([]byte, error) {
	r := bytes.NewReader(input)
	src, err := imaging.Decode(r)
	if err != nil {
		return nil, err
	}

	img := imaging.Fill(src, e.Width, e.Height, imaging.Center, imaging.Lanczos)

	w := bytes.Buffer{}
	imaging.Encode(&w, img, imaging.JPEG)

	return w.Bytes(), nil
}

func New(width, height int) *resize {
	return &resize{
		Width:  width,
		Height: height,
	}
}
