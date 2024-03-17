package query

import (
	"strings"

	"github.com/ermes-labs/api-go/infrastructure"
)

func CollectAreas(infra *infrastructure.Infrastructure, areasMap map[string]*infrastructure.Area, query string) ([]*infrastructure.Area, error) {
	q, err := ParseQuery(query)
	if err != nil {
		return nil, err
	}

	return collectAreas(areasMap, q.Set), nil
}

func ParseQuery(input string) (*Query, error) {
	query, err := parser.ParseString("query", input)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func collectAreas(
	areasMap map[string]*infrastructure.Area,
	set *Set,
) []*infrastructure.Area {
	var areas []*infrastructure.Area

	if set.AreaId != nil {
		// Handle the '#AreaId' case.
		if area, found := areasMap[strings.TrimPrefix(*set.AreaId, "#")]; found {
			areas = append(areas, area)
		}
	} else if set.NodesInArea != nil {
		if set.NodesInArea.AreaId == "*" {
			// Handle the '*' case.
			for _, area := range areasMap {
				if checkFilters(area, set.NodesInArea.Filters...) {
					areas = append(areas, area)
				}
			}
		} else {
			// Handle the 'NodesInArea' case.
			if area, found := areasMap[set.NodesInArea.AreaId]; found {
				areas = appendTree(areas, area, set.NodesInArea.Filters...)
			}
		}
	} else if set.Sets != nil {
		// Handle the 'Sets' case.
		areas = append(areas, collectAreas(areasMap, set.Sets.InitialSet)...)
		for _, opSet := range set.Sets.NextSets {
			additionalAreas := collectAreas(areasMap, opSet.Set)
			if opSet.Op == "+" {
				areas = union(areas, additionalAreas)
			} else if opSet.Op == "-" {
				areas = difference(areas, additionalAreas)
			}
		}
	}

	return areas
}

func checkFilters(area *infrastructure.Area, filters ...*Filter) bool {
	for _, filter := range filters {
		// if filter.Level != nil {
		// 	// TODO: Implement the level filter.
		// }

		if filter.Tag != nil {
			if area.Tags == nil {
				return false
			}

			if value, found := area.Tags[filter.Tag.Key]; !found || value != filter.Tag.Value {
				return false
			}
		}
	}

	return true
}

func appendTree(areas []*infrastructure.Area, area *infrastructure.Area, filters ...*Filter) []*infrastructure.Area {
	areas = append(areas, area)
	if area.Areas != nil {
		for _, subArea := range area.Areas {
			if checkFilters(&subArea, filters...) {
				areas = appendTree(areas, &subArea)
			}
		}
	}

	return areas
}

func union(a, b []*infrastructure.Area) []*infrastructure.Area {
	m := make(map[string]bool)
	var result []*infrastructure.Area

	for _, areaA := range a {
		m[areaA.AreaName] = true
		result = append(result, areaA)
	}

	for _, areaB := range b {
		if _, found := m[areaB.AreaName]; !found {
			result = append(result, areaB)
		}
	}

	return result
}

func difference(a, b []*infrastructure.Area) []*infrastructure.Area {
	m := make(map[string]bool)
	var result []*infrastructure.Area

	for _, item := range b {
		m[item.AreaName] = true
	}

	for _, item := range a {
		if _, found := m[item.AreaName]; !found {
			result = append(result, item)
		}
	}

	return result
}
