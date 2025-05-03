package empty

import (
	"bytes"
	"testing"
)

func TestEmptyPreprocessor_Run(t *testing.T) {
	t.Run("should return the input unchanged", func(t *testing.T) {
		// Arrange
		input := []byte("test input")
		preprocessor := New()

		// Act
		output, err := preprocessor.Run(input)

		// Assert
		if err != nil {
			t.Fatalf("Run should not return an error, got: %v", err)
		}

		if !bytes.Equal(input, output) {
			t.Errorf("Run should return the input unchanged, got: %s, expected: %s", string(output), string(input))
		}
	})

	t.Run("should return an empty slice unchanged", func(t *testing.T) {
		// Arrange
		input := []byte{}
		preprocessor := New()

		// Act
		output, err := preprocessor.Run(input)

		// Assert
		if err != nil {
			t.Fatalf("Run should not return an error, got: %v", err)
		}

		if !bytes.Equal(input, output) {
			t.Errorf("Run should return the input unchanged, got: %s, expected: %s", string(output), string(input))
		}
	})

	t.Run("should handle nil input", func(t *testing.T) {
		// Arrange
		var input []byte // nil slice
		preprocessor := New()

		// Act
		output, err := preprocessor.Run(input)

		// Assert
		if err != nil {
			t.Fatalf("Run should not return an error, got: %v", err)
		}

		if !bytes.Equal(input, output) {
			t.Errorf("Run should return the input unchanged, got: %s, expected: %s", string(output), string(input))
		}
	})
}
