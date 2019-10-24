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
