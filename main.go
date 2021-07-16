package main

import (
	"fmt"
	"os"
)

var header = "Kirby TV Channel Patcher\nBy: SketchMaster2001\n\n"


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

	selection(input, download, prePatch)
}

func download() {
	clear()
	fmt.Printf(header)
	fmt.Println("Patching Kirby TV Channel...")
	url := "soon_tm"

	err := DownloadFile("Kirby-TV-Channel-Patched(Striim).wad", url)
	if err != nil {
		clear()
		fmt.Printf(header)
		fmt.Println("An error has occurred. DM SketchMaster2001 #0713 on Discord for support, along with the error.")
		fmt.Printf("\nError: %s\n", err)
		fmt.Println("\nPress any key to exit the application.")

		// Only exists so the process doesn't die without user interaction
		var input string
		fmt.Scanln(&input)
		os.Exit(0)
	}

	finish()
}

func finish() {
	clear()
	fmt.Printf(header)
	fmt.Printf("Patching is complete! Please move the WAD from the WAD folder on to your SD Card, and install the channel as normal.\nFor a more in depth guide, go to 'wii.guide/kirby-tv'.\n\n")
	fmt.Println("Press any key to exit the application.")

	var input string
	fmt.Scanln(&input)
}
