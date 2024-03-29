package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

var ENEMY_ORDER []string = []string{
	"imp_idle_anim",
	"goblin_idle_anim",
	"masked_orc_idle_anim",
	"skelet_idle_anim",
	"muddy_idle_anim",
	"zombie_idle_anim",
	"ice_zombie_idle_anim",
	"orc_warrior_idle_anim",
	"orc_shaman_idle_anim",
	"wogol_idle_anim",
	"big_zombie_idle_anim",
	"ogre_idle_anim",
	"big_demon_idle_anim",
}
var ENEMY_ORDER_SPRITES []*Sprite

func loadEnemyOrderSprites() {
	ENEMY_ORDER_SPRITES = make([]*Sprite, len(ENEMY_ORDER), len(ENEMY_ORDER))
	for i, name := range ENEMY_ORDER {
		sprite := spriteByName(name)
		ENEMY_ORDER_SPRITES[i] = sprite
	}
}

var monsters []*PlacedEntity

type WillMove struct {
	X, Y int
}

func drawMonsterWillMove(r *sdl.Renderer) {
	for _, monster := range monsters {
		if willMove, ok := monster.Data.(WillMove); monster.Data != nil && ok {
			// full image
			arrowRect := sdl.Rect{X: 0, Y: 0, W: 16, H: 16}

			// default angle is facing down
			var angle float64
			if willMove.X-monster.X > 0 {
				angle = 0
			} else if willMove.Y-monster.Y > 0 {
				angle = 90
			} else if monster.Y-willMove.Y > 0 {
				angle = 270
			} else {
				angle = 180
			}

			// draw half way between the 2 points.
			tgtRect := sdl.Rect{
				X: int32(willMove.X+monster.X) * PIXELS_PER_CELL / 2,
				Y: int32(willMove.Y+monster.Y) * PIXELS_PER_CELL / 2,
				W: 16,
				H: 16,
			}
			err := r.CopyEx(moveArrowTexture, &arrowRect, &tgtRect, angle, nil, 0)
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
		willMove.X = monster.X + willMove.X
		willMove.Y = monster.Y

		// keep it if there is nothing there or the player there.
		if len(otherEntitiesAt(character, willMove.X, willMove.Y)) == 0 {
			monster.Data = willMove
			return
		}

		willMove.Y = character.Y - monster.Y
		if willMove.Y > 1 {
			willMove.Y = 1
		}
		if willMove.Y < -1 {
			willMove.Y = -1
		}
		willMove.Y = monster.Y + willMove.Y
		willMove.X = monster.X

		// keep it if there is nothing there or the player there.
		if len(otherEntitiesAt(character, willMove.X, willMove.Y)) == 0 {
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

type Position struct {
	X, Y int
}

func getEmptyPositions() []Position {
	positions := make([]Position, 0)
	for x := 0; x < MAX_X; x++ {
		for y := 0; y < MAX_Y; y++ {
			occupied := entitiesAt(x, y)
			if len(occupied) == 0 {
				positions = append(positions, Position{x, y})
			}
		}
	}
	return positions
}

func getNewMonsterPosition() (x, y int) {
	allPositions := getEmptyPositions()
	if len(allPositions) == 0 {
		panic("suicide")
	}

	bestPositions := make([]Position, 0)
	for _, pos := range allPositions {
		if (pos.X == character.X+1 || pos.X == character.X-1) && pos.Y == character.Y {
			// skip this position
			continue
		}
		if (pos.Y == character.Y+1 || pos.Y == character.Y-1) && pos.X == character.X {
			// skip this position
			continue
		}
		bestPositions = append(bestPositions, pos)
	}

	// take best position first
	if len(bestPositions) > 0 {
		i := rand.Intn(len(bestPositions))
		return bestPositions[i].X, allPositions[i].Y
	}

	i := rand.Intn(len(allPositions))
	return allPositions[i].X, allPositions[i].Y
}

func spawnMonster() {
	x, y := getNewMonsterPosition()
	newMonster := &PlacedEntity{
		Sprite: ENEMY_ORDER_SPRITES[0],
		X:      x,
		Y:      y,
	}
	monsters = append(monsters, newMonster)
}

func removeMonster(monster *PlacedEntity) {
	for i, item := range monsters {
		if item == monster {
			monsters = append(monsters[:i], monsters[i+1:]...)
		}
	}
}

func upgrade(monster *PlacedEntity) bool {
	for i, upgradeSprite := range ENEMY_ORDER_SPRITES {
		if upgradeSprite == monster.Sprite && i+1 < len(ENEMY_ORDER_SPRITES) {
			monster.Sprite = ENEMY_ORDER_SPRITES[i+1]
			return true
		}
	}
	return false
}

func downgrade(monster *PlacedEntity) bool {
	for i, upgradeSprite := range ENEMY_ORDER_SPRITES {
		if upgradeSprite == monster.Sprite && i-1 >= 0 {
			monster.Sprite = ENEMY_ORDER_SPRITES[i-1]
			return true
		}
	}
	return false
}
