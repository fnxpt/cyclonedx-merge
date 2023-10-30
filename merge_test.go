package main

import (
	"testing"

	cyclonedx "github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/assert"
)

// TestMergeComponents
func TestMergeComponents(t *testing.T) {
	firstObject := []cyclonedx.Component{
		{
			BOMRef:  "A:1",
			Name:    "A",
			Version: "1",
		},
		{
			BOMRef:  "B:2",
			Name:    "B",
			Version: "2",
		},
	}
	secondObject := []cyclonedx.Component{
		{
			BOMRef:  "A:1",
			Name:    "A",
			Version: "1",
		},
		{
			BOMRef:  "B:2",
			Name:    "B",
			Version: "2",
		},
	}

	Merge(&firstObject, &secondObject, "")
	assert.Equal(t, 2, len(firstObject))
	assert.Equal(t, "A:1", firstObject[0].BOMRef)
	assert.Equal(t, "A", firstObject[0].Name)
	assert.Equal(t, "1", firstObject[0].Version)
	assert.Equal(t, "B:2", firstObject[1].BOMRef)
	assert.Equal(t, "B", firstObject[1].Name)
	assert.Equal(t, "2", firstObject[1].Version)
}

func TestMergeComponentsWithDifferentObjects(t *testing.T) {
	firstObject := []cyclonedx.Component{
		{
			BOMRef:  "A:1",
			Name:    "A",
			Version: "1",
		},
		{
			BOMRef:  "B:2",
			Name:    "B",
			Version: "2",
		},
	}
	secondObject := []cyclonedx.Component{
		{
			BOMRef:  "A:1",
			Name:    "A",
			Version: "1",
		},
		{
			BOMRef:  "B:2",
			Name:    "B",
			Version: "3",
		},
	}

	Merge(&firstObject, &secondObject, "")
	assert.Equal(t, 2, len(firstObject))
	assert.Equal(t, "A:1", firstObject[0].BOMRef)
	assert.Equal(t, "A", firstObject[0].Name)
	assert.Equal(t, "1", firstObject[0].Version)
	assert.Equal(t, "B:2", firstObject[1].BOMRef)
	assert.Equal(t, "B", firstObject[1].Name)
	assert.Equal(t, "2", firstObject[1].Version)
}

// TestMergeComponentsNested
func TestMergeComponentsNested(t *testing.T) {
	firstObject := []cyclonedx.Component{
		{
			BOMRef:  "A:1",
			Name:    "A",
			Version: "1",
		},
		{
			BOMRef:  "B:2",
			Name:    "B",
			Version: "2",
		},
	}
	secondObject := []cyclonedx.Component{
		{
			BOMRef:  "A:1",
			Name:    "A",
			Version: "1",
		},
		{
			BOMRef:  "B:2",
			Name:    "B",
			Version: "2",
		},
	}

	Merge(&firstObject, &secondObject, "X|")
	assert.Equal(t, 4, len(firstObject))
	assert.Equal(t, "A:1", firstObject[0].BOMRef)
	assert.Equal(t, "A", firstObject[0].Name)
	assert.Equal(t, "1", firstObject[0].Version)
	assert.Equal(t, "B:2", firstObject[1].BOMRef)
	assert.Equal(t, "B", firstObject[1].Name)
	assert.Equal(t, "2", firstObject[1].Version)
	assert.Equal(t, "X|A:1", firstObject[2].BOMRef)
	assert.Equal(t, "A", firstObject[2].Name)
	assert.Equal(t, "1", firstObject[2].Version)
	assert.Equal(t, "X|B:2", firstObject[3].BOMRef)
	assert.Equal(t, "B", firstObject[3].Name)
	assert.Equal(t, "2", firstObject[3].Version)
}
