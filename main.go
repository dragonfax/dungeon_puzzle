package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Frame struct {
	X, Y, W, H int
	// PixelData
}

type Sprite struct {
	X, Y, W, H int
	FrameCount int
	Frames     []Frame
	Name       string
}

var sprites []Sprite

func read_tiles() {
	sprites = make([]Sprite, 0)

	parseLine := regexp.MustCompile(`^([a-z_]+) (\d+) (\d+) (\d+) (\d+)( (\d+))?$`)

	file, err := os.Open("sprites/tiles_list_v1")
	defer file.Close()

	rd := bufio.NewReader(file)

	for err == nil {
		line, isPrefix, err := rd.ReadLine()
		if isPrefix {
			panic(err)
		}
		if err != nil {
			panic(err)
		}

		matches := parseLine.FindStringSubmatch(string(line))
		if len(matches) == 0 {
			continue
		}
		name := matches[1]
		x, err := strconv.Atoi(matches[2])
		y, err := strconv.Atoi(matches[3])
		w, err := strconv.Atoi(matches[4])
		h, err := strconv.Atoi(matches[5])
		frames := 1
		if len(matches) == 7 {
			frames, err = strconv.Atoi(matches[6])
		}

		sprite := Sprite{
			Name:       name,
			X:          x,
			Y:          y,
			W:          w,
			H:          h,
			FrameCount: frames,
		}

		sprites = append(sprites, sprite)

	}

}

var pixelTex *sdl.Texture

func read_pixels(r *sdl.Renderer) {
	var err error
	pixelTex, err = img.LoadTexture(r, "sprites/0x72_DungeonTilesetII_v1.png")
	if err != nil {
		panic(err)
	}
}

func main() {
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
	read_pixels(r)
	err = r.Copy(pixelTex, nil, nil)
	if err != nil {
		panic(err)
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

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
