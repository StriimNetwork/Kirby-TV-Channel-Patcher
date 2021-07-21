package main

import (
	"Kirby-TV-Channel-Patcher/wad"
	"fmt"
	"github.com/gabstv/go-bsdiff/pkg/bspatch"
	"github.com/wii-tools/GoNUSD"
	"io/ioutil"
	"os"
)

const (
	header = "Kirby TV Channel Patcher\nBy: SketchMaster2001\n\n"
	sketch_url = "https://sketchmaster2001.github.io"
)


func main()  {
	clear()
	fmt.Printf(header)
	fmt.Printf("1.Start\n2.Exit\n\nChoose: ")
	var input string
	fmt.Scanln(&input)

	selection(input, prePatch, main)
}

func prePatch() {
	clear()
	fmt.Printf("%s", header)
	fmt.Printf("Welcome to the Kirby TV Channel Installation Process. The patcher will download the required files.\nThe entire process should take about 1 to 3 minutes depending on your computer CPU and internet speed.\n\n")
	fmt.Printf("1. Patch!\n2. Exit\n\nChoose: ")

	var input string
	fmt.Scanln(&input)

	selection(input, patch, prePatch)
}

func patch() {
	clear()
	fmt.Printf(header)
	fmt.Println("Patching Kirby TV Channel...")

	err := os.Mkdir("patching-dir", 0777)
	if err != nil {
		// Overwrite patching-dir folder if it exists
		if os.IsExist(err) {

		}
	}

	data, err := GoNUSD.Download(0x0001000148434d50, 257, false, false)
	if err != nil {
		errorFunc(err)
	}

	// Since NUS does not have a ticket for us and GoNUSD does not store the TMD data for some reason, we have to download the files ourselves.
	ticket, err := downloadFile(fmt.Sprintf("%s/kirby-tv/0001000148434d50.tik", sketch_url), "patching-dir/0001000148434d50.tik", true)
	if err != nil {
		errorFunc(err)
	}

	tmd, err := downloadFile(fmt.Sprintf("%s/kirby-tv/0001000148434d50.tmd", sketch_url), "patching-dir/0001000148434d50.tmd", true)
	if err != nil {
		errorFunc(err)
	}

	wadBuffer, err := getWadContents(tmd, ticket)
	if err != nil {
		errorFunc(err)
	}

	for i, content := range wadBuffer.TMD.Contents {
		// Create IV based on the content's Index
		iv := []byte{0x00, byte(content.Index), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		err := decryptAESCBC(wadBuffer.Ticket.TitleKey[:], iv, data[2+i], fmt.Sprintf("patching-dir/%08x.app", content.Index))
		if err != nil {
			errorFunc(err)
		}

		// Download and patch 00000002.app
		if content.Index == 00000002 {
			// Download patch
			patch, err := downloadFile(fmt.Sprintf("%s/kirby-tv/kirby-tv.patch", sketch_url), "", false)
			if err != nil {
				errorFunc(err)
			}

			// Create patch buffer
			patchedFile, err := bspatch.Bytes(data[2+i], patch)
			if err != nil {
				errorFunc(err)
			}

			// Write patch
			err = ioutil.WriteFile("patching-dir/00000002.app", patchedFile, 077)
			if err != nil {
				if os.IsExist(err) {
					// Overwrite current file.
				} else {
					// Error we cannot continue off of.
					errorFunc(err)
				}
			}
		}
	}

	// Download cert
	_, err = downloadFile(fmt.Sprintf("%s/cert", sketch_url), "patching-dir/0001000148434d50.certs", true)

	// Pack the WAD
	err = wad.Pack()
	if err != nil {
		errorFunc(err)
	}

	finish()
}

func finish() {
	clear()
	fmt.Printf(header)
	fmt.Printf("Patching is complete! Please move the WAD from the WAD folder on to your SD Card, and install the channel as normal.\nFor a more in depth guide, go to wii.guide/kirby-tv\n\n")
	fmt.Println("Press enter to exit the application.")

	var input string
	fmt.Scanln(&input)
}
