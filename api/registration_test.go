package main

import "testing"

func TestRegister(t *testing.T) {

	t.Log("something")
}

func TestRegister2(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	t.Log("something 2")

}
