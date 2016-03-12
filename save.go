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
const MAX_KEPT_BACKUP_FILES = 60

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
		p(fmt.Sprintf("Failed %q", err))
	} else {
		p(fmt.Sprintf("Saved %s", newFilename))
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
				p("Removed " + f)
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

const SAVE_DIRECTORY = "./"

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

func load() {

	err := backupCurrentSave()

	if err != nil {
		printF("Error creating backup %q", err)
	} else {
		latestSaveFileName, err := loadLatestSave()

		if err != nil {
			printF("Error loading latest save %q", err)
		} else if latestSaveFileName == "" {
			p("No saves located")
		} else {
			printF("Loaded %s\n", latestSaveFileName)
			err = os.Remove(latestSaveFileName)
		}
	}
}

func undo() {

}

func printF(msg string, args ...interface{}) {
	p(format(msg, args))
}

func format(msg string, args []interface{}) string {
	return fmt.Sprintf(msg, args)
}

var info os.FileInfo

// Checks save file every once in a while
// If the file was changed, backs it up
func watchSave() {

	defer p("File watcher exiting")

	info, _ = os.Stat(DEFAULT_SAVE)

	for {
		time.Sleep(time.Second * CHECK_FREQUENCY)

		currentInfo, err := os.Stat(DEFAULT_SAVE)

		if (err != nil) {
			continue
		}

		timeDiff := currentInfo.ModTime().Sub(info.ModTime())

		if timeDiff.Seconds() > 0 {
			p(fmt.Sprintf("Save changed, diff %f", timeDiff.Seconds()))
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

// Helper method for output
func p(msg string) {
	fmt.Println(msg)
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
	os.Remove(LOCK_FILE_NAME)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	go watchSave()
	lock := createLock()

	if lock {
		defer unlock()

		exit := false
		for !exit {
			latestSaveName, _ := getLatestSave()
			p("What would you like to do?")
			p(fmt.Sprintf("[%s] for save, [%s] for load, [%s] to undo load, [%s] for exit", SAVE_CHAR, LOAD_CHAR, UNDO_CHAR, EXIT_CHAR))
			p(fmt.Sprintf("Last save is %s", latestSaveName))
			p("Hit enter to confirm")
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
				p("Exiting")
				exit = true
			default:
				p("Unknown command: " + char)
			}
		}
	} else {
		p("Lock file found")
		p("Either you have another instance of this program already running, or the progam didn't quit correctly last time")
		p("If the later is the case, please remove the lock file manually")
	}

	fmt.Println("Bye!")
}
