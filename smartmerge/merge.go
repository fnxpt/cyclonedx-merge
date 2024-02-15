package smartmerge

import (
	"github.com/fnxpt/cyclonedx-merge/utils"

	"github.com/CycloneDX/cyclonedx-go"
)

func MergeSBOM(b1 *cyclonedx.BOM, b2 *cyclonedx.BOM) {
	if b1 == nil || b2 == nil {
		panic("cannot merge 2 nil objects")
	}

	if b2.Metadata != nil {
		if b2.Metadata.Component != nil {
			topComponents := []cyclonedx.Component{*b2.Metadata.Component}
			utils.Merge(b1.Components, &topComponents)
			utils.MergeDependenciesMissing(b1.Dependencies, &[]cyclonedx.Dependency{{Ref: b1.Metadata.Component.BOMRef, Dependencies: &[]string{b2.Metadata.Component.BOMRef}}})

			if b2.Metadata.Component.Components != nil {
				utils.Merge(b1.Components, b2.Metadata.Component.Components)
			}
		}
	}

	utils.Merge(b1.Components, b2.Components)
	utils.Merge(b1.Services, b2.Services)
	utils.Merge(b1.ExternalReferences, b2.ExternalReferences)
	utils.Merge(b1.Compositions, b2.Compositions)
	utils.Merge(b1.Properties, b2.Properties)
	utils.Merge(b1.Annotations, b2.Annotations)

	missing := utils.MergeDependenciesMissing(b1.Dependencies, b2.Dependencies)
	if len(missing) > 0 {
		for _, dep := range *b1.Dependencies {
			if dep.Ref == b2.Metadata.Component.BOMRef {
				*dep.Dependencies = append(*dep.Dependencies, missing...)
				break
			}
		}
	}

}
