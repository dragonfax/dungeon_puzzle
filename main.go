package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

const UNIT_SIZE = 16
const FIELD_WIDTH = 16 * UNIT_SIZE

func drawHorde(tick int, r *sdl.Renderer) {
	sprite := spriteByName("goblin_idle_anim")
	for y := 0; y < 2; y++ {
		for x := 0; x < 4; x++ {
			drawSpriteAt(tick, r, sprite, int32(x+4)*UNIT_SIZE, int32(y+2)*UNIT_SIZE, 0)
		}
	}
}

const Pi = 3.14159

func degrees2Radians(d float64) float64 {
	return d * (Pi / 180)
}

func main() {

	spriteMap := flag.Bool("sprite-map", false, "show the sprites")
	flag.Parse()

	read_tiles()
	weapons := gatherWeapons()
	placedWeapons, weaponSpace := placeWeapons(weapons)

	character := spriteByName("wizzard_m_idle_anim")
	characterHit := spriteByName("wizzart_m_hit_anim")
	attackTimer := 0
	characterShape := resolv.NewRectangle(4*UNIT_SIZE, 4*UNIT_SIZE, character.Frames[0].W, character.Frames[0].H)
	var weilded *PlacedEntity
	var weildedSwinging = false
	var weildedSwingAngle = 0.0

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

	cursor := sdl.CreateColorCursor(reticleSur, 4, 4)
	sdl.SetCursor(cursor)

	var floor [][]*Sprite

	running := true
	tick := 0
	for running {

		r.Clear()
		if *spriteMap {
			showSpriteMap(tick, r)
		} else {
			if floor == nil {
				floor = generateFloor()
			}
			showFloor(tick, r, floor)

			drawHorde(tick, r)

			weaponsColliding := weaponSpace.GetCollidingShapes(characterShape)
			if weaponsColliding.Length() > 0 {
				fmt.Printf("%d weapons colliding\n", weaponsColliding.Length())
				// take the weapon off the field.
				collidingWeapon := weaponsColliding.Get(0)
				weaponSpace.Remove(collidingWeapon)
				weapon, ok := collidingWeapon.GetData().(*PlacedEntity)
				if !ok {
					panic(fmt.Sprintf("wasm't a weapon (*PlacedEntity) was %T ", collidingWeapon.GetData()))
				}
				placedWeapons = removePlacedEntity(placedWeapons, weapon)
				// give the weapon to the player.j
				weilded = weapon
			}

			if weilded != nil {
				// TODO follow player

				if weildedSwinging {
					weildedSwingAngle += 30.0
					if weildedSwingAngle > 300 {
						weildedSwinging = false
						weildedSwingAngle = 0.0
					}
				}

				weildedX := int32(math.Cos(degrees2Radians(weildedSwingAngle)) * UNIT_SIZE)
				weildedY := int32(math.Sin(degrees2Radians(weildedSwingAngle)) * UNIT_SIZE)

				drawSpriteAt(tick, r, weilded.Sprite, weildedX+characterShape.X, weildedY+characterShape.Y, weildedSwingAngle)
			}

			drawEntities(tick, r, placedWeapons)

			// draw player
			if attackTimer > 0 {
				attackTimer--
				drawSpriteAt(tick, r, characterHit, characterShape.X, characterShape.Y, 0)
			} else {
				drawSpriteAt(tick, r, character, characterShape.X, characterShape.Y, 0)
			}

		}
		r.Present()

		/*
			surface, err := window.GetSurface()
			if err != nil {
				panic(err)
			}
			surface.FillRect(nil, 0)

			rect := sdl.Rect{0, 0, 200, 200}
			surface.FillRect(&rect, 0xffff0000)
			window.UpdateSurface()
		*/

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_LEFT, sdl.K_a:
						characterShape.X = characterShape.X - UNIT_SIZE
						if characterShape.X < 0 {
							characterShape.X = 0
						}
					case sdl.K_RIGHT, sdl.K_d:
						characterShape.X = characterShape.X + UNIT_SIZE
						if characterShape.X >= FIELD_WIDTH {
							characterShape.X = FIELD_WIDTH - UNIT_SIZE
						}
					case sdl.K_UP, sdl.K_w:
						characterShape.Y = characterShape.Y - UNIT_SIZE
						if characterShape.Y < 0 {
							characterShape.Y = 0
						}
					case sdl.K_DOWN, sdl.K_s:
						characterShape.Y = characterShape.Y + UNIT_SIZE
						if characterShape.Y >= FIELD_WIDTH {
							characterShape.Y = FIELD_WIDTH - UNIT_SIZE
						}
					case sdl.K_SPACE:
						// attack
						attackTimer = 3
						if weilded != nil {
							weildedSwinging = true
							weildedSwingAngle = 0.0
						}
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
