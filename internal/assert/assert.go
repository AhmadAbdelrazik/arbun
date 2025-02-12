package assert

import (
	"errors"
	"testing"
)

func Nil(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func Err(t *testing.T, got, want error) {
	t.Helper()

	if errors.Is(want, got) {
		t.Fatalf("got:  %v\nwant: %v", got, want)
	}
}

func Equal[T comparable](t *testing.T, got T, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("got:  %v\nwant: %v", got, want)
	}
}
