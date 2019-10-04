package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const UNIT_SIZE = 16
const FIELD_WIDTH = 16

type Sprite struct {
	FrameCount int
	Frames     []sdl.Rect
	Name       string
	Tags       []string
}

var sprites []Sprite

func includesTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func spritesWithTag(tag string) []*Sprite {
	tagged := make([]*Sprite, 0)
	for i := range sprites {
		sprite := &sprites[i]
		if includesTag(sprite.Tags, tag) {
			fmt.Printf("tags %v includes tag %s\n", sprite.Tags, tag)
			tagged = append(tagged, sprite)
		}
	}
	if len(tagged) == 0 {
		panic("no sprite with tag " + tag)
	}
	return tagged
}

func spriteWithTag(tag string) *Sprite {
	return spritesWithTag(tag)[0]
}

func spriteByName(name string) *Sprite {
	for _, sprite := range sprites {
		if sprite.Name == name {
			return &sprite
		}
	}
	panic("no sprite by name " + name)
}

func read_tiles() {
	sprites = make([]Sprite, 0)

	parseLine := regexp.MustCompile(`^([0-9a-z_]+)\s+(\d+) (\d+) (\d+) (\d+)( (\d+))?$`)

	file, err := os.Open("sprites/tiles_list_v1")
	defer file.Close()

	rd := bufio.NewReader(file)

	for err == nil {
		lineBytes, isPrefix, err := rd.ReadLine()
		if isPrefix {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		line := string(lineBytes)
		if line == "" {
			continue
		}

		matches := parseLine.FindStringSubmatch(line)
		if len(matches) == 0 {
			panic("didn't match " + line)
		}
		name := matches[1]
		x, err := strconv.Atoi(matches[2])
		y, err := strconv.Atoi(matches[3])
		w, err := strconv.Atoi(matches[4])
		h, err := strconv.Atoi(matches[5])
		frames := 1
		if matches[7] != "" {
			frames, err = strconv.Atoi(matches[7])
			if err != nil {
				panic(err)
			}
		}

		tags := strings.Split(name, "_")
		// fmt.Printf("sprite has tags %v\n", tags)

		sprite := Sprite{
			Name:       name,
			FrameCount: frames,
			Tags:       tags,
		}

		frameList := make([]sdl.Rect, 0)
		for x1 := 0; x1 < sprite.FrameCount; x1++ {

			frameList = append(frameList,
				sdl.Rect{
					X: int32(x + x1*w),
					Y: int32(y),
					W: int32(w),
					H: int32(h),
				},
			)
		}

		sprite.Frames = frameList

		sprites = append(sprites, sprite)

	}

}

var pixelTex *sdl.Texture
var reticleSur *sdl.Surface

func read_reticle(r *sdl.Renderer) {
	sur, err := img.Load("sprites/reticle.png")
	if err != nil {
		panic(err)
	}
	reticleSur = gfx.ZoomSurface(sur, 8, 8, 0)
}

func read_pixels(r *sdl.Renderer) {
	var err error
	pixelTex, err = img.LoadTexture(r, "sprites/0x72_DungeonTilesetII_v1.png")
	if err != nil {
		panic(err)
	}
}

func drawHorde(tick int, r *sdl.Renderer) {
	sprite := spriteByName("goblin_idle_anim")
	for y := 0; y < 2; y++ {
		for x := 0; x < 4; x++ {
			drawSpriteAt(tick, r, sprite, int32(x+4)*UNIT_SIZE, int32(y+2)*UNIT_SIZE)
		}
	}
}

func showFloor(tick int, r *sdl.Renderer, floor [][]*Sprite) {
	for y := 0; y < len(floor); y++ {
		for x := 0; x < len(floor[y]); x++ {
			sprite := floor[y][x]

			drawSpriteAt(tick, r, sprite, int32(x)*UNIT_SIZE, int32(y)*UNIT_SIZE)
		}
	}
}

func drawEntities(tick int, r *sdl.Renderer, entities []*PlacedEntity) {
	for _, entity := range entities {
		drawSpriteAt(tick, r, entity.Sprite, entity.Shape.X*UNIT_SIZE, entity.Shape.Y*UNIT_SIZE)
	}
}

func drawSpriteAt(tick int, r *sdl.Renderer, sprite *Sprite, x, y int32) {
	animIndex := tick % sprite.FrameCount
	if sprite.FrameCount == 1 {
		animIndex = 0
	}
	frame := sprite.Frames[animIndex]

	tgtRect := sdl.Rect{
		X: x,
		Y: y,
		W: frame.W,
		H: frame.H,
	}

	err := r.Copy(pixelTex, &frame, &tgtRect)
	if err != nil {
		panic(err)
	}
}

func showSpriteMap(tick int, r *sdl.Renderer) {
	x := int32(0)
	y := int32(0)
	for i := 0; i < len(sprites); i++ {

		sprite := &sprites[i]
		drawSpriteAt(tick, r, sprite, x, y)

		x = x + sprite.Frames[0].W

		if x > 200 {
			x = 0
			y = y + UNIT_SIZE
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

func gatherWeapons() []*Sprite {
	return spritesWithTag("weapon")
}

type PlacedEntity struct {
	Sprite *Sprite
	Shape  *resolv.Rectangle
}

func placeWeapons(weapons []*Sprite) ([]*PlacedEntity, *resolv.Space) {
	space := resolv.NewSpace()

	placedWeapons := make([]*PlacedEntity, 0)
	for i := range weapons {
		weapon := weapons[i]
		placedWeapon := &PlacedEntity{
			Sprite: weapon,
			Shape:  resolv.NewRectangle(rand.Int31n(FIELD_WIDTH), rand.Int31n(FIELD_WIDTH), weapon.Frames[0].W, weapon.Frames[0].H),
		}
		placedWeapon.Shape.SetData(placedWeapon)
		space.Add(placedWeapon.Shape)
		placedWeapons = append(placedWeapons, placedWeapon)
	}
	return placedWeapons, space
}

func removePlacedEntity(input []*PlacedEntity, entity *PlacedEntity) []*PlacedEntity {
	for i, e := range input {
		if e == entity {
			// delete the item
			return append(input[:i], input[i+1:]...)
		}
	}
	panic("can't remove entity")
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
	characterShape := resolv.NewRectangle(4, 4, character.Frames[0].W, character.Frames[0].H)
	var weilded *Sprite

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
				weilded = weapon.Sprite
			}

			if weilded != nil {
				drawSpriteAt(tick, r, weilded, (characterShape.X+1)*UNIT_SIZE, (characterShape.Y-1)*UNIT_SIZE)
			}

			drawEntities(tick, r, placedWeapons)

			// draw player
			if attackTimer > 0 {
				attackTimer--
				drawSpriteAt(tick, r, characterHit, characterShape.X*UNIT_SIZE, characterShape.Y*UNIT_SIZE)
			} else {
				drawSpriteAt(tick, r, character, characterShape.X*UNIT_SIZE, characterShape.Y*UNIT_SIZE)
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
						characterShape.X = characterShape.X - 1
						if characterShape.X < 0 {
							characterShape.X = 0
						}
					case sdl.K_RIGHT, sdl.K_d:
						characterShape.X = characterShape.X + 1
						if characterShape.X >= FIELD_WIDTH {
							characterShape.X = FIELD_WIDTH - 1
						}
					case sdl.K_UP, sdl.K_w:
						characterShape.Y = characterShape.Y - 1
						if characterShape.Y < 0 {
							characterShape.Y = 0
						}
					case sdl.K_DOWN, sdl.K_s:
						characterShape.Y = characterShape.Y + 1
						if characterShape.Y >= FIELD_WIDTH {
							characterShape.Y = FIELD_WIDTH - 1
						}
					case sdl.K_SPACE:
						// attack
						attackTimer = 3
					}

				}
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		time.Sleep(time.Second / 30)
		tick = tick + 1
	}
}
