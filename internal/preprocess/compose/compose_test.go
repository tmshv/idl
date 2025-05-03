package compose

import (
	"errors"
	"testing"
)

// Mock Preprocessor for testing
type mock struct {
	fn func([]byte) ([]byte, error)
}

func (m *mock) Run(input []byte) ([]byte, error) {
	return m.fn(input)
}

func TestCompose_Run(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// Arrange
		input := []byte("initial input")
		prep1 := &mock{fn: func(input []byte) ([]byte, error) {
			return []byte("prep1 applied: " + string(input)), nil
		}}
		prep2 := &mock{fn: func(input []byte) ([]byte, error) {
			return []byte("prep2 applied: " + string(input)), nil
		}}
		composer := New(prep1, prep2)

		// Act
		result, err := composer.Run(input)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []byte("prep2 applied: prep1 applied: initial input")
		if string(result) != string(expected) {
			t.Errorf("expected: %s, got: %s", string(expected), string(result))
		}
	})

	t.Run("error in first preprocessor", func(t *testing.T) {
		// Arrange
		input := []byte("initial input")
		prep1 := &mock{fn: func(input []byte) ([]byte, error) {
			return nil, errors.New("prep1 failed")
		}}
		prep2 := &mock{fn: func(input []byte) ([]byte, error) {
			return []byte("prep2 applied: " + string(input)), nil
		}}
		composer := New(prep1, prep2)

		// Act
		result, err := composer.Run(input)

		// Assert
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "prep1 failed" {
			t.Errorf("expected error message 'prep1 failed', got '%s'", err.Error())
		}
		if result != nil {
			t.Errorf("expected nil result, got '%s'", string(result))
		}
	})

	t.Run("error in second preprocessor", func(t *testing.T) {
		// Arrange
		input := []byte("initial input")
		prep1 := &mock{fn: func(input []byte) ([]byte, error) {
			return []byte("prep1 applied: " + string(input)), nil
		}}
		prep2 := &mock{fn: func(input []byte) ([]byte, error) {
			return nil, errors.New("prep2 failed")
		}}
		composer := New(prep1, prep2)

		// Act
		result, err := composer.Run(input)

		// Assert
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "prep2 failed" {
			t.Errorf("expected error message 'prep2 failed', got '%s'", err.Error())
		}
		if result != nil {
			t.Errorf("expected nil result, got '%s'", string(result))
		}
	})

	t.Run("empty composer", func(t *testing.T) {
		// Arrange
		input := []byte("initial input")
		composer := New()

		// Act
		result, err := composer.Run(input)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != string(input) {
			t.Errorf("expected input to be returned unchanged, got: %s", string(result))
		}
	})

	t.Run("preprocessor returns empty byte slice", func(t *testing.T) {
		// Arrange
		input := []byte("initial input")
		prep1 := &mock{fn: func(input []byte) ([]byte, error) {
			return []byte(""), nil
		}}
		prep2 := &mock{fn: func(input []byte) ([]byte, error) {
			return []byte("prep2 applied: " + string(input)), nil
		}}
		composer := New(prep1, prep2)

		// Act
		result, err := composer.Run(input)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []byte("prep2 applied: ")
		if string(result) != string(expected) {
			t.Errorf("expected: %s, got: %s", string(expected), string(result))
		}
	})
}
