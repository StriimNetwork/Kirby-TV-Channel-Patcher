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
func selection(input string, goToFunc func(), currentFunc func()) {
	switch input {
	case "1":
		goToFunc()
	case "2":
		os.Exit(0)
	default:
		// They shouldn't have entered this
		currentFunc()
	}
}

func errorFunc(err error) {
	clear()
	fmt.Printf(header)
	fmt.Println("An error has occurred. DM SketchMaster2001 #0713 on Discord for support, along with the error.")
	fmt.Printf("\nError: %s\n", err)
	fmt.Println("\nPress enter to exit the application.")

	// Only exists so the process doesn't die without user interaction
	var input string
	fmt.Scanln(&input)
	os.Exit(0)
}

// downloadFile downloads a file from the specified URL.
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		errorFunc(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorFunc(err)
	}

	return data, nil
}
