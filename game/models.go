package main

// World - is a general top struct, describes the whole game world.
// World consists of several amount of Room.
type World struct {
	Rooms map[string]*Room
}

// Room - is a part of the World, describes a part of the World.
// Each Room has it`s own condition, items, exits and Door to another Room.
type Room struct {
	Actions           map[string]ActionFunc
	Doors             map[string]*Door
	Items             []*Item
	Exits             []string
	Name              string
	Description       string
	LookAroundMessage string
}

// ActionFunc describes an action, which can be done inside a current Room.
// ActionFunc can change Room condition and items.
// ActionFunc can change Player condition, like Player`s current room, inventory condition or worn items.
type ActionFunc func(args ...string) string

// Door connects to Room between each other.
type Door struct {
	IsLocked bool
}

type Item struct {
	WhatAppliesTo map[string]struct{}
	Name          string
	Location      string
	CanBeWorn     bool
	Apply         ApplyFunc
}

// ApplyFunc describes an action, which can be proceeded when an Item is applied by a Player.
// ApplyFunc can change Room condition and items.
// ApplyFunc can change Player condition, like Player`s current room, inventory condition or worn items.
type ApplyFunc func(player *Player) string

// Player represents user in the World.
type Player struct {
	Inventory     map[string]*Item
	ItemsWorn     map[string]struct{}
	CurrentRoom   *Room
	InventorySize int
}
