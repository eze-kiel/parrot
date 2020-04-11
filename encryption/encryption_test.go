package encryption

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	cipher := []struct {
		data       string
		passphrase string
	}{
		{"Hello, friend", "abc"},
		{"Hello, friend", ""},
		{"", "123azerty"},
		{"/date", "abc"},
	}

	for _, test := range cipher {
		result := Encrypt([]byte(test.data), test.passphrase)
		result = Decrypt(result, test.passphrase)

		if string(result) != test.data {
			t.Errorf("Encrypt and decrypt were incorrect, got: %s, want: %s.", result, test.data)
		}
	}
}
