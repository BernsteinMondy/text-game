package main

import (
	"fmt"
	"slices"
	"strings"
)

var (
	world = &World{
		Rooms: make(map[string]*Room),
	}
	player = &Player{}
)

func main() {
	initGame()
}

func initGame() {
	// Setup world.
	setupWorld()

	// Player starts from Kitchen room.
	player.CurrentRoom = world.Rooms[roomNameKitchen]
	player.InventorySize = 0
	player.Inventory = make(map[string]*Item)
	player.ItemsWorn = make(map[string]struct{})
}

// setupWorld sets up a new World.
func setupWorld() {
	// Setup rooms.

	// Setup doors.
	doorBetweenCorridorAndOutside := &Door{
		IsLocked: true,
	}

	corridorDoors := map[string]*Door{
		roomNameOutside: doorBetweenCorridorAndOutside,
	}
	outsideDoors := map[string]*Door{
		roomNameCorridor: doorBetweenCorridorAndOutside,
	}

	// Kitchen room.
	kitchenItems := []*Item{
		{
			Name:          "чай",
			WhatAppliesTo: map[string]struct{}{},
			Location:      "на столе",
		},
	}
	kitchenExits := []string{
		roomNameCorridor,
	}
	world.Rooms[roomNameKitchen] = &Room{
		Name:              roomNameKitchen,
		Description:       roomDescKitchen,
		LookAroundMessage: roomLookAroundKitchen,
		Doors:             make(map[string]*Door),
		Items:             kitchenItems,
		Exits:             kitchenExits,
	}

	// Corridor room.
	var corridorItems []*Item
	corridorExits := []string{
		roomNameKitchen,
		roomNameMyRoom,
		roomNameOutside,
	}
	world.Rooms[roomNameCorridor] = &Room{
		Name:              roomNameCorridor,
		Description:       roomDescCorridor,
		LookAroundMessage: roomLookAroundCorridor,
		Doors:             corridorDoors,
		Items:             corridorItems,
		Exits:             corridorExits,
	}

	// My Room.
	myRoomItems := []*Item{
		{
			Name: "ключи",
			WhatAppliesTo: map[string]struct{}{
				"дверь": {},
			},
			Location:  "на столе",
			CanBeWorn: false,
			Apply: func(player *Player) string {
				door, ok := player.CurrentRoom.Doors[roomNameOutside]
				if !ok {
					return responseItemCantBeApplied
				}

				door.IsLocked = !door.IsLocked
				if door.IsLocked {
					return responseDoorIsLocked
				} else {
					return responseDoorIsUnLocked
				}
			},
		},
		{
			Name:          "конспекты",
			WhatAppliesTo: map[string]struct{}{},
			Location:      "на столе",
			CanBeWorn:     false,
		},
		{
			Name:          "рюкзак",
			WhatAppliesTo: map[string]struct{}{},
			Location:      "на стуле",
			CanBeWorn:     true,
			Apply: func(player *Player) string {
				_, ok := player.ItemsWorn["рюкзак"]
				if !ok {
					player.InventorySize += 32
					player.ItemsWorn["рюкзак"] = struct{}{}

					return fmt.Sprintf(responseItemIsWorn, "рюкзак")
				}

				player.InventorySize -= 32
				delete(player.ItemsWorn, "рюкзак")

				return fmt.Sprintf(responseItemIsRemoved, "рюкзак")
			},
		}}
	myRoomExits := []string{
		roomNameCorridor,
	}
	world.Rooms[roomNameMyRoom] = &Room{
		Name:              roomNameMyRoom,
		Description:       roomDescMyRoom,
		LookAroundMessage: roomLookAroundMyRoom,
		Doors:             make(map[string]*Door),
		Items:             myRoomItems,
		Exits:             myRoomExits,
	}

	// Outside room.
	var outsideItems []*Item
	outsideExits := []string{
		"домой",
	}
	world.Rooms[roomNameOutside] = &Room{
		Name:        roomNameOutside,
		Description: roomDescOutside,
		Doors:       outsideDoors,
		Items:       outsideItems,
		Exits:       outsideExits,
	}

	// Setup actions for each Room in World.
	for _, room := range world.Rooms {
		setupRoomActions(room)
	}
}

// setupRoomActions setup Actions field (map[string]ActionFunc) dynamically for each Room.
func setupRoomActions(room *Room) {
	room.Actions = map[string]ActionFunc{
		actionLookAround: func(args ...string) string {
			if len(args) != 1 {
				return responseUnknownCommand
			}

			var resp string
			itemsResp := buildStringAboutItemsInRoom(room.Items)

			// If Player is located in kitchen, should build a dynamic message, depending on Player condition.
			if room.Name == roomNameKitchen {
				resp = buildResponseForKitchen(room.LookAroundMessage, itemsResp, strings.Join(room.Exits, ", "))
			} else {
				resp = fmt.Sprintf(room.LookAroundMessage, itemsResp, strings.Join(room.Exits, ", "))
			}

			return resp
		},
		actionGo: func(args ...string) string {
			if len(args) != 2 {
				return responseUnknownCommand
			}

			roomToGoName := args[1]
			roomToGo, ok := world.Rooms[roomToGoName]
			// Check if no exit from current room or room does not exist at all.
			if !slices.Contains(player.CurrentRoom.Exits, roomToGoName) || !ok {
				resp := fmt.Sprintf(responseNoWayToRoom, roomToGoName)
				return resp
			}

			door, exists := player.CurrentRoom.Doors[roomToGoName]
			if exists && door.IsLocked {
				return responseDoorIsLocked
			}

			player.CurrentRoom = roomToGo

			resp := fmt.Sprintf(roomToGo.Description, strings.Join(roomToGo.Exits, ", "))
			return resp
		},
		actionPickUp: func(args ...string) string {
			if len(args) != 2 {
				return responseUnknownCommand
			}

			itemToPickUp := args[1]
			itemInRoom := isItemInRoom(player.CurrentRoom, itemToPickUp)
			if !itemInRoom {
				return responseNoItem
			}

			var it *Item
			for _, item := range player.CurrentRoom.Items {
				if item.Name == itemToPickUp {
					it = item
				}
			}

			if player.InventorySize == len(player.Inventory) {
				return responseNoSpaceInInventory
			}

			// Add item to the player`s inventory
			player.Inventory[itemToPickUp] = it
			// Remove item from room.
			deleteItemFromRoom(player.CurrentRoom, itemToPickUp)

			resp := fmt.Sprintf(responseItemPickedUp, itemToPickUp)
			return resp
		},
		actionApply: func(args ...string) string {
			if len(args) != 3 {
				return responseUnknownCommand
			}

			var resp string
			itemToApply, onWhatToApply := args[1], args[2]

			item, ok := player.Inventory[itemToApply]
			if !ok {
				resp = fmt.Sprintf(responseNoItemInInventory, itemToApply)
				return resp
			}

			_, applies := item.WhatAppliesTo[onWhatToApply]
			if !applies {
				return responseItemCantBeApplied
			}

			resp = item.Apply(player)
			return resp
		},
		actionWear: func(args ...string) string {
			if len(args) != 2 {
				return responseUnknownCommand
			}

			itemToWear := args[1]
			itemInRoom := isItemInRoom(player.CurrentRoom, itemToWear)
			if !itemInRoom {
				return responseNoItem
			}

			var it *Item
			for _, item := range player.CurrentRoom.Items {
				if item.Name == itemToWear && !item.CanBeWorn {
					return responseItemCantBeWorn
				} else {
					it = item
				}
			}

			// Item was taken by a Player, so it should be deleted from Room.
			deleteItemFromRoom(player.CurrentRoom, itemToWear)

			// Apply Item.
			resp := it.Apply(player)
			return resp
		},
	}
}

// handleCommand handles player`s command according to command params length.
func handleCommand(command string) string {
	params := strings.Split(command, " ")
	if len(params) > 3 || len(params) < 1 {
		return responseUnknownCommand
	}

	// Find action
	actionParam := params[0]
	action, ok := player.CurrentRoom.Actions[actionParam]
	if !ok {
		return responseUnknownCommand
	}

	actionResult := action(params...)
	return actionResult
}
