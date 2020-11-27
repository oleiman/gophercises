package main

import (
	"testing"
)

var REF = "1234567890"
var WS = " 123 456 7890"
var PAREN = "(123)4567890"
var HYPH = "123-456-7890"
var ALL = "(123) 456-7890"

func TestNormalizeWS(t *testing.T) {
	result, err := normalize(WS)
	if err != nil {
		t.Error(err)
	}
	if result != REF {
		t.Errorf("Expected %s, got %s", REF, result)
	}
}

func TestNormalizePAREN(t *testing.T) {
	result, err := normalize(PAREN)
	if err != nil {
		t.Error(err)
	}
	if result != REF {
		t.Errorf("Expected %s, got %s", REF, result)
	}
}

func TestNormalizeHYPH(t *testing.T) {
	result, err := normalize(HYPH)
	if err != nil {
		t.Error(err)
	}
	if result != REF {
		t.Errorf("Expected %s, got %s", REF, result)
	}
}

func TestNormalizeALL(t *testing.T) {
	result, err := normalize(ALL)
	if err != nil {
		t.Error(err)
	}
	if result != REF {
		t.Errorf("Expected %s, got %s", REF, result)
	}
}
