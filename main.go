package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fnxpt/cyclonedx-merge/flatmerge"
	"github.com/fnxpt/cyclonedx-merge/merge"
	"github.com/fnxpt/cyclonedx-merge/utils"

	"github.com/CycloneDX/cyclonedx-go"
)

type MergeMode int

const (
	MergeModeNormal MergeMode = iota
	MergeModeFlat
	MergeModeSmart
)

var version = "0.0.6"

var rootComponent = &cyclonedx.Component{
	BOMRef: "root",
	Name:   "root",
	Type:   cyclonedx.ComponentTypeApplication,
}

var sbom *cyclonedx.BOM
var mode = MergeModeNormal
var outputFormat = cyclonedx.BOMFileFormatJSON
var output = os.Stdout

func main() {
	parseArguments()
}

func showHelpMenu() {
	fmt.Printf("usage: cyclonedx-merge(%s) [options]\n", version)
	fmt.Println("options:")
	os.Exit(0)
}

func parseArguments() {
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	flag.Usage = func() {
		showHelpMenu()
		flag.PrintDefaults()
	}

	flag.Func("group", "group of the merged parent component", func(value string) error {
		rootComponent.Group = value
		return nil
	})

	flag.Func("name", "name of the merged parent component", func(value string) error {
		rootComponent.Name = value
		return nil
	})

	flag.Func("version", "version of the merged parent component", func(value string) error {
		rootComponent.Version = value
		return nil
	})

	flag.Func("bomref", "BOMRef of the merged parent component", func(value string) error {
		rootComponent.BOMRef = value
		return nil
	})

	flag.Func("type", "type of the aggregator component", func(value string) error {
		rootComponent.Type = cyclonedx.ComponentType(value)
		return nil
	})

	flag.Func("file", "merges file", func(value string) error {
		sbom = utils.NewBOM(rootComponent)
		return fileMerge(value)
	})

	flag.Func("dir", "merges files in directory", func(value string) error {
		sbom = utils.NewBOM(rootComponent)
		return dirMerge(value)
	})

	flag.Func("mode", "merge mode - normal/flat/smart (default: normal)", func(value string) error {
		switch value {
		case "normal":
			mode = MergeModeNormal
		case "flat":
			mode = MergeModeFlat
		case "smart":
			mode = MergeModeSmart
		default:
			return fmt.Errorf("invalid mode")
		}
		return nil
	})

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

	switch mode {
	case MergeModeNormal:
		merge.MergeSBOM(sbom, bom)
	case MergeModeFlat:
		flatmerge.MergeSBOM(sbom, bom)
	case MergeModeSmart:
		panic("not implemented yet")
		// smartmerge.MergeSBOM(sbom, bom)
	}

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
