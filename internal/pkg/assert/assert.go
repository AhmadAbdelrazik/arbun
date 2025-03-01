package assert

import (
	"errors"
	"testing"

	"github.com/Rhymond/go-money"
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

func True(t *testing.T, condition bool) {
	t.Helper()

	if condition != true {
		t.Fatal("got:  false\nwant: true")
	}
}

func False(t *testing.T, condition bool) {
	t.Helper()

	if condition != false {
		t.Fatal("got:  true\nwant: false")
	}
}

func MoneyEqual(t *testing.T, got, want *money.Money) {
	t.Helper()

	ok, err := got.Equals(want)
	Nil(t, err)
	True(t, ok)
}
