package flatmerge

import (
	"slices"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
)

func MergeSBOM(b1 *cyclonedx.BOM, b2 *cyclonedx.BOM) {
	if b1 == nil || b2 == nil {
		panic("cannot merge 2 nil objects")
	}

	if b2.Metadata != nil {
		if b2.Metadata.Component != nil {
			topComponents := []cyclonedx.Component{*b2.Metadata.Component}
			merge(b1.Components, &topComponents)

			mergeDependencies(b1.Dependencies, &[]cyclonedx.Component{*b2.Metadata.Component}, "root")

			if b2.Metadata.Component.Components != nil {
				merge(b1.Components, b2.Metadata.Component.Components)
			}
		}
	}

	merge(b1.Components, b2.Components)
	merge(b1.Services, b2.Services)
	merge(b1.ExternalReferences, b2.ExternalReferences)
	merge(b1.Compositions, b2.Compositions)
	merge(b1.Properties, b2.Properties)
	merge(b1.Annotations, b2.Annotations)

	mergeDependencies(b1.Dependencies, b2.Components, b2.Metadata.Component.BOMRef)
}

func NewBOM() *cyclonedx.BOM {
	sbom := cyclonedx.NewBOM()
	sbom.Metadata = &cyclonedx.Metadata{
		Tools: &[]cyclonedx.Tool{{
			Vendor:  "fnxpt",
			Name:    "cyclonedx-merge",
			Version: "0.0.2",
		}},
		Timestamp: time.Now().Format(time.RFC3339),
		Component: &cyclonedx.Component{
			BOMRef: "root",
			Name:   "root",
			Type:   cyclonedx.ComponentTypeApplication,
		},
	}
	sbom.Dependencies = &[]cyclonedx.Dependency{
		{
			Ref: "root",
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

func merge[T comparableType](items *[]T, input *[]T) {
	if items != nil && input != nil {
		for _, item := range *input {
			if !contains(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}

func mergeDependencies(d1 *[]cyclonedx.Dependency, d2 *[]cyclonedx.Component, key string) {
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
