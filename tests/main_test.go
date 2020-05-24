package main_test

import (
  "testing"
)

func TestMain(t *testing.T) {
  if 2*2 != 4 {
    t.Errorf("FALED")
  }
}
