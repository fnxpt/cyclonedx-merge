package main

import (
	"cyclonedx-merge/merge"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
)

var version = "0.0.1"

var sbom *cyclonedx.BOM = merge.NewBOM()
var outputFormat = cyclonedx.BOMFileFormatJSON
var output = os.Stdout

func main() {
	parseArguments()
}

func showHelpMenu() {
	fmt.Println("usage: cyclonedx-merge [options]")
	fmt.Println("options:")
	os.Exit(0)
}

func parseArguments() {
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	flag.Usage = func() {
		showHelpMenu()
		flag.PrintDefaults()
	}

	flag.Func("file", "merges file", fileMerge)
	flag.Func("dir", "merges files in directory", dirMerge)
	flag.Func("format", "output format - json/xml (default: json)", func(value string) error {
		switch value {
		case "json":
			outputFormat = cyclonedx.BOMFileFormatJSON
		case "xml":
			outputFormat = cyclonedx.BOMFileFormatXML
		default:
			return fmt.Errorf("invalid output format")
		}
		return nil
	})
	flag.Func("output", "output file (default: stdout)", func(value string) error {
		file, err := os.Create(value)

		if err != nil {
			fmt.Printf("unable to create file %s\n", value)
			return err
		}
		output = file
		return nil
	})

	flag.Parse()
	writeSBOM()
}

func fileMerge(value string) error {
	// fmt.Printf("Processing file %s\n", value)
	if _, err := os.Stat(value); os.IsNotExist(err) {
		fmt.Printf("file %s doesn't exist\n", value)
		return err
	}

	file, err := os.Open(value)

	if err != nil {
		fmt.Printf("unable to open file %s\n", value)
		return err
	}

	bom, err := parseSBOM(file)

	if err != nil {
		fmt.Printf("unable to parse file %s\n", value)
		return err
	}

	merge.MergeSBOM(sbom, bom)

	return nil
}

func writeSBOM() {

	encoder := cyclonedx.NewBOMEncoder(output, outputFormat)
	encoder.Encode(sbom)
}

func parseSBOM(input io.Reader) (*cyclonedx.BOM, error) {

	bom := &cyclonedx.BOM{}
	decoder := cyclonedx.NewBOMDecoder(input, cyclonedx.BOMFileFormatJSON)
	err := decoder.Decode(bom)

	if err != nil {
		return nil, err
	}

	return bom, err
}

func dirMerge(value string) error {
	if _, err := os.Stat(value); os.IsNotExist(err) {
		fmt.Printf("directory %s doesn't exist\n", value)
		return err
	}

	entries, err := os.ReadDir(value)
	if err != nil {
		fmt.Printf("unable to read directory %s\n", value)
		return err
	}

	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") || strings.HasSuffix(e.Name(), ".xml") {
			fileMerge(fmt.Sprintf("%s/%s", value, e.Name()))
		}
	}

	return nil
}
