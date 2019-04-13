package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

var (
	conf struct {
		res  string
		pref string
		mono bool
	}

	si       *SI
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	timer    uint32
	saveslot int
	paused   bool
)

func main() {
	runtime.LockOSThread()
	si = NewSI()
	parseFlags()
	initSDL()
	loadRes()
	reset()
	for {
		event()
		update()
		blit()
	}
}

func parseFlags() {
	conf.res = sdl.GetBasePath()
	conf.pref = sdl.GetPrefPath("", "space_invaders")
	flag.StringVar(&conf.res, "res", conf.res, "resource directory")
	flag.StringVar(&conf.pref, "pref", conf.pref, "preference directory")
	flag.BoolVar(&conf.mono, "mono", conf.mono, "monochrome color")
	flag.Parse()
}

func loadRes() {
	err := si.LoadROM(filepath.Join(conf.res, "rom"))
	ck(err)

	err = si.LoadSound(filepath.Join(conf.res, "snd"))
	ck(err)
}

func initSDL() {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER)
	ck(err)

	err = sdl.InitSubSystem(sdl.INIT_AUDIO)
	ek(err)

	err = sdlmixer.OpenAudio(sdlmixer.DEFAULT_FREQUENCY, sdlmixer.DEFAULT_FORMAT, 2, 4096)
	ek(err)

	width, height := 256, 224
	window, err = sdl.CreateWindow("Space Invaders", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width*2, height*2, sdl.WINDOW_RESIZABLE)
	ck(err)

	window.SetMinimumSize(width, height)
	sdl.ShowCursor(0)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	ck(err)

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	ck(err)

	renderer.SetLogicalSize(width, height)
}

func reset() {
	timer = sdl.GetTicks()
	si.Reset()
	si.colored = !conf.mono
}

func event() {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				os.Exit(0)
			case sdl.K_1:
				saveslot++
				fmt.Printf("Saveslot %d\n", saveslot)
			case sdl.K_2:
				if saveslot > 0 {
					saveslot--
				}
				fmt.Printf("Saveslot %d\n", saveslot)
			case sdl.K_F2:
				si.SaveState(filepath.Join(conf.pref, fmt.Sprintf("save%d", saveslot)))
				fmt.Println("Save State")
			case sdl.K_F4:
				si.LoadState(filepath.Join(conf.pref, fmt.Sprintf("save%d", saveslot)))
				fmt.Println("Load State")
			case sdl.K_r:
				reset()
			case sdl.K_p:
				paused = !paused
				fmt.Println("Paused:", paused)
			case sdl.K_c:
				// coin
				si.port1 |= 1 << 0
			case sdl.K_BACKSPACE:
				// p2 start
				si.port1 |= 1 << 1
			case sdl.K_RETURN:
				// p1 start
				si.port1 |= 1 << 2
			case sdl.K_SPACE:
				// p1 shoot
				si.port1 |= 1 << 4
				// p2 shoot
				si.port2 |= 1 << 4
			case sdl.K_LEFT:
				// p1 left
				si.port1 |= 1 << 5
				// p2 left
				si.port2 |= 1 << 5
			case sdl.K_RIGHT:
				// p1 right
				si.port1 |= 1 << 6
				// p2 right
				si.port2 |= 1 << 6
			case sdl.K_t:
				// tilt
				si.port2 |= 1 << 2
			case sdl.K_3:
				si.colored = !si.colored
			}
		case sdl.KeyUpEvent:
			switch ev.Sym {
			case sdl.K_c:
				// coin
				si.port1 &^= 1 << 0
			case sdl.K_BACKSPACE:
				// p2 start
				si.port1 &^= 1 << 1
			case sdl.K_RETURN:
				// p1 start
				si.port1 &^= 1 << 2
			case sdl.K_SPACE:
				// p1 shoot
				si.port1 &^= 1 << 4
				// p2 shoot
				si.port2 &^= 1 << 4
			case sdl.K_LEFT:
				// p1 left
				si.port1 &^= 1 << 5
				// p1 right
				si.port2 &^= 1 << 5
			case sdl.K_RIGHT:
				// p1 right
				si.port1 &^= 1 << 6
				// p2 right
				si.port2 &^= 1 << 6
			case sdl.K_t:
				// tilt
				si.port2 &^= 1 << 2
			}
		}
	}
}

func update() {
	if float64(sdl.GetTicks()-timer) > (1/float64(si.fps))*1000 && !paused {
		si.Update()
		si.render()
		texture.Update(nil, si.fb.Pix, 4*256)
	}
}

func blit() {
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ek(err error) bool {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return true
	}
	return false
}

type SI struct {
	*CPU
	rom      [4][0x800]byte
	fb       *image.RGBA
	snd      [16]*sdlmixer.Chunk
	nextif   uint16
	port1    uint8
	port2    uint8
	shift0   uint8
	shift1   uint8
	shiftoff uint8
	loport3  uint8
	loport5  uint8
	colored  bool
	fps      uint64
	cpf      uint64
}

func NewSI() *SI {
	si := &SI{
		CPU: NewCPU(),
		fb:  image.NewRGBA(image.Rect(0, 0, 256, 224)),
	}
	si.IO = si
	return si
}

func (si *SI) Reset() {
	copy(si.M[0x0000:], si.rom[0][:])
	copy(si.M[0x0800:], si.rom[1][:])
	copy(si.M[0x1000:], si.rom[2][:])
	copy(si.M[0x1800:], si.rom[3][:])
	si.fps = 60
	si.cpf = 2000000 / si.fps

	si.nextif = 0x08
	si.port1 = 1 << 3
	si.port2 = 0
	si.shift0 = 0
	si.shift1 = 0
	si.loport3 = 0
	si.loport5 = 0
	si.CPU.Reset()
}

func (si *SI) Update() {
	hcpf := si.cpf / 2

	var cyc uint64
	for cyc <= si.cpf {
		scyc := si.Cycles
		si.Step()
		cyc += si.Cycles - scyc
		if si.Cycles >= hcpf {
			si.Interrupt(si.nextif)
			si.Cycles -= hcpf
			if si.nextif == 0x8 {
				si.nextif = 0x10
			} else {
				si.nextif = 0x8
			}
		}
	}
}

func (si *SI) render() {
	// the screen is 256 * 224 pixels, and is rotated anti-clockwise.
	// these are the overlay dimensions:
	// ,_______________________________.
	// |WHITE            ^             |
	// |                32             |
	// |                 v             |
	// |-------------------------------|
	// |RED              ^             |
	// |                32             |
	// |                 v             |
	// |-------------------------------|
	// |WHITE                          |
	// |         < 224 >               |
	// |                               |
	// |                 ^             |
	// |                120            |
	// |                 v             |
	// |                               |
	// |                               |
	// |                               |
	// |-------------------------------|
	// |GREEN                          |
	// | ^                  ^          |
	// |56        ^        56          |
	// | v       72         v          |
	// |____      v      ______________|
	// |  ^  |          | ^            |
	// |<16> |  < 118 > |16   < 122 >  |
	// |  v  |          | v            |
	// |WHITE|          |         WHITE|
	// `-------------------------------'

	// the screen is 256 * 224 pixels, and 1 byte contains 8 pixels

	const VRAM_ADDR = 0x2400
	for i := uint(0); i < 256*224/8; i++ {
		y := i * 8 / 256
		base_x := (i * 8) % 256
		cur_byte := si.M[VRAM_ADDR+i]

		for bit := uint(0); bit < 8; bit++ {
			px := base_x + bit
			py := y
			is_pixel_lit := (cur_byte >> bit) & 1

			var r, g, b uint8
			if !si.colored && is_pixel_lit != 0 {
				r, g, b = 255, 255, 255
			} else if si.colored && is_pixel_lit != 0 {
				if px < 16 {
					if py < 16 || py > 118+16 {
						r, g, b = 255, 255, 255
					} else {
						g = 255
					}
				} else if px >= 16 && px <= 16+56 {
					g = 255
				} else if px >= 16+56+120 && px < 16+56+120+32 {
					r = 255
				} else {
					r, g, b = 255, 255, 255
				}
			}

			temp_x := px
			px = py
			py = -temp_x + 256 - 1
			si.fb.SetRGBA(int(px), int(py), color.RGBA{r, g, b, 255})
		}
	}
}

func (si *SI) In(port uint8) {
	var v uint8
	switch port {
	case 1:
		v = si.port1
	case 2:
		v = si.port2
	case 3:
		s := uint16(si.shift1)<<8 | uint16(si.shift0)
		v = uint8((s >> (8 - si.shiftoff)))
	default:
		panic(fmt.Errorf("unknown in port %d", port))
	}
	si.R[A] = v
}

func (si *SI) Out(port uint8) {
	v := si.R[A]
	switch port {
	case 2:
		si.shiftoff = v & 0x7
	case 3:
		si.playSound(1)
	case 4:
		si.shift0 = si.shift1
		si.shift1 = v
	case 5:
		si.playSound(2)
	case 6:
		// debug port?
	default:
		panic(fmt.Errorf("unknown out port %d", port))
	}
}

func (si *SI) LoadROM(dir string) error {
	files := []string{
		"invaders.h",
		"invaders.g",
		"invaders.f",
		"invaders.e",
	}
	for i, file := range files {
		name := filepath.Join(dir, file)
		buf, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		if len(buf) != len(si.rom[i]) {
			return fmt.Errorf("%s: invalid rom file", name)
		}
		copy(si.rom[i][:], buf)
	}
	return nil
}

func (si *SI) LoadSound(dir string) error {
	for i := 0; i <= 10; i++ {
		name := filepath.Join(dir, fmt.Sprintf("%d.wav", i))
		snd, err := sdlmixer.LoadWAV(name)
		if ek(err) {
			continue
		}
		si.snd[i] = snd
	}
	si.snd[0] = si.snd[8]
	si.snd[8] = si.snd[10]
	return nil
}

func (si *SI) playSound(bank int) {
	// play a sound if corresponding bit changed from 0 to 1
	data := si.R[A]
	switch bank {
	case 1:
		if data != si.loport3 {
			for i := uint(0); i < 4; i++ {
				if data&(1<<i) != 0 && si.loport3&(1<<i) == 0 {
					si.snd[i].PlayChannel(-1, 0)
				}
			}
			si.loport3 = data
		}
	case 2:
		if data != si.loport5 {
			for i := uint(0); i < 4; i++ {
				if data&(1<<i) != 0 && si.loport5&(1<<i) != 0 {
					si.snd[4+i].PlayChannel(-1, 0)
				}
			}
			si.loport5 = data
		}
	}
}

func (si *SI) SaveState(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	binary.Write(w, binary.LittleEndian, si.R)
	binary.Write(w, binary.LittleEndian, si.M)
	binary.Write(w, binary.LittleEndian, si.PC)
	binary.Write(w, binary.LittleEndian, si.SP)
	binary.Write(w, binary.LittleEndian, si.CY)
	binary.Write(w, binary.LittleEndian, si.HC)
	binary.Write(w, binary.LittleEndian, si.IF)
	binary.Write(w, binary.LittleEndian, si.HLT)
	binary.Write(w, binary.LittleEndian, si.Cycles)

	w.Flush()
	f.Close()

	return nil
}

func (si *SI) LoadState(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	binary.Read(r, binary.LittleEndian, &si.R)
	binary.Read(r, binary.LittleEndian, &si.M)
	binary.Read(r, binary.LittleEndian, &si.PC)
	binary.Read(r, binary.LittleEndian, &si.SP)
	binary.Read(r, binary.LittleEndian, &si.CY)
	binary.Read(r, binary.LittleEndian, &si.HC)
	binary.Read(r, binary.LittleEndian, &si.Z)
	binary.Read(r, binary.LittleEndian, &si.S)
	binary.Read(r, binary.LittleEndian, &si.P)
	binary.Read(r, binary.LittleEndian, &si.IF)
	binary.Read(r, binary.LittleEndian, &si.HLT)
	binary.Read(r, binary.LittleEndian, &si.Cycles)

	return nil
}
