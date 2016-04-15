package main

import "testing"

import (
	"os"
)

func TestLock(t *testing.T) {
	unlock()
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
	defer unlock()
	lock := createLock()

	if !lock {
		t.Error("Unable to create lock")
	}

	lock = createLock()

	if lock {
		t.Error("Able to lock for second time")
	}
}

func createFakeSave() {
	os.Create(DEFAULT_SAVE)
}

func deleteFakeSave() {
	os.Remove(DEFAULT_SAVE)
}

func TestBackupCurrentSave(t *testing.T) {
	createFakeSave()
	MAX_KEPT_BACKUP_FILES = 0
	defer cleanup()
	defer deleteFakeSave()
	err := backupCurrentSave()

	if err != nil {
		t.Error("Unable to backup")
	}
}

func TestBackupCleanup(t *testing.T) {
	MAX_KEPT_BACKUP_FILES = 0

	cleanup()

	MAX_KEPT_BACKUP_FILES = DEFAULT_MAX_KEPT_BACKUP_FILES
}

func TestLoad(t *testing.T) {

	createFakeSave()
	defer deleteFakeSave()
	err := load()

	if err != nil {
		t.Error("Unable to load")
	}
}
