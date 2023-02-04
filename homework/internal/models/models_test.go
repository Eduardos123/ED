package models

import (
	"testing"
)

func TestUserStructCreation(t *testing.T) {
	NormalUser := User{3, "Test", "18", nil}
	t.Log(NormalUser)
}
