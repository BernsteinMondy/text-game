package main

const (
	roomLookAroundKitchen  = "ты находишься на кухне, %s, надо %s. можно пройти - %s"
	roomLookAroundCorridor = "%s можно пройти - %s"
	roomLookAroundMyRoom   = "%s. можно пройти - %s"
)

const (
	kitchenActionNotReady = "собрать рюкзак и идти в универ"
	kitchenActionReady    = "идти в универ"
)

const (
	roomDescKitchen  = "кухня, ничего интересного. можно пройти - %s"
	roomDescCorridor = "ничего интересного. можно пройти - %s"
	roomDescMyRoom   = "ты в своей комнате. можно пройти - %s"
	roomDescOutside  = "на улице весна. можно пройти - %s"
)

const (
	roomNameKitchen  = "кухня"
	roomNameOutside  = "улица"
	roomNameCorridor = "коридор"
	roomNameMyRoom   = "комната"
)

const (
	actionLookAround = "осмотреться"
	actionGo         = "идти"
	actionPickUp     = "взять"
	actionApply      = "применить"
	actionWear       = "надеть"
)

const (
	responseUnknownCommand = "неизвестная команда"

	responseNoItemInInventory  = "нет предмета в инвентаре - %s"
	responseNoSpaceInInventory = "некуда класть"

	responseItemCantBeApplied = "не к чему применить"
	responseItemCantBeWorn    = "не может быть надето"
	responseItemPickedUp      = "предмет добавлен в инвентарь: %s"
	responseNoItem            = "нет такого"

	responseNoWayToRoom = "нет пути в %s"
	responseEmptyRoom   = "пустая комната"

	responseDoorIsLocked   = "дверь закрыта"
	responseDoorIsUnLocked = "дверь открыта"
)
