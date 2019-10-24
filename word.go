package gohocr

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

// Populates all fields that get their value from the Title property
func (word *Word) populateTitleFields() {
	word.Confidence = getConfidenceFromTitle(word.Title)
	word.BoundingBox = getBoundingBoxFromTitle(word.Title)
}
