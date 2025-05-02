package preprocess

type Preprocessor interface {
	Run([]byte) ([]byte, error)
}
