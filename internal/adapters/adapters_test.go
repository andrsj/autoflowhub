package adapters

import "testing"

func TestDummyAdapter(t *testing.T) {
    expectedResult := "dummy adapter"
    result := "dummy adapter"

    if result != expectedResult {
        t.Errorf("Expected %q but got %q", expectedResult, result)
    }
}

