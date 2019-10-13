package gohocr

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	_, err := Parse("./test/test.hocr")
	if err != nil {
		t.Error("Parse failed")
	}
}

func TestWordFields(t *testing.T) {
	page, err := Parse("./test/test.hocr")
	if err != nil {
		t.Error("Parse failed")
	}
	word := page.Words[3]
	if word.Content != "Lorem" {
		t.Errorf("Have: %s -- Expected: Lorem", word.Content)
	}
	if word.ID != "word_1_4" {
		t.Errorf("Have: %s -- Expected: word_1_4", word.ID)
	}
	if word.Lang != "eng" {
		t.Errorf("Have: %s -- Expected: eng", word.Lang)
	}
	if word.Title != "bbox 299 432 422 465; x_wconf 97.01" {
		t.Errorf("Have: %s -- Expected: bbox 299 432 422 465; x_wconf 97.01", word.Title)
	}
	if word.Direction != "ltr" {
		t.Errorf("Have: %s -- Expected: ltr", word.Direction)
	}
}

func TestWrongPath(t *testing.T) {
	_, err := Parse("./foo/test.hocr")
	if err == nil {
		t.Error("Parse should have failed with wrong path")
	}
}

func TestNotAHOCR(t *testing.T) {
	_, err := Parse("./test/notahocr.hocr")
	if err == nil {
		t.Error("Parse should have failed with invalid hocr")
	}
}

func TestOSFilePoiner(t *testing.T) {
	file, err := os.Open("./test/test.hocr")
	if err != nil {
		t.Error("os.Open should return valid pointer")
	}
	_, err = Parse(file)
	if err != nil {
		t.Error("Parse failed")
	}
}

func TestByteSlice(t *testing.T) {
	xmlFile, err := os.Open("./test/test.hocr")
	if err != nil {
		t.Error("os.Open should return valid pointer")
	}
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		t.Error("ioutil.ReadAll should return a valid byte slice")
	}

	_, err = Parse(byteValue)
	if err != nil {
		t.Error("Parse failed")
	}
}

func TestInvalidInput(t *testing.T) {
	_, err := Parse(1)
	if err == nil {
		t.Error("Parse should have returned error with invalid input")
	}
}
