package resize

import (
	"bytes"
	"errors"
	"image"
	"io"
	"os"
	"path"
	"testing"

	"github.com/disintegration/imaging"
)

func TestResize_Run(t *testing.T) {
	type args struct {
		input  []byte
		width  int
		height int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				input:  readImageFile("example.jpg"),
				width:  100,
				height: 100,
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			args: args{
				input:  []byte("invalid image data"),
				width:  100,
				height: 100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resize := New(tt.args.width, tt.args.height)
			_, err := resize.Run(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && err == nil {
				// Perform basic validation on the output image.  Check dimensions are not zero.
				outputBytes, _ := resize.Run(tt.args.input)
				img, err := imaging.Decode(bytes.NewReader(outputBytes))
				if err != nil {
					t.Fatalf("failed to decode output image: %v", err)
				}

				if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
					t.Errorf("output image dimensions are zero")
				}
			}
		})
	}
}

func readImageFile(filename string) []byte {
	path := path.Join("../../../testdata", filename)
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err) // Or handle the error more gracefully in a real test
	}
	return f
}

// Mock imaging.Decode to return an error for testing purposes
func mockDecode(r io.Reader) (image.Image, error) {
	return nil, errors.New("mock decode error")
}
