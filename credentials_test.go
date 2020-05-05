package JadeSocks

import (
	"testing"
)

func TestStaticCredentials(t *testing.T) {
	credentials := StaticCredentials{
		"foo": "bar",
		"baz": "",
	}

	if !credentials.Valid("foo", "bar") {
		t.Error()
	}

	if !credentials.Valid("baz", "") {
		t.Error()
	}

	if credentials.Valid("foo", "") {
		t.Error()
	}
}