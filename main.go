package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const PIXELS_PER_CELL = 16
const CELLS_PER_BOARD = 6
const MAX_X = CELLS_PER_BOARD - 1
const MAX_Y = MAX_X
const TICKS_PER_SPAWN = 30

var character *PlacedEntity
var moveArrowTexture *sdl.Texture

func main() {

	read_tiles()

	characterSprite := spriteByName("necromancer_idle_anim")
	characterHitSprite := spriteByName("necromancer_run_anim")
	character = &PlacedEntity{
		Sprite:    characterSprite,
		HitSprite: characterHitSprite,
	}

	attackTimer := 0

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	r.SetScale(4.0, 4.0)

	read_pixels(r)
	read_reticle(r)
	loadEnemyOrderSprites()

	cursor := sdl.CreateColorCursor(reticleSur, 4, 4)
	sdl.SetCursor(cursor)

	var floor [][]*Sprite = generateFloor(CELLS_PER_BOARD)

	monsters = make([]*PlacedEntity, 0)

	running := true
	tick := 0
	for running {

		monsterThinkTick := tick % TICKS_PER_SPAWN
		if monsterThinkTick == 0 {
			monstersMove()
			spawnMonster()
		}
		if monsterThinkTick == TICKS_PER_SPAWN/2 {
			monstersThink()
		}

		r.Clear()
		showFloor(tick, r, floor)
		drawMonsterWillMove(r)

		drawEntities(tick, r, monsters)

		// draw player
		if attackTimer > 0 {
			attackTimer--
			drawSpriteAt(tick, r, character.HitSprite, int32(character.X), int32(character.Y), 0)
		} else {
			drawSpriteAt(tick, r, character.Sprite, int32(character.X), int32(character.Y), 0)
		}

		r.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			moveToX := character.X
			moveToY := character.Y
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_a:
						moveToX = character.X - 1
						if moveToX < 0 {
							moveToX = 0
						}
					case sdl.K_d:
						moveToX = character.X + 1
						if moveToX > MAX_X {
							moveToX = MAX_X
						}
					case sdl.K_w:
						moveToY = character.Y - 1
						if moveToY < 0 {
							moveToY = 0
						}
					case sdl.K_s:
						moveToY = character.Y + 1
						if moveToY > MAX_Y {
							moveToY = MAX_Y
						}
					case sdl.K_SPACE:
						// attack
						attackTimer = 3
					case sdl.K_LEFT:
						swipe(LEFT)
					case sdl.K_RIGHT:
						swipe(RIGHT)
					case sdl.K_UP:
						swipe(UP)
					case sdl.K_DOWN:
						swipe(DOWN)
					}

					if moveToX != character.X || moveToY != character.Y {
						monstersAt := otherEntitiesAt(character, moveToX, moveToY)
						if len(monstersAt) == 0 {
							character.X = moveToX
							character.Y = moveToY
						} else {
							// attack monster
							monster := monstersAt[0]
							if !downgrade(monster) {
								removeMonster(monster)
							}
						}
					}

					if len(otherEntitiesAt(character, character.X, character.Y)) != 0 {
						panic("suicide")
					}

				}
			case *sdl.QuitEvent:
				running = false
			}
		}

		time.Sleep(time.Second / 30)
		tick = tick + 1
	}
}
