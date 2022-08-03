package main

import (
	"bufio"
	"flag"
	"github.com/nyudlts/go-aspace"
	"log"
	"os"
	"strings"
)

var (
	config      string
	environment string
	inputFile   string
	httpHandle  = "http://hdl.handle.net/2333.1/"
)

func init() {
	flag.StringVar(&config, "config", "", "")
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&inputFile, "input-file", "", "")
}

func main() {
	flag.Parse()

	client, err := aspace.NewClient(config, environment, 20)
	if err != nil {
		panic(err)
	}

	inFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	logFile, err := os.Create("handle-update.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		repoID, doID, err := aspace.URISplit(scanner.Text())
		if err != nil {
			log.Printf("[ERROR] %s", strings.ReplaceAll(err.Error(), "\n", " "))
			continue
		}

		do, err := client.GetDigitalObject(repoID, doID)
		if err != nil {
			log.Printf("[ERROR] %s", strings.ReplaceAll(err.Error(), "\n", " "))
			continue
		}
		oldFileVersions := do.FileVersions
		newFileVersions := []aspace.FileVersion{}
		for _, fv := range oldFileVersions {
			if strings.Contains(fv.FileURI, httpHandle) == true {
				fv.FileURI = strings.ReplaceAll(fv.FileURI, "http", "https")
				newFileVersions = append(newFileVersions, fv)
			} else {
				newFileVersions = append(newFileVersions, fv)
			}
		}
		do.FileVersions = newFileVersions

		msg, err := client.UpdateDigitalObject(repoID, doID, do)
		if err != nil {
			log.Printf("[ERROR] %s", strings.ReplaceAll(err.Error(), "\n", " "))
			continue
		}
		log.Printf("%s\t%s", do.URI, strings.ReplaceAll(msg, "\n", ""))
	}

}
