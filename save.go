package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const DEFAULT_SAVE = "DS2SOFS0000.sl2"

const SAVE_SUFFIX = ".sav"
const BACKUP_SUFFIX = ".bak"
const SAVE_FORMAT = "%02d%02d_%02d%02d%02d"

const CHECK_FREQUENCY = 1
const SAVE_FREQUENCY = 10

const MAX_KEPT_SAVE_FILES = 30
const DEFAULT_MAX_KEPT_BACKUP_FILES = 60

const SAVE_DIRECTORY = "./"

var MAX_KEPT_BACKUP_FILES = DEFAULT_MAX_KEPT_BACKUP_FILES

const MAX_SAVES = 100
const MAX_BACKUPS = 100

func copyFiles(oldFilename string, newFilename string) error {
	src, err := os.Open(oldFilename)
	defer src.Close()

	dst, err := os.Create(newFilename)
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}

func save() {
	t := time.Now()

	newFilename := fmt.Sprintf(SAVE_FORMAT+SAVE_SUFFIX,
		t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	err := copyFiles(DEFAULT_SAVE, newFilename)

	if err != nil {
		message(fmt.Sprintf("Failed %q", err))
	} else {
		message(fmt.Sprintf("Saved %s", newFilename))
	}

	cleanup()
}

func isSaveFile(fileName string) bool {
	return strings.Contains(fileName, SAVE_SUFFIX) && fileName != DEFAULT_SAVE
}

func isBackupFile(fileName string) bool {
	return strings.Contains(fileName, BACKUP_SUFFIX)
}

func deleteOldFiles(fileNames []string, numFilesToKeep int) {
	if len(fileNames) > numFilesToKeep {
		filesToDelete := fileNames[0 : len(fileNames)-MAX_KEPT_SAVE_FILES-1]
		for _, f := range filesToDelete {
			err := os.Remove(f)
			if err == nil {
				debug("Removed " + f)
			}
		}
	}
}

func cleanup() {
	files, _ := ioutil.ReadDir(SAVE_DIRECTORY)

	allSaves := make([]string, MAX_SAVES)
	allBackups := make([]string, MAX_BACKUPS)
	for _, f := range files {
		currentFileName := f.Name()
		if isSaveFile(currentFileName) {
			allSaves = append(allSaves, currentFileName)
		} else if isBackupFile(currentFileName) {
			allBackups = append(allBackups, currentFileName)
		}
	}

	deleteOldFiles(allSaves, MAX_KEPT_SAVE_FILES)
	deleteOldFiles(allBackups, MAX_KEPT_BACKUP_FILES)
}

func backupCurrentSave() error {

	t := time.Now()

	newFilename := fmt.Sprintf(SAVE_FORMAT+BACKUP_SUFFIX,
		t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return copyFiles(DEFAULT_SAVE, newFilename)
}

func getLatestSave() (fileName string, file os.FileInfo) {

	latestSaveFileName := ""
	var latestSaveFile os.FileInfo

	files, _ := ioutil.ReadDir(SAVE_DIRECTORY)

	for _, f := range files {
		currentFileName := f.Name()

		if isSaveFile(currentFileName) {
			latestSaveFileName = currentFileName
			latestSaveFile = f
		}
	}

	return latestSaveFileName, latestSaveFile
}

func loadLatestSave() (latestSaveFileName string, err error) {
	latestSaveFileName, _ = getLatestSave()
	err = os.Remove(DEFAULT_SAVE)

	copyFiles(latestSaveFileName, DEFAULT_SAVE)
	info, _ = os.Stat(DEFAULT_SAVE)

	return latestSaveFileName, err
}

// Attempts to create a backup, then load latest save file
func load() error {

	err := backupCurrentSave()

	if err != nil {
		warning(fmt.Sprintf("Error creating backup %q", err))
	} else {
		latestSaveFileName, err := loadLatestSave()

		if err != nil {
			warning(fmt.Sprintf("Error loading latest save %q", err))
		} else if latestSaveFileName == "" {
			message("No saves located")
		} else {
			message(fmt.Sprintf("Loaded %s\n", latestSaveFileName))
			err = os.Remove(latestSaveFileName)
		}
	}

	return err
}

func undo() {

}

var info os.FileInfo

// Checks save file every once in a while
// If the file was changed, backs it up
func watchSave() {

	defer debug("File watcher exiting")

	info, _ = os.Stat(DEFAULT_SAVE)

	for {
		time.Sleep(time.Second * CHECK_FREQUENCY)

		currentInfo, err := os.Stat(DEFAULT_SAVE)

		if err != nil {
			continue
		}

		timeDiff := currentInfo.ModTime().Sub(info.ModTime())

		if timeDiff.Seconds() > 0 {

			debug(fmt.Sprintf("Save changed, diff %f", timeDiff.Seconds()))
			if timeDiff > SAVE_FREQUENCY {
				info = currentInfo
				save()
			}
		}
	}
}

const SAVE_CHAR = "s"
const LOAD_CHAR = "l"
const EXIT_CHAR = "x"
const UNDO_CHAR = "u"

const LOG_WARNING = 2
const LOG_INFO = LOG_WARNING - 1
const LOG_DEBUG = LOG_INFO - 1

var logLevel = LOG_INFO

func warning(msg string) {
	if logLevel <= LOG_WARNING {
		fmt.Println("Error: " + msg)
	}
}

// Helper method for output
func message(msg string) {
	if logLevel <= LOG_INFO {
		fmt.Println(msg)
	}
}

func debug(msg string) {
	if logLevel <= LOG_DEBUG {
		fmt.Println("DEBUG: " + msg)
	}
}

const LOCK_FILE_NAME = "saves.lock"

func createLock() bool {

	_, err := os.Stat(LOCK_FILE_NAME)
	if !os.IsNotExist(err) {
		fmt.Println("Lock file already exists")
		return false
	}

	_, err = os.Create(LOCK_FILE_NAME)

	if err != nil {
		fmt.Println("Unable to create lock file: ", err)
		return false
	}

	return true
}

func unlock() {
	err := os.Remove(LOCK_FILE_NAME)

	if err != nil {
		warning("Unable to unlock the file")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	go watchSave()
	lock := createLock()

	if lock {
		defer unlock()

		exit := false
		message("What would you like to do?")
		for !exit {
			latestSaveName, _ := getLatestSave()
			message(fmt.Sprintf("Last save is %s", latestSaveName))
			message(fmt.Sprintf("[%s] for save, [%s] for load, [%s] to undo load, [%s] for exit", SAVE_CHAR, LOAD_CHAR, UNDO_CHAR, EXIT_CHAR))
			message("Hit enter to confirm")
			input, _ := reader.ReadString('\n')

			char := input[0:1]

			switch char {
			case SAVE_CHAR:
				save()
			case LOAD_CHAR:
				load()
			case UNDO_CHAR:
				undo()
			case EXIT_CHAR:
				message("Exiting")
				exit = true
			default:
				message("Unknown command: " + char)
			}
		}
	} else {
		message("Lock file found")
		message("Either you have another instance of this program already running, or the progam didn't quit correctly last time")
		message("If the later is the case, please remove the lock file manually")
	}

	fmt.Println("Bye!")
}
