package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

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
