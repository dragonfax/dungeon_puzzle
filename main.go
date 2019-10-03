package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"

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
}

var sprites []Sprite

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

		sprite := Sprite{
			Name: name,
			Rect: sdl.Rect{
				X: int32(x),
				Y: int32(y),
				W: int32(w),
				H: int32(h),
			},
			FrameCount: frames,
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

func read_pixels(r *sdl.Renderer) {
	var err error
	pixelTex, err = img.LoadTexture(r, "sprites/0x72_DungeonTilesetII_v1.png")
	if err != nil {
		panic(err)
	}
}

func main() {
	read_tiles()

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

	running := true
	tick := 0
	for running {

		r.Clear()

		x := int32(0)
		y := int32(0)
		for i := 0; i < len(sprites); i++ {

			sprite := sprites[i]

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

			x = x + frame.Rect.W

			if x > 200 {
				x = 0
				y = y + 16
			}

			err = r.Copy(pixelTex, &frame.Rect, &tgtRect)
			if err != nil {
				panic(err)
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
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		time.Sleep(time.Second / 30)
		tick = tick + 1
	}
}
