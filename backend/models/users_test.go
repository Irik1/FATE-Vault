package models

import "testing"

func TestUsers_SetPasswordAndCheckPassword_Success(t *testing.T) {
	u := &Users{}

	plain := "my-secure-password"
	if err := u.SetPassword(plain); err != nil {
		t.Fatalf("SetPassword returned error: %v", err)
	}

	if u.HashedPassword == "" {
		t.Fatalf("expected HashedPassword to be set, got empty string")
	}

	if !u.CheckPassword(plain) {
		t.Fatalf("CheckPassword should return true for correct password")
	}
}

func TestUsers_CheckPassword_WrongPassword(t *testing.T) {
	u := &Users{}

	plain := "my-secure-password"
	if err := u.SetPassword(plain); err != nil {
		t.Fatalf("SetPassword returned error: %v", err)
	}

	if u.CheckPassword("wrong-password") {
		t.Fatalf("CheckPassword should return false for incorrect password")
	}
}

func TestUsers_CheckPassword_EmptyHash(t *testing.T) {
	u := &Users{}

	if u.CheckPassword("anything") {
		t.Fatalf("CheckPassword should return false when HashedPassword is empty")
	}
}
