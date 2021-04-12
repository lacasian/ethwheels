package main

import (
	"flag"

	"github.com/lacasian/ethwheels/ethgen"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting to generate ABI bindings")
	abisDir := flag.String("abi-folder", "ethgen/testdata/_source", "Folder containing ABI JSONs")
	packagePath := flag.String("package-path", "ethgen/testdata", "Path where to generate packages. Final folder represents package name")
	flag.Parse()

	if abisDir == nil || packagePath == nil {
		log.Fatal("please specify \"abi-folder\" and \"package-path \"")
	}
	err := ethgen.NewFromABIs(*abisDir, *packagePath)
	if err != nil {
		log.Fatal(err)
	}

}
