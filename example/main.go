package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// Ball radius.
	radius = 20
	// Ball default speed in px/ms.
	friction = 0.999
)

// Point is a struct for representing 2D vectors.
type Point struct {
	x, y float64
}

type Ball struct {
	// Ball position on a screen.
	pos Point
	// Ball speed in px/ms.
	vel   Point
	track []Point
	color color.RGBA
}

func randv() (float64, float64) {
	v1 := rand.Intn(2)
	v2 := rand.Intn(2)
	if int(v1)%2 == 0 && int(v2)%2 == 0 {
		v1 = -1
		v2 = -1
		return float64(v1), float64(v2)
	} else if int(v1)%2 == 0 && int(v2)%2 != 0 {
		v1 = 1
		v2 = -1
		return float64(v1), float64(v2)
	} else if int(v1)%2 != 0 && int(v2)%2 == 0 {
		v1 = -1
		v2 = 1
		return float64(v1), float64(v2)
	} else {
		v1 = 1
		v2 = 1
		return float64(v1), float64(v2)
	}
}

// NewBall initializes and returns a new Ball instance.
func NewBall(x, y int) *Ball {
	a := rand.Float64()
	//random vector
	v1, v2 := randv()
	//
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(math.Pi/4) * a * v1,
			y: math.Sin(math.Pi/4) * a * v2,
		},
		color: color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		},
	}
}

// Update changes a ball state.
//
// dtMs defines a time interval in microseconds between now and a previous time
// when Update was called.
func (b *Ball) Update(dtMs float64, fieldWidth, fieldHeight int) {
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
	switch {
	case b.pos.x+radius >= float64(fieldWidth):
		b.pos.x = float64(fieldWidth) - radius
		b.vel.x = -b.vel.x
	case b.pos.x-radius <= 0:
		b.pos.x = radius
		b.vel.x = -b.vel.x
	case b.pos.y+radius >= float64(fieldHeight):
		b.pos.y = float64(fieldHeight) - radius
		b.vel.y = -b.vel.y
	case b.pos.y-radius <= 0:
		b.pos.y = radius
		b.vel.y = -b.vel.y
	}
	b.vel.x, b.vel.y = b.vel.x*friction, b.vel.y*friction
}

// Draw renders a ball on a screen.
func (b *Ball) Draw(screen *ebiten.Image) {
	track_clr := b.color
	track_clr.A = 56
	for i := len(b.track) - 1; i > 0; i-- {
		ebitenutil.DrawLine(screen, b.track[i].x, b.track[i].y, b.track[i/2].x, b.track[i/2].y, track_clr)
		ebitenutil.DrawCircle(screen, b.track[i].x, b.track[i].y, radius, track_clr)
		track_clr.A -= 5
	}
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
}

// Game is a game instance.
type Game struct {
	width, height int
	ball          []*Ball
	// last is a timestamp when Update was called last time.
	last time.Time
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		last:   time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
	t := time.Now()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.ball = append(g.ball, NewBall(ebiten.CursorPosition()))
	}
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for i, b := range g.ball {
		b.Update(dt, g.width, g.height)
		for _, ba := range g.ball[i+1:] {
			Collided := b.pos.x+2*radius > ba.pos.x &&
				b.pos.y < ba.pos.y+2*radius &&
				ba.pos.x+2*radius > b.pos.x &&
				ba.pos.y < b.pos.y+2*radius
			if Collided {
				b.vel, ba.vel = ba.vel, b.vel
			}
		}
		if len(b.track) < 10 {
			b.track = append(b.track, b.pos)
			continue
		}
		b.track = b.track[1:]

	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.ball {
		b.Draw(screen)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
