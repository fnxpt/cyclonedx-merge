package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
)

var version = "0.0.1"

var nested = false
var sbom *cyclonedx.BOM
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

	flag.BoolVar(&nested, "nested", false, "nested merge")
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

	fillSBOM()
	flag.Parse()
	postSBOM()
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

	mergeSBOM(bom)

	return nil
}

var topDependencies = make([]string, 0)

func mergeSBOM(value *cyclonedx.BOM) {
	var prefix string

	if value.Metadata != nil {
		if value.Metadata.Component != nil {
			topComponents := []cyclonedx.Component{*value.Metadata.Component}
			if nested {
				prefix = fmt.Sprintf("%s|", value.Metadata.Component.BOMRef)
			}

			addIfNew(sbom.Components, &topComponents, "")
			topDependencies = append(topDependencies, value.Metadata.Component.BOMRef)

			if value.Metadata.Component.Components != nil {
				addIfNew(sbom.Components, value.Metadata.Component.Components, prefix)
			}
		}
	}

	addIfNew(sbom.Components, value.Components, prefix)
	// addIfNew(sbom.Services, value.Services, prefix)
	// addIfNew(sbom.ExternalReferences, value.ExternalReferences, prefix)

	addIfNewMap(value.Dependencies, prefix)
	// addIfNew(sbom.Compositions, value.Compositions)
	// addIfNew(sbom.Properties, value.Properties)
	// addIfNew(sbom.Vulnerabilities, value.Vulnerabilities, prefix)
	// addIfNew(sbom.Annotations, value.Annotations, prefix)
}

func fillSBOM() {

	sbom = cyclonedx.NewBOM()
	sbom.Metadata = &cyclonedx.Metadata{
		Tools: &[]cyclonedx.Tool{{
			Vendor:  "fnxpt",
			Name:    "cyclonedx-merge",
			Version: version,
		}},
		Timestamp: time.Now().String(), //TODO: RIGHT FORMAT
		Component: &cyclonedx.Component{
			BOMRef: "root",
			Name:   "root",
			Type:   cyclonedx.ComponentTypeApplication,
		},
	}

	annotations := make([]cyclonedx.Annotation, 0)
	components := make([]cyclonedx.Component, 0)
	compositions := make([]cyclonedx.Composition, 0)
	dependencies := make([]cyclonedx.Dependency, 0)
	externalReferences := make([]cyclonedx.ExternalReference, 0)
	properties := make([]cyclonedx.Property, 0)
	services := make([]cyclonedx.Service, 0)
	vulnerabilities := make([]cyclonedx.Vulnerability, 0)

	sbom.Annotations = &annotations
	sbom.Components = &components
	sbom.Compositions = &compositions
	sbom.Dependencies = &dependencies
	sbom.ExternalReferences = &externalReferences
	sbom.Properties = &properties
	sbom.Services = &services
	sbom.Vulnerabilities = &vulnerabilities

}

func postSBOM() {

	for _, item := range tmp {
		*sbom.Dependencies = append(*sbom.Dependencies, item)
	}
	*sbom.Dependencies = append(*sbom.Dependencies, cyclonedx.Dependency{
		Ref:          "root",
		Dependencies: &topDependencies,
	})

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

var tmp = make(map[string]cyclonedx.Dependency)

func addIfNewMap(input *[]cyclonedx.Dependency, prefix string) {
	if input != nil {

		for _, item := range *input {
			key := item.Ref

			if nested && prefix[:len(prefix)-1] != item.Ref {
				key = fmt.Sprintf("%s%s", prefix, key)
			}
			if _, ok := tmp[key]; ok {
				for _, dependency := range *item.Dependencies {
					dependency = fmt.Sprintf("%s%s", prefix, dependency)
					if !slices.Contains(*tmp[key].Dependencies, dependency) {
						*tmp[item.Ref].Dependencies = append(*tmp[item.Ref].Dependencies, dependency)
					}
				}
			} else {
				item.Ref = key
				//TODO: PREFIX DEPENDENCIES
				if len(prefix) > 0 {
					tmp[key] = cyclonedx.Dependency{
						Ref:          key,
						Dependencies: &[]string{},
					}
					if item.Dependencies != nil {
						for _, dep := range *item.Dependencies {
							dependency := fmt.Sprintf("%s%s", prefix, dep)
							*tmp[key].Dependencies = append(*tmp[key].Dependencies, dependency)
						}
					}
				} else {
					tmp[key] = item
				}
			}
		}
	}
}
func addIfNew(items *[]cyclonedx.Component, input *[]cyclonedx.Component, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			item.BOMRef = fmt.Sprintf("%s%s", prefix, item.BOMRef)
			if !has(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}

func has(items *[]cyclonedx.Component, input *cyclonedx.Component) bool {
	if items != nil {
		for _, item := range *items {
			if item.BOMRef == input.BOMRef {
				return true
			}
		}
	}

	return false
}
