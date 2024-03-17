package query

import (
	"testing"

	"github.com/ermes-labs/api-go/infrastructure"
)

var infra infrastructure.Infrastructure = infrastructure.Infrastructure{
	Areas: []infrastructure.Area{
		{
			Node: infrastructure.Node{
				AreaName: "Area1",
				Tags: map[string]string{
					"tag1": "value1",
				},
			},
		},
		{
			Node: infrastructure.Node{
				AreaName: "Area2",
				Tags: map[string]string{
					"tag2": "value2",
				},
			},
			Areas: []infrastructure.Area{
				{
					Node: infrastructure.Node{
						AreaName: "Area2.1",
						Tags: map[string]string{
							"tag3": "value3",
						},
					},
				},
			},
		},
	},
}

var areasMap = map[string]*infrastructure.Area{
	infra.Areas[0].AreaName:          &infra.Areas[0],
	infra.Areas[1].AreaName:          &infra.Areas[1],
	infra.Areas[1].Areas[0].AreaName: &infra.Areas[1].Areas[0],
}

func TestParser(t *testing.T) {
	query, err := parser.ParseString("", "{ #Area1 + #Area2 - #Area3 + { #Area4 - #Area5 } }")
	if err != nil {
		t.Errorf("Error parsing query: %v", err)
		return
	}

	areas := collectAreas(areasMap, query.Set)

	if len(areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(areas))
	}

	if areas[0].AreaName != "Area1" {
		t.Errorf("Expected Area1, got %s", areas[0].AreaName)
	}

	if areas[1].AreaName != "Area2" {
		t.Errorf("Expected Area2, got %s", areas[1].AreaName)
	}

	query, err = parser.ParseString("", "Area2")

	if err != nil {
		t.Errorf("Error parsing query: %v", err)
	}

	areas = collectAreas(areasMap, query.Set)

	if len(areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(areas))
	} else {
		if areas[0].AreaName != "Area2" {
			t.Errorf("Expected Area2, got %s", areas[0].AreaName)
		}

		if areas[1].AreaName != "Area2.1" {
			t.Errorf("Expected Area2.1, got %s", areas[1].AreaName)
		}
	}

	query, err = parser.ParseString("", "*")

	if err != nil {
		t.Errorf("Error parsing query: %v", err)
		return
	}

	areas = collectAreas(areasMap, query.Set)

	if len(areas) != 3 {
		t.Errorf("Expected 3 areas, got %d", len(areas))
	}

	query, err = parser.ParseString("", "*(tag:tag1=value1)")

	if err != nil {
		t.Errorf("Error parsing query: %v", err)
		return
	}

	areas = collectAreas(areasMap, query.Set)

	if len(areas) != 1 {
		t.Errorf("Expected 1 area, got %d", len(areas))
	} else if areas[0].AreaName != "Area1" {
		t.Errorf("Expected Area1, got %s", areas[0].AreaName)
	}
}
