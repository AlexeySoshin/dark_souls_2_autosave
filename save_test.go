package main
import "testing"

import "os"

func TestLock(t *testing.T) {
	createLock()

	_, err := os.Stat(LOCK_FILE_NAME)

	if (err != nil && os.IsNotExist(err)) {
		t.Fail()
	}
}

func TestUnlock(t *testing.T) {
	unlock()

	_, err := os.Stat(LOCK_FILE_NAME)

	if (err == nil) {
		t.Fail()
	}
}