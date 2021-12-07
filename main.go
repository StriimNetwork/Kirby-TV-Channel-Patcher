package main

import (
	"fmt"
	"github.com/gabstv/go-bsdiff/pkg/bspatch"
	"github.com/wii-tools/GoNUSD"
	"github.com/wii-tools/wadlib"
	"io/ioutil"
	"os"
)

const (
	header     = "Kirby TV Channel Patcher\nBy: SketchMaster2001\n\n"
	sketch_url = "https://sketchmaster2001.github.io"
)

func main() {
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

	wadData, err := GoNUSD.Download(0x0001000148434d50, 257, false)
	if err != nil {
		errorFunc(err)
	}

	// Since NUS does not have a ticket for us and GoNUSD does not store the TMD data for some reason, we have to download the files ourselves.
	ticket, err := downloadFile(fmt.Sprintf("%s/kirby-tv/0001000148434d50.tik", sketch_url))
	if err != nil {
		errorFunc(err)
	}

	err = wadData.LoadTicket(ticket)
	if err != nil {
		return
	}

	// Download cert
	cert, err := downloadFile(fmt.Sprintf("%s/cert", sketch_url))

	wadData.CertificateChain = cert

	wadContents, err := wadData.GetWAD(wadlib.WADTypeCommon)
	if err != nil {
		errorFunc(err)
	}

	// Download patch
	patch, err := downloadFile(fmt.Sprintf("%s/kirby-tv/kirbo.patch", sketch_url))
	if err != nil {
		errorFunc(err)
	}

	patchedWad, err := bspatch.Bytes(wadContents, patch)
	if err != nil {
		errorFunc(err)
	}

	// Write the WAD to file. First create the WAD directory
	err = os.Mkdir("WAD", 0777)
	if err != nil {
		// If the directory exists, do nothing
		if os.IsExist(err) {

		} else {
			// We cannot handle this
			errorFunc(err)
		}
	}

	err = ioutil.WriteFile("WAD/Kirby-TV-Channel(Striim Network).wad", patchedWad, os.ModePerm)
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
