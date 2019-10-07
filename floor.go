package main

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

func showFloor(tick int, r *sdl.Renderer, floor [][]*Sprite) {
	for y := 0; y < len(floor); y++ {
		for x := 0; x < len(floor[y]); x++ {
			sprite := floor[y][x]

			drawSpriteAt(tick, r, sprite, int32(x)*UNIT_SIZE, int32(y)*UNIT_SIZE, 0)
		}
	}
}

var floorTiles []*Sprite

func chooseRandomFloorSprite() *Sprite {
	if floorTiles == nil {
		floorTiles = spritesWithTag("floor")
		fmt.Printf("found %d floor tiles\n", len(floorTiles))
		fmt.Printf("%v\n%v\n", floorTiles[0], floorTiles[1])
	}
	n := rand.Intn(len(floorTiles))
	// fmt.Printf("choosing floor tile %d %s\n", n, floorTiles[n].Name)
	return floorTiles[n]
}

func generateFloor() [][]*Sprite {
	floor := make([][]*Sprite, 10)
	for y := 0; y < 10; y++ {
		floor[y] = make([]*Sprite, 10)
		for x := 0; x < 10; x++ {
			floor[y][x] = chooseRandomFloorSprite()
		}
	}
	return floor
}
