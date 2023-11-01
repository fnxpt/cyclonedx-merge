package merge

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

	merge(&firstObject, &secondObject)
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

	merge(&firstObject, &secondObject)
	assert.Equal(t, 2, len(firstObject))
	assert.Equal(t, "A:1", firstObject[0].BOMRef)
	assert.Equal(t, "A", firstObject[0].Name)
	assert.Equal(t, "1", firstObject[0].Version)
	assert.Equal(t, "B:2", firstObject[1].BOMRef)
	assert.Equal(t, "B", firstObject[1].Name)
	assert.Equal(t, "2", firstObject[1].Version)
}

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
