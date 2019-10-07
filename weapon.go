package main

import (
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
)

func gatherWeapons() []*Sprite {
	return spritesWithTag("weapon")
}

func placeWeapons(weapons []*Sprite) ([]*PlacedEntity, *resolv.Space) {
	space := resolv.NewSpace()

	placedWeapons := make([]*PlacedEntity, 0)
	for i := range weapons {
		weapon := weapons[i]
		placedWeapon := &PlacedEntity{
			Sprite: weapon,
			Shape:  resolv.NewRectangle(rand.Int31n(FIELD_WIDTH/UNIT_SIZE)*UNIT_SIZE, rand.Int31n(FIELD_WIDTH/UNIT_SIZE)*UNIT_SIZE, weapon.Frames[0].W, weapon.Frames[0].H),
		}
		placedWeapon.Shape.SetData(placedWeapon)
		space.Add(placedWeapon.Shape)
		placedWeapons = append(placedWeapons, placedWeapon)
	}
	return placedWeapons, space
}
