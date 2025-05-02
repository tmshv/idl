package compose

import "github.com/tmshv/idl/internal/preprocess"

type compose struct {
	items []preprocess.Preprocessor
}

func (c *compose) Run(input []byte) ([]byte, error) {
	for _, item := range c.items {
		var err error
		input, err = item.Run(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

func New(preprocessors ...preprocess.Preprocessor) *compose {
	return &compose{items: preprocessors}
}
