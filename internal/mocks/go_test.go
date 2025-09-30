package mocks

import (
	"testing"

	. "github.com/switch/birb"
)

func TestGreeter(t *testing.T) {
	mockGreeter := NewMockGreeter(t)
	WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("Hello, John")

	actual := mockGreeter.Greet()

	if actual != "Hello, John" {
		t.Errorf("Expected 'Hello, John', got '%s'", actual)
	}
	Verify(mockGreeter, Once()).CALLED_Greet()
	VerifyNoOtherInteractions(mockGreeter)
}
