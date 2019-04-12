package file

import "runtime"

func GetCurrentFilePath() (string, bool) {
	_, filename, _, ok := runtime.Caller(1)
	return filename, ok
}