package flatmerge

import (
	"encoding/xml"
	"testing"

	"github.com/fnxpt/cyclonedx-merge/utils"

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

	sbom := utils.NewBOM()
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

	assert.Equal(t, 2, len(*sbom.Dependencies))
	//TODO: SORTING
	// assert.Equal(t, "root", (*sbom.Dependencies)[0].Ref)
	// assert.Equal(t, "topLib:1.0", (*(*sbom.Dependencies)[0].Dependencies)[0])
	// assert.Equal(t, "topLib:1.0", (*sbom.Dependencies)[1].Ref)
	// assert.Equal(t, "libA:1.0", (*(*sbom.Dependencies)[1].Dependencies)[0])
	// assert.Equal(t, "libB:2.0", (*(*sbom.Dependencies)[1].Dependencies)[1])
}

func TestMergeTwoSBOM(t *testing.T) {
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
				Ref:          "topLib:1.0",
				Dependencies: &[]string{"libA:1.0"},
			},
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

	secondObject := cyclonedx.BOM{
		XMLName:      xml.Name{},
		XMLNS:        "",
		JSONSchema:   "",
		BOMFormat:    "",
		SpecVersion:  0,
		SerialNumber: "",
		Version:      0,
		Metadata: &cyclonedx.Metadata{
			Component: &cyclonedx.Component{BOMRef: "topLibB:2.0", Name: "topLibB", Version: "2.0"},
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
			{
				BOMRef:  "libC:3.0",
				Name:    "libC",
				Version: "3.0",
			},
		},
		Services:           &[]cyclonedx.Service{},
		ExternalReferences: &[]cyclonedx.ExternalReference{},
		Dependencies: &[]cyclonedx.Dependency{
			{
				Ref:          "topLibB:2.0",
				Dependencies: &[]string{"libA:1.0"},
			},
			{
				Ref:          "libA:1.0",
				Dependencies: &[]string{"libB:2.0", "libC:3.0"},
			},
		},
		Compositions:    &[]cyclonedx.Composition{},
		Properties:      &[]cyclonedx.Property{},
		Vulnerabilities: &[]cyclonedx.Vulnerability{},
		Annotations:     &[]cyclonedx.Annotation{},
	}

	sbom := utils.NewBOM()
	MergeSBOM(sbom, &firstObject)
	MergeSBOM(sbom, &secondObject)

	assert.NotNil(t, sbom)
	assert.Equal(t, "root", sbom.Metadata.Component.BOMRef)
	assert.Equal(t, "root", sbom.Metadata.Component.Name)
	assert.Equal(t, "", sbom.Metadata.Component.Version)

	assert.Equal(t, 5, len(*sbom.Components))
	assert.Equal(t, "topLib:1.0", (*sbom.Components)[0].BOMRef)
	assert.Equal(t, "topLib", (*sbom.Components)[0].Name)
	assert.Equal(t, "1.0", (*sbom.Components)[0].Version)
	assert.Equal(t, "libA:1.0", (*sbom.Components)[1].BOMRef)
	assert.Equal(t, "libA", (*sbom.Components)[1].Name)
	assert.Equal(t, "1.0", (*sbom.Components)[1].Version)
	assert.Equal(t, "libB:2.0", (*sbom.Components)[2].BOMRef)
	assert.Equal(t, "libB", (*sbom.Components)[2].Name)
	assert.Equal(t, "2.0", (*sbom.Components)[2].Version)
	assert.Equal(t, "topLibB:2.0", (*sbom.Components)[3].BOMRef)
	assert.Equal(t, "topLibB", (*sbom.Components)[3].Name)
	assert.Equal(t, "2.0", (*sbom.Components)[3].Version)
	assert.Equal(t, "libC:3.0", (*sbom.Components)[4].BOMRef)
	assert.Equal(t, "libC", (*sbom.Components)[4].Name)
	assert.Equal(t, "3.0", (*sbom.Components)[4].Version)

	assert.Equal(t, 3, len(*sbom.Dependencies))
	//TODO: SORTING
	// assert.Equal(t, "root", (*sbom.Dependencies)[0].Ref)
	// assert.Equal(t, "topLib:1.0", (*(*sbom.Dependencies)[0].Dependencies)[0])
	// assert.Equal(t, "topLibB:2.0", (*(*sbom.Dependencies)[0].Dependencies)[1])
	// assert.Equal(t, "topLib:1.0", (*sbom.Dependencies)[1].Ref)
	// assert.Equal(t, "libA:1.0", (*(*sbom.Dependencies)[1].Dependencies)[0])
	// assert.Equal(t, "libB:2.0", (*(*sbom.Dependencies)[1].Dependencies)[1])
	// assert.Equal(t, "topLibB:2.0", (*sbom.Dependencies)[2].Ref)
	// assert.Equal(t, "libA:1.0", (*(*sbom.Dependencies)[2].Dependencies)[0])
	// assert.Equal(t, "libB:2.0", (*(*sbom.Dependencies)[2].Dependencies)[1])
	// assert.Equal(t, "libC:3.0", (*(*sbom.Dependencies)[2].Dependencies)[2])
}
