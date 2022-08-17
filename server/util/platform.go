package util

import "runtime"

func ExtBasedOnPlatform() string {
	if runtime.GOOS == "linux" {
		return ".tar.gz"
	} else {
		return ".zip"
	}
}
