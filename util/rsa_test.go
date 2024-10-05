package util

import (
	"testing"
)

func TestSigner_SignLease(t *testing.T) {
	// jGydenKNH1HWKHo9lFm3lKVCSbBAIYBob/QbvnhdW9FnkHmM5RPvSrd62gqH8rasP1CL++cgTGG0qthFXaUP+m9y42ulmvadzC2os/ZknEi4Xx2fTrH4DwtACsp9o7HkI6BcWvIuBL/8eNyrcf50T14fMWpGZ9WJHGJwP7/qy4w=
	signer := NewSigner()
	if lease, err := signer.SignLease("bTqakeUt614==",
		"wdpQy7vZRgTWQIXjqA1VAAbSBLp9PsmQMHSMj6c8w6E==",
		false, "validFrom", "validUntil"); err == nil {
		t.Log(lease)
	} else {
		t.Fatal(err)
	}
}
