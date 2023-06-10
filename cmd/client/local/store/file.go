package store

import "time"

func FetchUser() (string, bool) {
	time.Sleep(time.Second)
	return "lordvidex", true
}
