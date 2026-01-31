package image

import (
	"github.com/pRizz/dive/dive/filetree"
)

type Image struct {
	Request string
	Trees   []*filetree.FileTree
	Layers  []*Layer
}
