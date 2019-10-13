package gohocr

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	_, err := Parse("./test/test.hocr")
	if err != nil {
		t.Error("Parse failed")
	}
}

func TestDefaultWordFields(t *testing.T) {
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

func TestPopulatedWordFields(t *testing.T) {
	page, err := Parse("./test/test.hocr")
	if err != nil {
		t.Error("Parse failed")
	}
	word := page.Words[3]

	// Title="bbox 299 432 422 465; x_wconf 97.01"

	if word.Confidence != 97.01 {
		t.Errorf("Have: %f -- Expected: 97.01", word.Confidence)
	}

	expectedBoundingBox := BoundingBox{
		X0: 299,
		X1: 422,
		Y0: 432,
		Y1: 465,
	}
	if !reflect.DeepEqual(word.BoundingBox, expectedBoundingBox) {
		t.Errorf("Have: %v -- Expected: %v", word.BoundingBox, expectedBoundingBox)
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

func Test_getConfidenceFromTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "should get int",
			args: args{title: "bbox 299 432 422 465; x_wconf 97"},
			want: 97,
		},
		{
			name: "should get float",
			args: args{title: "bbox 299 432 422 465; x_wconf 97.01"},
			want: 97.01,
		},
		{
			name: "should return 0 if missing",
			args: args{title: "bbox 299 432 422 465"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfidenceFromTitle(tt.args.title); got != tt.want {
				t.Errorf("getConfidenceFromTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBoundingBoxFromTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		wantBox BoundingBox
	}{
		{
			name: "should use bbox",
			args: args{title: "bbox 299 432 422 465; x_wconf 97.01"},
			wantBox: BoundingBox{
				X0: 299,
				X1: 422,
				Y0: 432,
				Y1: 465,
			},
		},
		{
			name: "should have zero coordinates if not in title",
			args: args{title: "x_wconf 97.01"},
			wantBox: BoundingBox{
				X0: 0,
				X1: 0,
				Y0: 0,
				Y1: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBox := getBoundingBoxFromTitle(tt.args.title); !reflect.DeepEqual(gotBox, tt.wantBox) {
				t.Errorf("getBoundingBoxFromTitle() = %v, want %v", gotBox, tt.wantBox)
			}
		})
	}
}
