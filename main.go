package main

import (
	"fmt"
	"log"
	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"os"
	"strconv"
)

func main() {
	pathToFolderPtr := flag.String("path", "", "The path to the folder with video and srt files. If non is specified the folder with the executable is used.")
	searchWordPtr := flag.String(".", "", "Provide a unique word to identify the files by. E.g. 'Queens' or 'Mother'")
	videoFileExtensionPtr := flag.String("videoFileExtension", ".mkv", "Provide the video file extension. Defaults to .mkv.")
	subtitleFileExtensionPtr := flag.String("subtitleFileExtension", ".srt", "Provide the subtitle file extension. Defaults to .srt.")
	confirmationRequiredPtr := flag.String("disableConfirmation", "", "Disable the confirmation before every rename.")

	flag.Parse()

	pathToFolder :=  *pathToFolderPtr
	pathToFolder = ""
	searchWord :=  *searchWordPtr
	videoFileExtension :=  *videoFileExtensionPtr
	confirmationRequired := *confirmationRequiredPtr
	subtitleFileExtension := *subtitleFileExtensionPtr
	fileRenameCount := 0

	// If no path was specified, use the current folder.
	if pathToFolder == "" {
		currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		pathToFolder = currentPath;
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

	if videoFileExtension != "" {
		fmt.Println("Looking for video files with file extension: " + videoFileExtension)
	} else {
		fmt.Println("Looking for all video files")
	}

	if subtitleFileExtension != "" {
		fmt.Println("Looking for subtitle files with file extension: " + subtitleFileExtension)
	} else {
		fmt.Println("Looking for all subtitle files")
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

		// Check if the file even has a file extension and if that file extension is the video type.
		if len(videoFilename) < 2 || !checkForFileExtension(videoFileExtension, videoFilename[2]) {
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

			// Skip all files without the defined subtitle extension.
			if !checkForFileExtension(subtitleFileExtension, srtFilename[2]) {
				continue
			}

			// Skip all files with file names that do not contain YYxYY (e.g. 02x12)
			match, _ := regexp.MatchString(strings.ToUpper(strings.Join(result, "x")), strings.ToUpper(g.Name()))

			if !match {
				// Nothing found? Try to find subtitle files that contain SYYEYY (e.g. S02E12)
				match, _ = regexp.MatchString(strings.ToUpper("S" + strings.Join(result, "E")), strings.ToUpper(g.Name()))

				if !match {
					continue
				}
			}

			fmt.Println("Found a match between " + f.Name() + " and " + g.Name())
			fmt.Println("Renaming " + g.Name() + " to " + videoFilename[1] + ".srt" + " (y/n?)")

			if confirmationRequiredBool && !askForConfirmation() {
				continue
			}

			err := os.Rename(pathToFolder + g.Name(), pathToFolder + videoFilename[1] + ".srt")

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("File was renamed to " + videoFilename[1] + ".srt")
				fileRenameCount++
			}
		}
	}

	if fileRenameCount > 0 {
		fmt.Println(strconv.Itoa(fileRenameCount) + " files renamed.")
	} else {
		fmt.Println("No files found to rename.")
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