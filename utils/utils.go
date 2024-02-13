package utils

import (
	"fmt"
	"hash/fnv"
	"slices"
	"strings"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func NewBOM(rootComponent *cyclonedx.Component) *cyclonedx.BOM {
	sbom := cyclonedx.NewBOM()
	sbom.Metadata = &cyclonedx.Metadata{
		Tools: &[]cyclonedx.Tool{{
			Vendor:  "fnxpt",
			Name:    "cyclonedx-merge",
			Version: "0.0.2",
		}},
		Timestamp: time.Now().Format(time.RFC3339),
		Component: rootComponent,
	}
	sbom.Dependencies = &[]cyclonedx.Dependency{
		{
			Ref: rootComponent.BOMRef,
		},
	}

	annotations := make([]cyclonedx.Annotation, 0)
	components := make([]cyclonedx.Component, 0)
	compositions := make([]cyclonedx.Composition, 0)
	externalReferences := make([]cyclonedx.ExternalReference, 0)
	properties := make([]cyclonedx.Property, 0)
	services := make([]cyclonedx.Service, 0)
	vulnerabilities := make([]cyclonedx.Vulnerability, 0)

	sbom.Annotations = &annotations
	sbom.Components = &components
	sbom.Compositions = &compositions
	sbom.ExternalReferences = &externalReferences
	sbom.Properties = &properties
	sbom.Services = &services
	sbom.Vulnerabilities = &vulnerabilities

	return sbom
}

type comparableType interface {
	cyclonedx.Annotation | cyclonedx.Component | cyclonedx.Composition | cyclonedx.ExternalReference | cyclonedx.Property | cyclonedx.Service
}

func Merge[T comparableType](items *[]T, input *[]T) {
	if items != nil && input != nil {
		for _, item := range *input {
			if !contains(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}

func MergeDependencies(d1 *[]cyclonedx.Dependency, d2 *[]cyclonedx.Dependency) {
	if d2 != nil {

		var tmp = make(map[string]cyclonedx.Dependency)

		for _, item := range *d1 {
			tmp[item.Ref] = item
		}

		for _, item := range *d2 {
			key := item.Ref

			if _, ok := tmp[key]; ok {
				if item.Dependencies != nil {
					if tmp[key].Dependencies != nil {
						for _, dependency := range *item.Dependencies {
							if !slices.Contains(*tmp[key].Dependencies, dependency) {
								*tmp[item.Ref].Dependencies = append(*tmp[item.Ref].Dependencies, dependency)
							}
						}
					} else {
						tmp[key] = item
					}
				}
			} else {
				tmp[item.Ref] = item
			}
		}

		newDependencies := make([]cyclonedx.Dependency, 0)
		for _, item := range tmp {
			newDependencies = append(newDependencies, item)
		}
		*d1 = newDependencies

	}
}

func MergeDependenciesFromComponents(d1 *[]cyclonedx.Dependency, d2 *[]cyclonedx.Component, key string) {
	if d2 != nil {
		tmp := make(map[string]cyclonedx.Dependency)

		for _, item := range *d1 {
			tmp[item.Ref] = item
		}

		if _, ok := tmp[key]; !ok || tmp[key].Dependencies == nil {
			dependencies := []string{}
			tmp[key] = cyclonedx.Dependency{Ref: key, Dependencies: &dependencies}
		}

		for _, item := range *d2 {

			if !slices.Contains(*tmp[key].Dependencies, item.BOMRef) {
				*tmp[key].Dependencies = append(*tmp[key].Dependencies, item.BOMRef)
			}
		}

		newDependencies := make([]cyclonedx.Dependency, 0)
		for _, item := range tmp {
			newDependencies = append(newDependencies, item)
		}
		*d1 = newDependencies
	}
}

func MergeDependenciesMissing(d1 *[]cyclonedx.Dependency, d2 *[]cyclonedx.Dependency) []string {
	missingDependencies := make([]string, 0)

	if d2 != nil {
		tmp := make(map[string]cyclonedx.Dependency)

		for _, item := range *d1 {
			tmp[item.Ref] = item
		}

		for _, item := range *d2 {
			key := item.Ref

			if _, ok := tmp[key]; ok {
				if item.Dependencies != nil {
					if tmp[key].Dependencies != nil {
						if item.Ref != "root" {

							//CHECK IF THERE IS ANY KEY THAT HAS KEY+|
							//IF KEY+| exists
							//calculate KEY+| FOR ITEM :KEY2
							//IF KEY2 doesnt exist
							//ADD KEY2 AND KEY3 to TMP
							//CREATE COMPONENT WITH KEY2
							//ELSE
							//check if tmp[key] dependencies and item.dependencies are the same
							//if !equal
							//calculate KEY+| FOR ITEM :KEY2
							//calculate KEY+| FOR TMP[KEY] :KEY3
							//ADD KEY2 AND KEY3 to TMP
							//CREATE COMPONENT WITH KEY2
							//CREATE COMPONENT WITH KEY3
							//DELETE ORIGINAL COMPONENT

							less := func(a, b string) bool { return a < b }
							equalIgnoreOrder := cmp.Diff(tmp[key].Dependencies, item.Dependencies, cmpopts.SortSlices(less)) == ""
							if !equalIgnoreOrder {

								slices.Sort(*tmp[key].Dependencies)
								slices.Sort(*item.Dependencies)

								value1 := strings.Join(*tmp[key].Dependencies, "|")
								value2 := strings.Join(*item.Dependencies, "|")

								h1 := fnv.New32a()
								h1.Write([]byte(value1))

								h2 := fnv.New32a()
								h2.Write([]byte(value2))

								fmt.Printf("%s|%d != %s|%d\n", key, h1.Sum32(), key, h2.Sum32())
							}
						} else {
							*tmp[item.Ref].Dependencies = append(*tmp[item.Ref].Dependencies, *item.Dependencies...)
						}

						// for _, dependency := range *item.Dependencies {
						// 	if !slices.Contains(*tmp[key].Dependencies, dependency) {

						// 		tmp[fmt.Sprintf("%s|%s", key, dependency)] =
						// 			fmt.Printf("WARN: %s doensn't have %s\n", item.Ref, dependency)
						// 		missingDependencies = append(missingDependencies, dependency)
						// 	} else {
						// 		*tmp[item.Ref].Dependencies = append(*tmp[item.Ref].Dependencies, dependency)
						// 	}

						// }
					}
				} else {
					tmp[key] = item
				}
			} else {
				tmp[item.Ref] = item
			}
		}

		newDependencies := make([]cyclonedx.Dependency, 0)
		for _, item := range tmp {
			newDependencies = append(newDependencies, item)
		}
		*d1 = newDependencies

	}

	return missingDependencies
}

func contains[T comparableType](items *[]T, input *T) bool {
	if items != nil {
		for _, item := range *items {
			if isEqual([]*T{&item, input}) {
				return true
			}
		}
	}

	return false
}

func isEqual[T comparableType](input []*T) bool {
	switch values := any(input).(type) {
	case []*cyclonedx.Annotation:
		return values[0].BOMRef == values[1].BOMRef
	case []*cyclonedx.Component:
		return values[0].BOMRef == values[1].BOMRef
	case []*cyclonedx.Service:
		return values[0].BOMRef == values[1].BOMRef
	case []*cyclonedx.ExternalReference:
		return values[0].URL == values[1].URL && values[0].Type == values[1].Type
	case []*cyclonedx.Property:
		return values[0].Name == values[1].Name
	default:
		return false
	}
}
