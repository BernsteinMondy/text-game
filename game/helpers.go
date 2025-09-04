package main

import (
	"slices"
	"sort"
	"strings"
)

// isItemInRoom checks if the target item exists in Room.
func isItemInRoom(room *Room, itemName string) bool {
	itemNames := make([]string, 0, len(room.Items))
	for _, item := range room.Items {
		itemNames = append(itemNames, item.Name)
	}

	return slices.Contains(itemNames, itemName)
}

// deleteItemFromRoom deletes target item from Room.
func deleteItemFromRoom(room *Room, itemName string) {
	result := make([]*Item, 0, len(room.Items))

	for _, item := range room.Items {
		if item != nil && item.Name != itemName {
			result = append(result, item)
		}
	}

	room.Items = result
}

// buildStringAboutItemsInRoom builds a complex string in format
// location1 : item-1, item-2, location2 : item-3, item-4.
//
// buildStringAboutItemsInRoom returns responseEmptyRoom if no items in a Room.
func buildStringAboutItemsInRoom(items []*Item) string {
	locationMap := make(map[string][]string)

	if len(items) == 0 {
		return responseEmptyRoom
	}

	for _, item := range items {
		if item.Location == "" {
			continue
		}
		locationMap[item.Location] = append(locationMap[item.Location], item.Name)
	}

	var result strings.Builder

	locations := make([]string, 0, len(locationMap))
	for location := range locationMap {
		locations = append(locations, location)
	}
	sort.Strings(locations)

	for i, location := range locations {
		items := locationMap[location]
		sort.Strings(items)

		if i > 0 {
			result.WriteString(", ")
		}

		result.WriteString(location)
		result.WriteString(": ")
		result.WriteString(strings.Join(items, ", "))
	}

	return result.String()
}
