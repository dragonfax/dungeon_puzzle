package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var ENEMY_ORDER []string = []string{
	"imp",
	"goblin",
	"masked_orc",
	"skelet",
	"muddy",
	"zombie",
	"ice_zombie",
	"orc_warrior",
	"orc_shaman",
	"wogol",
	"big_zombie",
	"ogre",
	"big_damon",
}

var monsters []*PlacedEntity

type WillMove struct {
	X, Y int
}

func drawMonsterWillMove(r *sdl.Renderer) {
	for _, monster := range monsters {
		if willMove, ok := monster.Data.(WillMove); monster.Data != nil && ok {
			// full image
			arrowRect := sdl.Rect{0, 0, 16, 16}

			// default angle is facing down
			angle := 0.0
			if willMove.X-monster.X > 0 {
				angle = 180
			} else if willMove.Y-monster.Y > 0 {
				angle = 280
			} else if monster.Y-willMove.Y > 0 {
				angle = 90
			}

			// draw half way between the 2 points.
			tgtRect := sdl.Rect{
				int32(willMove.X+monster.X) * PIXELS_PER_CELL / 2,
				int32(willMove.Y+monster.Y) * PIXELS_PER_CELL / 2,
				16,
				16}
			err := r.CopyEx(pixelTex, &arrowRect, &tgtRect, angle, nil, 0)
			if err != nil {
				panic(err)
			}
		}
	}
}

func monstersThink() {
	for _, monster := range monsters {

		willMove := WillMove{}

		// next space towards player
		willMove.X = character.X - monster.X
		if willMove.X > 1 {
			willMove.X = 1
		}
		if willMove.X < -1 {
			willMove.X = -1
		}
		willMove.Y = character.Y - monster.Y
		if willMove.Y > 1 {
			willMove.Y = 1
		}
		if willMove.Y < -1 {
			willMove.Y = -1
		}

		willMove.X = monster.X + willMove.X
		willMove.Y = monster.Y + willMove.Y

		// keep it if there is nothing there or the player there.
		if len(otherEntitiesAt(character, willMove.X, willMove.Y)) == 0 {
			fmt.Printf("monster will move to %v\n", willMove)
			monster.Data = willMove
		}
	}
}

func monstersMove() {
	for _, monster := range monsters {
		if willMove, ok := monster.Data.(WillMove); monster.Data != nil && ok {
			// if still empty
			if len(otherEntitiesAt(monster, willMove.X, willMove.Y)) == 0 {
				// move there
				monster.X = willMove.X
				monster.Y = willMove.Y
			} else if character.X == willMove.X && character.Y == willMove.Y {
				panic("end game")
			}
			// clear the willMove
			monster.Data = nil
		}
	}
}

func findEmptyPosition() (x, y int) {
	for x := 0; x < MAX_X; x++ {
		for y := 0; y < MAX_Y; y++ {
			occupied := entitiesAt(x, y)
			if len(occupied) == 0 {
				fmt.Printf("found empty cell at %d,%d\n", x, y)
				return x, y
			}
		}
	}
	panic("board full")
}

func spawnMonster() {
	x, y := findEmptyPosition()
	newMonster := &PlacedEntity{
		Sprite: spriteByName("skelet_idle_anim"),
		X:      x,
		Y:      y,
	}
	monsters = append(monsters, newMonster)
}
