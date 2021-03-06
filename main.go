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

var pathToFolder string
var searchWord string
var videoFileExtension string
var videoFileExtensionSlice []string
var confirmationRequired string
var subtitleFileExtension string
var subtitleFileExtensionSlice []string
var confirmationRequiredBool bool

func main() {
	readCommandLineArguments()

	fileRenameCount := 0

	files, err := ioutil.ReadDir(pathToFolder)
	if err != nil {
		log.Fatal(err)
	}

	// Find matching video files.
	for _, f := range files {

		// Skip all files that did not match the search string.
		if !checkForFileName(searchWord, f.Name()) {
			continue
		}

		filenameExpression, _ := regexp.Compile("(.*)(\\.\\w*)$")
		videoFilename := filenameExpression.FindStringSubmatch(f.Name())

		// Check if the file even has a file extension and if that file extension is the video type.
		if len(videoFilename) < 2 {
			continue
		}

		// Check if the file even has a file extension and if that file extension is the video type.
		if checkForFileExtension(videoFileExtensionSlice, videoFilename[2]) == "" {
			continue
		}

		// Skip all files that do not contain "SYYEYY" (e.g. S02E12)
		match, _ := regexp.MatchString("(?i)(S\\w{2,2}E\\w{2,2})", f.Name())
		if !match {
			continue
		}

		// Get the season and episode number from the file name.
		r, _ := regexp.Compile("(?i)S(\\w{2,2})E(\\w{2,2})")
		seasonAndEpisodeNumber := r.FindStringSubmatch(f.Name())
		seasonAndEpisodeNumber = append(seasonAndEpisodeNumber[:0], seasonAndEpisodeNumber[1:]...)

		// Find matching .srt files.
		for _, g := range files {
			// Skip all files that did not match the search string.
			if !checkForFileName(searchWord, g.Name()) {
				continue
			}

			filenameExpression, _ := regexp.Compile("(.*)(\\.\\w*)$")
			srtFilename := filenameExpression.FindStringSubmatch(g.Name())

			// Skip all files without the defined subtitle extension.
			if len(srtFilename) < 2 {
				continue
			}

			// Skip all files without the defined subtitle extensions.
			subtitleFileExtension = checkForFileExtension(subtitleFileExtensionSlice, srtFilename[2])
			if subtitleFileExtension == "" {
				continue
			}

			// Skip all files with file names that do not contain YYxYY (e.g. 02x12) or SYYEYY (e.g. S02E12)
			subtitleNameMatchOne, _ := regexp.MatchString(strings.ToUpper(strings.Join(seasonAndEpisodeNumber, "x")), strings.ToUpper(g.Name()))
			subtitleNameMatchTwo, _ := regexp.MatchString(strings.ToUpper("S" + strings.Join(seasonAndEpisodeNumber, "E")), strings.ToUpper(g.Name()))

			if !subtitleNameMatchOne && !subtitleNameMatchTwo {
				continue
			}

			fmt.Println("Found a match between " + f.Name() + " and " + g.Name())
			fmt.Println("Renaming " + g.Name() + " to " + videoFilename[1] + subtitleFileExtension + " (y/n?)")

			if confirmationRequiredBool && !askForConfirmation() {
				continue
			}

			err := os.Rename(pathToFolder + g.Name(), pathToFolder + videoFilename[1] + subtitleFileExtension)

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

func readCommandLineArguments () {
	pathToFolderPtr := flag.String("path", "", "The path to the folder with video and srt files or directly to a file. If none is specified the folder with the executable is used.")
	searchWordPtr := flag.String("searchWord", "", "Provide a unique word to identify the files by. E.g. 'Queens' or 'Mother'")
	videoFileExtensionPtr := flag.String("videoFileExtension", ".mkv,.mp4", "Provide the video file extension. E.g. '.mkv,.mp4'. Defaults to .mkv and .mp4.")
	subtitleFileExtensionPtr := flag.String("subtitleFileExtension", ".srt,.sub", "Provide the subtitle file extension. E.g. '.srt,.sub'. Defaults to .srt and .sub.")
	confirmationRequiredPtr := flag.String("enableConfirmation", "", "Enable the confirmation before every rename. If enabled every rename needs to be confirmed by typing 'y' or denied by typing 'n'")

	flag.Parse()

	pathToFolder =  *pathToFolderPtr
	searchWord =  *searchWordPtr
	videoFileExtension =  *videoFileExtensionPtr
	confirmationRequired = *confirmationRequiredPtr
	subtitleFileExtension = *subtitleFileExtensionPtr

	// If no path was specified, use the current folder.
	if pathToFolder == "" {
		currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		pathToFolder = currentPath
	} else {
		pathToFolder = filepath.Dir(pathToFolder)
	}


	// Check if the last character in the path is a slash. If not, add it.
	if !strings.HasSuffix(pathToFolder, "/") && !strings.HasSuffix(pathToFolder, "\\") {
		pathToFolder = pathToFolder + "\\"
	}

	confirmationRequiredBool = false
	if confirmationRequired != "" {
		confirmationRequiredBool = true
	}

	fmt.Println("Looking for files in " + pathToFolder)

	if searchWord != "" {
		fmt.Println("Looking for files that include the expression: " + searchWord)
	} else {
		fmt.Println("Looking for all files")
	}

	if videoFileExtension != "" {
		videoFileExtensionSlice = strings.Split(videoFileExtension, ",")

		for i := range videoFileExtensionSlice {
			fmt.Println("Looking for video files with file extension: " + videoFileExtensionSlice[i])
		}
	} else {
		fmt.Println("Looking for all video files")
	}

	if subtitleFileExtension != "" {
		subtitleFileExtensionSlice = strings.Split(subtitleFileExtension, ",")

		for i := range subtitleFileExtensionSlice {
			fmt.Println("Looking for subtitle files with file extension: " + subtitleFileExtensionSlice[i])
		}
	} else {
		fmt.Println("Looking for all subtitle files")
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

func checkForFileExtension(searchExtensions []string, sourceExtension string) string {

	if len(searchExtensions) == 0 {
		return sourceExtension;
	}

	// Skip files that do not match the given video file extension.
	for i := range searchExtensions {
		if sourceExtension == searchExtensions[i] {
			return searchExtensions[i]
		}
	}

	return ""
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