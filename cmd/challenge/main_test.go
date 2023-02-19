package main

import (
	"testing"
)

func TestReadJsonString(t *testing.T) {
	testString := `{"address":"0x38894302a6eabea6f2b29b508031d2ed75f0be22"}`
	expected := "0x38894302a6eabea6f2b29b508031d2ed75f0be22"
	addrs := readJsonString(testString)
	
	if expected != addrs[0] {
		t.Fatalf(`readJsonString(%s) did not return address "%s"`, testString, expected)
	}
}

func TestReadJsonStringMultipleAddrs(t *testing.T) {
	testString := `
		{"address":"0x5d88f42412dfd3bdef9d17ab0f77ea3b2077502f"}
		{"address":"0x38894302a6eabea6f2b29b508031d2ed75f0be22"}`
	
	expected1 := "0x5d88f42412dfd3bdef9d17ab0f77ea3b2077502f"
	expected2 := "0x38894302a6eabea6f2b29b508031d2ed75f0be22"
	addrs := readJsonString(testString)

	if expected1 != addrs[0] {
		t.Fatalf(`readJsonString(%s) did not return address "%s"`, testString, expected1)
	}

	if expected2 != addrs[1] {
		t.Fatalf(`readJsonString(%s) did not return address "%s"`, testString, expected2)
	}
}