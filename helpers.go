package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/inancgumus/screen"
	"github.com/wii-tools/wadlib"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func clear() {
	screen.Clear()
	screen.MoveTopLeft()
}

// Since we only give 2 options for the user, selection will do all the work for us without having to re-write the code.
func selection(input string, goToFunc func(), currentFunc func())  {
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
func downloadFile(url string, outpath string, keepFile bool) ([]byte, error){
	resp, err := http.Get(url)
	if err != nil {
		errorFunc(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorFunc(err)
	}

	if keepFile {
		err = ioutil.WriteFile(outpath, data, 0777)
		if err != nil {
			errorFunc(err)
		}
	}

	return data, nil
}

func getWadContents(tmd []byte, ticket []byte) (wadlib.WAD, error) {
	// Create empty WAD
	wad := wadlib.WAD{}

	// Load the tmd
	err := wad.LoadTMD(tmd)
	if err != nil {
		return wadlib.WAD{}, err
	}

	// Load the ticket
	err = wad.LoadTicket(ticket)
	if err != nil {
		return wadlib.WAD{}, err
	}

	return wad, nil
}

// decryptAESCBC is used to decrypt the app files into a format that can be packed into a WAD.
func decryptAESCBC(key []byte, iv []byte, data []byte, outPath string) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(data, data)

	err = ioutil.WriteFile(outPath, data, 0777)
	if err != nil {
		return err
	}

	return nil
}
