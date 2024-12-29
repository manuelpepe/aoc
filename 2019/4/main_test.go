package main

import "testing"

func TestFulfills(t *testing.T) {
	if !fulfills(111111) {
		panic("expected 111111 to pass")
	}
	if fulfills(223450) {
		panic("expected 223450 to fail")
	}
	if fulfills(123789) {
		panic("expected 123789 to fail")
	}

}

func TestFulfills2(t *testing.T) {
	if !fulfills2(112233) {
		panic("expected 112233 to pass")
	}
	if fulfills2(123444) {
		panic("expected 123444 to fail")
	}
	if !fulfills2(111122) {
		panic("expected 111122 to fail")
	}
}
