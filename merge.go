package cyclonedxmerge

import (
	"fmt"
	"slices"

	"github.com/CycloneDX/cyclonedx-go"
)

func MergeMap(input *[]cyclonedx.Dependency, prefix string) {
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

func Merge(items *[]cyclonedx.Component, input *[]cyclonedx.Component, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			item.BOMRef = fmt.Sprintf("%s%s", prefix, item.BOMRef)
			if !Has(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}
func MergeS(items *[]cyclonedx.Service, input *[]cyclonedx.Service, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			item.BOMRef = fmt.Sprintf("%s%s", prefix, item.BOMRef)
			if !HasS(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}
func MergeE(items *[]cyclonedx.ExternalReference, input *[]cyclonedx.ExternalReference, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			if !HasE(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}
func MergeC(items *[]cyclonedx.Composition, input *[]cyclonedx.Composition, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			if !HasC(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}
func MergeP(items *[]cyclonedx.Property, input *[]cyclonedx.Property, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			item.Name = fmt.Sprintf("%s%s", prefix, item.Name)
			if !HasP(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}
func MergeA(items *[]cyclonedx.Annotation, input *[]cyclonedx.Annotation, prefix string) {
	if items != nil && input != nil {
		for _, item := range *input {
			item.BOMRef = fmt.Sprintf("%s%s", prefix, item.BOMRef)
			if !HasA(items, &item) {
				*items = append(*items, item)
			}
		}
	}
}

func Has(items *[]cyclonedx.Component, input *cyclonedx.Component) bool {
	if items != nil {
		for _, item := range *items {
			if item.BOMRef == input.BOMRef {
				return true
			}
		}
	}

	return false
}

func HasS(items *[]cyclonedx.Service, input *cyclonedx.Service) bool {
	if items != nil {
		for _, item := range *items {
			if item.BOMRef == input.BOMRef {
				return true
			}
		}
	}

	return false
}

func HasE(items *[]cyclonedx.ExternalReference, input *cyclonedx.ExternalReference) bool {
	if items != nil {
		for _, item := range *items {
			if item.URL == input.URL && item.Type == input.Type {
				return true
			}
		}
	}

	return false
}

func HasC(items *[]cyclonedx.Composition, input *cyclonedx.Composition) bool {
	if items != nil {
		for _, item := range *items {
			if item.Aggregate == input.Aggregate {
				return true
			}
		}
	}

	return false
}

func HasP(items *[]cyclonedx.Property, input *cyclonedx.Property) bool {
	if items != nil {
		for _, item := range *items {
			if item.Name == input.Name && item.Value == input.Value {
				return true
			}
		}
	}

	return false
}

func HasA(items *[]cyclonedx.Annotation, input *cyclonedx.Annotation) bool {
	if items != nil {
		for _, item := range *items {
			if item.BOMRef == input.BOMRef {
				return true
			}
		}
	}

	return false
}
