package payload

import "github.com/pRizz/dive/dive/image"

type Explore struct {
	Analysis image.Analysis
	Content  image.ContentReader
}
