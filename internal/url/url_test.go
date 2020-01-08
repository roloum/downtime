package url

import (
	"testing"
)

//TestURL
func TestURL(t *testing.T) {
	t.Log("Testing URLs")

	if err := Check("wrongformat", false); err != ErrInvalidURL {
		t.Fatal("Failed checking invalid URL")
	}

	if err := Check("http://google.com", false); err != nil {
		t.Fatal("Failed checking working domain")
	}

	if err := Check("https://httpstat.us/400", false); err != ErrCouldNotCheck {
		t.Fatal("Failed checking 400 response")
	}

	if err := Check("https://httpstat.us/400", true); err != nil {
		t.Fatal("Failed checking domain only for 400 response ")
	}

	if err := Check("http://google.com/whatever", true); err != nil {
		t.Fatal("Failed checking domain only")
	}
}
