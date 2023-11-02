package merge

import (
	"encoding/xml"
	"testing"

	cyclonedx "github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/assert"
)

func TestMergeOneSBOM(t *testing.T) {
	firstObject := cyclonedx.BOM{
		XMLName:      xml.Name{},
		XMLNS:        "",
		JSONSchema:   "",
		BOMFormat:    "",
		SpecVersion:  0,
		SerialNumber: "",
		Version:      0,
		Metadata: &cyclonedx.Metadata{
			Component: &cyclonedx.Component{BOMRef: "topLib:1.0", Name: "topLib", Version: "1.0"},
		},
		Components: &[]cyclonedx.Component{
			{
				BOMRef:  "libA:1.0",
				Name:    "libA",
				Version: "1.0",
			},
			{
				BOMRef:  "libB:2.0",
				Name:    "libB",
				Version: "2.0",
			},
		},
		Services:           &[]cyclonedx.Service{},
		ExternalReferences: &[]cyclonedx.ExternalReference{},
		Dependencies: &[]cyclonedx.Dependency{
			{
				Ref:          "libA:1.0",
				Dependencies: &[]string{"libB:2.0"},
			},
		},
		Compositions:    &[]cyclonedx.Composition{},
		Properties:      &[]cyclonedx.Property{},
		Vulnerabilities: &[]cyclonedx.Vulnerability{},
		Annotations:     &[]cyclonedx.Annotation{},
	}

	sbom := NewBOM()
	MergeSBOM(sbom, &firstObject)

	assert.NotNil(t, sbom)
	assert.Equal(t, "root", sbom.Metadata.Component.BOMRef)
	assert.Equal(t, "root", sbom.Metadata.Component.Name)
	assert.Equal(t, "", sbom.Metadata.Component.Version)

	assert.Equal(t, 3, len(*sbom.Components))
	assert.Equal(t, "topLib:1.0", (*sbom.Components)[0].BOMRef)
	assert.Equal(t, "topLib", (*sbom.Components)[0].Name)
	assert.Equal(t, "1.0", (*sbom.Components)[0].Version)
	assert.Equal(t, "libA:1.0", (*sbom.Components)[1].BOMRef)
	assert.Equal(t, "libA", (*sbom.Components)[1].Name)
	assert.Equal(t, "1.0", (*sbom.Components)[1].Version)
	assert.Equal(t, "libB:2.0", (*sbom.Components)[2].BOMRef)
	assert.Equal(t, "libB", (*sbom.Components)[2].Name)
	assert.Equal(t, "2.0", (*sbom.Components)[2].Version)

	assert.Equal(t, 4, len(*sbom.Dependencies))
	assert.Equal(t, "root", (*sbom.Dependencies)[0].Ref)
	assert.Equal(t, "topLib:1.0", (*(*sbom.Dependencies)[0].Dependencies)[0])
	assert.Equal(t, "topLib:1.0", (*sbom.Dependencies)[1].Ref)
	assert.Equal(t, "libA:1.0", (*(*sbom.Dependencies)[1].Dependencies)[0])
	assert.Equal(t, "libA:1.0", (*sbom.Dependencies)[2].Ref)
	assert.Equal(t, "libB:2.0", (*(*sbom.Dependencies)[2].Dependencies)[0])
	assert.Equal(t, "libB:2.0", (*sbom.Dependencies)[3].Ref)
	assert.Nil(t, (*(*sbom.Dependencies)[2].Dependencies))
}

// func TestMergeTwoSBOM(t *testing.T) {
// 	firstObject := []cyclonedx.BOM{
// 		{
// 			BOMRef:  "libA:1.0",
// 			Name:    "libA",
// 			Version: "1.0",
// 		},
// 		{
// 			BOMRef:  "libB:2.0",
// 			Name:    "libB",
// 			Version: "2.0",
// 		},
// 	}
// 	secondObject := []cyclonedx.Component{
// 		{
// 			BOMRef:  "libA:1.0",
// 			Name:    "libA",
// 			Version: "1.0",
// 		},
// 		{
// 			BOMRef:  "libB:2.0",
// 			Name:    "libB",
// 			Version: "2.0",
// 		},
// 	}
// }

// TestMergeComponents
func TestMergeComponents(t *testing.T) {
	firstObject := []cyclonedx.Component{
		{
			BOMRef:  "libA:1.0",
			Name:    "libA",
			Version: "1.0",
		},
		{
			BOMRef:  "libB:2.0",
			Name:    "libB",
			Version: "2.0",
		},
	}
	secondObject := []cyclonedx.Component{
		{
			BOMRef:  "libA:1.0",
			Name:    "libA",
			Version: "1.0",
		},
		{
			BOMRef:  "libB:2.0",
			Name:    "libB",
			Version: "2.0",
		},
	}

	merge(&firstObject, &secondObject)
	assert.Equal(t, 2, len(firstObject))
	assert.Equal(t, "libA:1.0", firstObject[0].BOMRef)
	assert.Equal(t, "libA", firstObject[0].Name)
	assert.Equal(t, "1.0", firstObject[0].Version)
	assert.Equal(t, "libB:2.0", firstObject[1].BOMRef)
	assert.Equal(t, "libB", firstObject[1].Name)
	assert.Equal(t, "2.0", firstObject[1].Version)
}

func TestMergeComponentsWithDifferentObjects(t *testing.T) {
	firstObject := []cyclonedx.Component{
		{
			BOMRef:  "libA:1.0",
			Name:    "libA",
			Version: "1.0",
		},
		{
			BOMRef:  "libB:2.0",
			Name:    "libB",
			Version: "2.0",
		},
	}
	secondObject := []cyclonedx.Component{
		{
			BOMRef:  "libA:1.0",
			Name:    "libA",
			Version: "1.0",
		},
		{
			BOMRef:  "libB:2.0",
			Name:    "libB",
			Version: "3.0",
		},
	}

	merge(&firstObject, &secondObject)
	assert.Equal(t, 2, len(firstObject))
	assert.Equal(t, "libA:1.0", firstObject[0].BOMRef)
	assert.Equal(t, "libA", firstObject[0].Name)
	assert.Equal(t, "1.0", firstObject[0].Version)
	assert.Equal(t, "libB:2.0", firstObject[1].BOMRef)
	assert.Equal(t, "libB", firstObject[1].Name)
	assert.Equal(t, "2.0", firstObject[1].Version)
}

// func TestMergeComponentsWithMultipleVersions(t *testing.T) {
// 	firstObject := []cyclonedx.Component{
// 		{
// 			BOMRef:  "libA:1.0",
// 			Name:    "libA",
// 			Version: "1.0",
// 		},
// 		{
// 			BOMRef:  "libB:2.0",
// 			Name:    "libB",
// 			Version: "2.0",
// 		},
// 	}
// 	secondObject := []cyclonedx.Component{
// 		{
// 			BOMRef:  "libA:1.0",
// 			Name:    "libA",
// 			Version: "1.0",
// 		},
// 		{
// 			BOMRef:  "libB:3.0",
// 			Name:    "libB",
// 			Version: "3.0",
// 		},
// 	}

// 	merge(&firstObject, &secondObject)
// 	assert.Equal(t, 1, len(firstObject))
// 	assert.Equal(t, "libA:1.0", firstObject[0].BOMRef)
// 	assert.Equal(t, "libA", firstObject[0].Name)
// 	assert.Equal(t, "1.0", firstObject[0].Version)
// }

// func TestMergeDependencies(t *testing.T) {
// 	dependencies := []cyclonedx.Dependency{
// 		{
// 			Ref:          "A",
// 			Dependencies: &[]string{"B"},
// 		},
// 		{
// 			Ref:          "B",
// 			Dependencies: &[]string{"C", "D"},
// 		},
// 	}
// 	MergeMap(&dependencies, "", false)
// }
