package merge

import (
	"slices"

	"github.com/CycloneDX/cyclonedx-go"
)

func MergeSBOM(b1 *cyclonedx.BOM, b2 *cyclonedx.BOM) {

	if b2.Metadata != nil {
		if b2.Metadata.Component != nil {
			topComponents := []cyclonedx.Component{*b2.Metadata.Component}
			merge(b1.Components, &topComponents)
			mergeD(b1.Dependencies, &[]cyclonedx.Dependency{{Ref: "root", Dependencies: &[]string{b2.Metadata.Component.BOMRef}}})

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

	mergeD(b1.Dependencies, b2.Dependencies)

}

func mergeD(d1 *[]cyclonedx.Dependency, d2 *[]cyclonedx.Dependency) {
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
