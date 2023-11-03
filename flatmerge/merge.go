package flatmerge

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

			utils.MergeDependenciesFromComponents(b1.Dependencies, &[]cyclonedx.Component{*b2.Metadata.Component}, "root")

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

	utils.MergeDependenciesFromComponents(b1.Dependencies, b2.Components, b2.Metadata.Component.BOMRef)
}
