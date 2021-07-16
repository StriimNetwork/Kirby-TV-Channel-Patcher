package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"io"
	"net/http"
	"os"
)

func clear() {
	screen.Clear()
	screen.MoveTopLeft()
}

// Since we only give 2 options for the user, selection will do all the work for us without having to re-write the code.
func selection(input string, function func())  {
	switch input {
	case "1":
		function()
	case "2":
		os.Exit(0)
	default:
		// They shouldn't have entered this
		main()
	}
}

// DownloadFile downloads the WAD from the specified url
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create WAD directory
	err = os.Mkdir("WAD", 0777)
	if err != nil {
		// Handle if WAD directory exists
		if os.IsExist(err) {
		}
	}

	filepath = fmt.Sprintf("WAD/%s", filepath)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}