# gohocr

**This is a fork of [github.com/stefanhengl/gohocr](https://github.com/stefanhengl/gohocr)**. The original repository was missing a lot of features
that I need for a project. This fork was created to be able to improve the quality of the parsed
data model. Feel free to collaborate.

gohocr parses Tesseract's *.hocr files to structs. It can be used as a library as well as from the command line.

# Getting Started

## Install

    go get github.com/fabiante/gohocr

## Use as Library

    import (
        "github.com/fabiante/gohocr
    )

    func main() {
        page, err := gohocr.Parse("path/to/hocr/file.hocr")
    }

The hOCR can either be provided as `string` (path to a hOCR file), as `*os.File` or `[]byte`. All of the following calls are equivalent
    
    page, err := gohocr.Parse("path/to/hocr/file.hocr")
    page, err := gohocr.Parse(os.File)
    page, err := gohocr.Parse([]byte)

`page` is a struct which can be marshalled into json

    pageJSON, err := json.Marshal(page)
    if err != nil {
        log.Fatal(err)
    }
	ioutil.WriteFile("page.json", pageJSON, 0644)

## Command Line Interface

The CLI will create a json in the same folder as the input file.

    go install github.com/fabiante/gohocr/..
    ghocr -f /path/to/hocr/file.hocr

