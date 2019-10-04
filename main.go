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

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Frame struct {
	sdl.Rect
	// PixelData
}

type Sprite struct {
	FrameCount int
	Frames     []Frame
	Name       string
	Rect       sdl.Rect
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
			Name: name,
			Rect: sdl.Rect{
				X: int32(x),
				Y: int32(y),
				W: int32(w),
				H: int32(h),
			},
			FrameCount: frames,
			Tags:       tags,
		}

		frameList := make([]Frame, 0)
		for x := 0; x < sprite.FrameCount; x++ {

			frameList = append(frameList, Frame{
				sdl.Rect{
					X: sprite.Rect.X + int32(x)*sprite.Rect.W,
					Y: sprite.Rect.Y,
					W: sprite.Rect.W,
					H: sprite.Rect.H,
				},
			})
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

func showFloor(tick int, r *sdl.Renderer, floor [][]*Sprite) {
	for y := 0; y < len(floor); y++ {
		for x := 0; x < len(floor[y]); x++ {
			sprite := floor[y][x]

			drawSpriteAt(tick, r, sprite, int32(x)*16, int32(y)*16)
		}
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
		W: frame.Rect.W,
		H: frame.Rect.H,
	}

	err := r.Copy(pixelTex, &frame.Rect, &tgtRect)
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

		x = x + sprite.Rect.W

		if x > 200 {
			x = 0
			y = y + 16
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

func main() {

	spriteMap := flag.Bool("sprite-map", false, "show the sprites")
	flag.Parse()

	read_tiles()

	character := spriteByName("wizzard_m_idle_anim")
	characterHit := spriteByName("wizzart_m_hit_anim")
	var charX int32 = 4
	var charY int32 = 4
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
		}

		// draw player
		if attackTimer > 0 {
			attackTimer--
			drawSpriteAt(tick, r, characterHit, charX*16, charY*16)
		} else {
			drawSpriteAt(tick, r, character, charX*16, charY*16)
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
					if e.Keysym.Sym == sdl.K_LEFT {
						charX = charX - 1
						if charX < 0 {
							charX = 0
						}
					}
					if e.Keysym.Sym == sdl.K_RIGHT {
						charX = charX + 1
						if charX > 15 {
							charX = 15
						}
					}
					if e.Keysym.Sym == sdl.K_UP {
						charY = charY - 1
						if charY < 0 {
							charY = 0
						}
					}
					if e.Keysym.Sym == sdl.K_DOWN {
						charY = charY + 1
						if charY > 15 {
							charY = 15
						}
					}
					if e.Keysym.Sym == sdl.K_SPACE {
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
