package main

import (
	"fmt"
	"log"
	"flag"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
)

func main() {
	wordPtr := flag.String("path", "", "The path to the folder with video and srt files.")
	searchWordPtr := flag.String(".", "", "Provide a unique word to identify the files by.")
	videoFileExtensionPtr := flag.String("videoFileExtension", ".mkv", "Provide the video file extension. Defaults to .mkv.")
	confirmationRequiredPtr := flag.String("confirmation", "", "Disable the confirmation before every rename.")

	flag.Parse()

	pathToFolder :=  *wordPtr
	searchWord :=  *searchWordPtr
	videoFileExtension :=  *videoFileExtensionPtr
	confirmationRequired := *confirmationRequiredPtr

	if pathToFolder == "" {
		fmt.Println("Please define a -path")
		return
	}

	confirmationRequiredBool := false
	if confirmationRequired == "" {
		confirmationRequiredBool = true
	}

	fmt.Println("Looking for files in " + pathToFolder)

	if searchWord != "" {
		fmt.Println("Looking for files that include the expression: " + searchWord)
	} else {
		fmt.Println("Looking for all files")
	}

	if searchWord != "" {
		fmt.Println("Looking for video files with file extension: " + videoFileExtension)
	} else {
		fmt.Println("Looking for all video files")
	}

	files, err := ioutil.ReadDir(pathToFolder)
	if err != nil {
		log.Fatal(err)
	}

	// Find matching .mkv (video) files.
	for _, f := range files {

		// Skip all files that did not match the search string.
		if !checkForFileName(searchWord, f.Name()) {
			continue
		}

		filenameExpression, _ := regexp.Compile("(.*)(\\.\\w*)$")
		videoFilename := filenameExpression.FindStringSubmatch(f.Name())

		// Skip all files without the defined video extension.
		if !checkForFileExtension(videoFileExtension, videoFilename[2]) {
			continue
		}

		match, _ := regexp.MatchString("(?i)(S\\w{2,2}E\\w{2,2})", f.Name())

		// Skip all files that do not contains "SYYEYY" (e.g. S02E12)
		if !match {
			continue
		}

		r, _ := regexp.Compile("(?i)S(\\w{2,2})E(\\w{2,2})")
		result := r.FindStringSubmatch(f.Name())
		result = append(result[:0], result[1:]...)

		// Find matching .srt files.
		for _, g := range files {
			// Skip all files that did not match the search string.
			if !checkForFileName(searchWord, g.Name()) {
				continue
			}

			filenameExpression, _ := regexp.Compile("(.*)(\\.\\w*)$")
			srtFilename := filenameExpression.FindStringSubmatch(g.Name())

			// Skip all files with the defined video extension to only get subtitle files. TODO: Maybe only search for .srt files?
			if checkForFileExtension(videoFileExtension, srtFilename[2]) {
				continue
			}

			// Skip all files with file names that do not contain YYxYY (e.g. 02x12) TODO: Maybe find other subtitle files too?
			match, _ := regexp.MatchString(strings.Join(result, "x"), g.Name())

			if !match {
				continue
			}

			fmt.Println("Found a match between " + f.Name() + " and " + g.Name())
			fmt.Println("Renaming " + g.Name() + " to " + videoFilename[1] + ".srt" + " (y/n?)")

			if confirmationRequiredBool && !askForConfirmation() {
				continue
			}

			err := os.Rename(pathToFolder + g.Name(), pathToFolder + videoFilename[1] + ".srt")
			fmt.Println("File was renamed to " + videoFilename[1] + ".srt")

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkForFileName(searchWord string, sourceWord string) bool {
	// Find all files with filename that contain the search word
	if searchWord != "" {
		containsSpecificFiles := strings.Contains(sourceWord, searchWord)

		if !containsSpecificFiles {
			return false
		}
	}

	return true
}

func checkForFileExtension(searchExtension string, sourceExtension string) bool {
	// Skip files that do not match the given video file extension.
	if searchExtension != "" {
		if searchExtension != sourceExtension {
			return false
		}
	}

	return true
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// containsString returns true if slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}