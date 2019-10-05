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

// for regular entities these subscriptions might be in the Entity struct
var weaponSwingSubscription *Subscription

// cancel the swing before swinging again, or interrupted by something like picking up another weapon
func cancelSwingWeapon() {

	if weaponSwingSubscription != nil {
		weaponSwingSubscription.cancel()
	}
	weaponSwingSubscription = nil
}

func swingWeapon() {

	cancelSwingWeapon()

	subscription := entityMovementEvent.subscribe()
	weaponSwingSubscription = subscription
	go func() {
		defer cancelSwingWeapon()
		for {
			subscription.wait()
			if subscription.cancelled {
				break
			}
			if weaponRotationDone() {
				break
			}
			rotateWeapon()
		}
		returnWeapon()
	}()
}
