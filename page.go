package gohocr

// Page represents one hocr file
type Page struct {
	Words []Word `xml:"body>div>div>p>span>span"`
}
