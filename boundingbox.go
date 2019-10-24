package gohocr

// Respective Spec: http://kba.cloud/hocr-spec/1.2/#bbox
type BoundingBox struct {
	X0, X1, Y0, Y1 uint
}
