/*
Package gohocr parses hocr files.
*/
package gohocr

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

// Page represents one hocr file
type Page struct {
	Words []Word `xml:"body>div>div>p>span>span"`
}

// Word represents a single word within an hocr file
// <span class='ocrx_word' id='word_1_11' title='bbox 572 568 641 684; x_wconf 57' lang='eng' dir='ltr'>Foo</span>
type Word struct {
	Lang        string      `xml:"lang,attr"`
	Direction   string      `xml:"dir,attr"`
	Title       string      `xml:"title,attr"`
	ID          string      `xml:"id,attr"`
	Class       string      `xml:"class,attr"`
	Content     string      `xml:",innerxml"`
	BoundingBox BoundingBox ``
	Confidence  float64     ``
}

// Respective Spec: http://kba.cloud/hocr-spec/1.2/#bbox
type BoundingBox struct {
	X0, X1, Y0, Y1 uint
}

// Populates all fields that get their value from the Title property
func (word *Word) populateTitleFields() {
	word.Confidence = getConfidenceFromTitle(word.Title)
	word.BoundingBox = getBoundingBoxFromTitle(word.Title)
}

func getConfidenceFromTitle(title string) float64 {
	if matches := regexConfidence.FindStringSubmatch(title); matches != nil {
		if confidence, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return confidence
		}
	}

	return 0
}

func getBoundingBoxFromTitle(title string) (box BoundingBox) {
	box = BoundingBox{}

	if matches := regexBoundingBox.FindStringSubmatch(title); matches != nil {
		if x0, err := strconv.ParseUint(matches[1], 10, 32); err == nil {
			box.X0 = uint(x0)
		}
		if y0, err := strconv.ParseUint(matches[2], 10, 32); err == nil {
			box.Y0 = uint(y0)
		}
		if x1, err := strconv.ParseUint(matches[3], 10, 32); err == nil {
			box.X1 = uint(x1)
		}
		if y1, err := strconv.ParseUint(matches[4], 10, 32); err == nil {
			box.Y1 = uint(y1)
		}
	}

	return
}

var (
	// Reference: http://kba.cloud/hocr-spec/1.2/#x_wconf - Example: x_wconf 97
	regexConfidence *regexp.Regexp
	// Reference: http://kba.cloud/hocr-spec/1.2/#bbox - Example: bbox 607 552 733 587
	regexBoundingBox *regexp.Regexp
)

func init() {
	regexConfidence = regexp.MustCompile(`x_wconf (\d*\.?\d*)`)
	regexBoundingBox = regexp.MustCompile(`bbox (\d*) (\d*) (\d*) (\d*)`)
}

// Parse takes either a string, *os.File, or []byte and returns a Page object
func Parse(value interface{}) (Page, error) {
	var byteValue []byte
	var err error
	switch str := value.(type) {
	case string:
		xmlFile, err := os.Open(str)
		if err != nil {
			return Page{}, err
		}
		defer xmlFile.Close()
		byteValue, err = ioutil.ReadAll(xmlFile)
		if err != nil {
			return Page{}, err
		}
	case *os.File:
		byteValue, err = ioutil.ReadAll(str)
		if err != nil {
			return Page{}, err
		}
	case []byte:
		byteValue = str
	default:
		return Page{}, errors.New("Invalid input for Parse. Submit either a string, *os.File, or []byte")
	}

	var page Page
	err = xml.Unmarshal(byteValue, &page)
	if err != nil {
		return Page{}, err
	}

	for i, _ := range page.Words {
		page.Words[i].populateTitleFields()
	}

	return page, nil
}
