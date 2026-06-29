package file

import "errors"

var ErrFileLimitExceeded = errors.New("exceeds the number of files that can be uploaded")
