package image

import "errors"

var ErrImageLimitExceeded = errors.New("exceeds the number of images that can be uploaded")
