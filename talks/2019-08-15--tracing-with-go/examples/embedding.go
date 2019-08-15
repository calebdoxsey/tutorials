package main

import "fmt"

func main() {
	type Knife interface {
		Cut()
	}
	type CanOpener interface {
		OpenCan()
	}

	{
		var knife Knife
		_, isKnife := knife.(Knife)
		_, isCanOpener := knife.(CanOpener)
		fmt.Println("knife?", isKnife, "can-opener?", isCanOpener)
	}

	{
		var canOpener CanOpener
		_, isKnife := canOpener.(Knife)
		_, isCanOpener := canOpener.(CanOpener)
		fmt.Println("knife?", isKnife, "can-opener?", isCanOpener)
	}

	{
		// not allowed:
		// var swissArmyKnife struct {
		// 	Knife
		// 	CanOpener
		// }
		// _, isKnife := swissArmyKnife.(Knife)
		// _, isCanOpener := swissArmyKnife.(CanOpener)
		// fmt.Println("knife?", isKnife, "can-opener?", isCanOpener)
	}

	{
		// allowed:
		var swissArmyKnife interface{} = struct {
			Knife
			CanOpener
		}{}
		_, isKnife := swissArmyKnife.(Knife)
		_, isCanOpener := swissArmyKnife.(CanOpener)
		fmt.Println("knife?", isKnife, "can-opener?", isCanOpener)
	}
}
