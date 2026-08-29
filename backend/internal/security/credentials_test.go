package security

import "testing"

func TestCredentialsRoundTrip(t *testing.T) {
	ciphertext, err := EncryptCredentials(`{"username":"device","password":"secret"}`, "test-key")
	if err != nil {
		t.Fatal(err)
	}
	if ciphertext == `{"username":"device","password":"secret"}` {
		t.Fatal("credentials were not encrypted")
	}
	plaintext, err := DecryptCredentials(ciphertext, "test-key")
	if err != nil || plaintext != `{"username":"device","password":"secret"}` {
		t.Fatalf("unexpected plaintext %q: %v", plaintext, err)
	}
}
