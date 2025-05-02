package empty

type EmptyPreprocessor struct{}

func (e *EmptyPreprocessor) Run(input []byte) ([]byte, error) {
	return input, nil
}

func New() *EmptyPreprocessor {
	return &EmptyPreprocessor{}
}
