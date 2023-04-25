package notifications

import "testing"

func TestDummyNotification(t *testing.T) {
    expectedResult := "dummy notification"
    result := "dummy notification"

    if result != expectedResult {
        t.Errorf("Expected %q but got %q", expectedResult, result)
    }
}

