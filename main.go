package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func main() {

	var filePath = "./Dockerfile"
	var hardCheck = "false"
	var regPattern = "registry.gitlab.com/mycompany.de/infrastructure/images/*"
	color.NoColor = false // Allow to color output in CI job logs (Gitlab-CI, GitHub Actions)

	filePathEnv, exists := os.LookupEnv("DOCKERFILE_PATH") // Use if Dockerfile path is defined by external env
	if exists {
		exists = true
	} else {
		filePathEnv = filePath
	}
	hardCheckEnv, exists := os.LookupEnv("HARD_CHECK") // "true" – if Dockerfile has external image, get exit 1
	if exists {
		exists = true
	} else {
		hardCheckEnv = hardCheck
	}
	regPatternEnv, exists := os.LookupEnv("REG_PATTERN") // Use if registry pattern is defined by external env
	if exists {
		exists = true
	} else {
		regPatternEnv = regPattern
	}

	filePathPtr := flag.String("f", filePathEnv, "Path to Dockerfile")                           // "file" flag
	hardCheckPtr := flag.String("m", hardCheckEnv, "Enable(true) and disable(false) Hard check") // "hard-check" flag
	regPtr := flag.String("p", regPatternEnv, "Pattern to find correct image")                   // "pattern" flag
	flag.Parse()

	file, err := os.ReadFile(*filePathPtr)
	if err != nil {
		log.Fatal(err)
	}

	refString := string(file)
	lookFor, _ := regexp.Compile("(?m:^FROM .+)") // Check for images in Dockerfile
	result := lookFor.FindAllString(refString, -1)
	legalImages := make([]string, len(result))
	var imageString string

loop:
	for i := 0; i < len(result); i++ {
		imageString = result[i]
		matched, _ := regexp.Compile(*regPtr)
		imageString = strings.TrimPrefix(imageString, "FROM ")
		imageRegExp, _ := regexp.Compile("(^\\S+\\b)")
		justImage := imageRegExp.FindAllString(imageString, -1) // Create images list
		if matched.MatchString(imageString) {
			legalImages[i] = justImage[0]
			if justImage[0] != imageString {
				newImageRegExp, _ := regexp.Compile("\\S+$")
				newImage := newImageRegExp.FindAllString(imageString, -1)
				legalImages = append(legalImages, newImage[0]) // Create local images list
			}
		}

		alertImage := color.New(color.FgRed, color.Bold) // Set colors to display in output
		wrongImage := color.New(color.FgHiYellow, color.Bold)
		var trueImage bool
		if err != nil {
			log.Fatal(err)
		}
		for k := 0; k < len(legalImages); k++ {
			if justImage[0] == legalImages[k] {
				trueImage = true
			}
		}
		if trueImage == true {
			fmt.Println(justImage[0], "is the correct image to use.")
		} else {
			if *hardCheckPtr == "true" {
				err1, _ := alertImage.Printf("%q is NOT a local image. You can't use it! ", justImage[0])
				log.Fatal(err1)
			} else {
				wrongImage.Printf("%q is NOT a local image. Don't use it!\n", justImage[0])
			}
			break loop
		}
	}
}
