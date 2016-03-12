package main

import "testing"

import "os"

func TestLock(t *testing.T) {
	createLock()
	defer unlock()

	_, err := os.Stat(LOCK_FILE_NAME)

	if err != nil && os.IsNotExist(err) {
		t.Fail()
	}
}

func TestUnlock(t *testing.T) {
	createLock()
	unlock()

	file, err := os.Stat(LOCK_FILE_NAME)

	if err == nil {
		t.Fail()
	}

	if !os.IsNotExist(err) {
		t.Fail()
	}

	if file != nil {
		t.Fail()
	}
}

func TestLockTwice(t *testing.T) {
	lock := createLock()

	if !lock {
		t.Error("Unable to create lock")
	}

	lock = createLock()

	if lock {
		t.Error("Able to lock for second time")
	}
}
