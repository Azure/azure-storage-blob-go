package azblob

import "io"

// Exports the method so we can use it in testing
// We don't want the underlying method in the public API, due to its specialist nature,
// so it remains private, accessed only by this _test file
func GetForceRetryFuncOrNil(rr io.Reader) func(){
	return rr.(*retryReader).getForceRetryFuncOrNil()
}