package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// Ball radius.
	radius = 30
	// Ball default speed in px/ms.
	speed = 0.4
)

// Point is a struct for representing 2D vectors.
type Point struct {
	x, y float64
}

type Ball struct {
	curentSpeed float64
	// Ball position on a screen.
	pos Point
	// Ball speed in px/ms.
	vel   Point
	color color.RGBA
}

// NewBall initializes and returns a new Ball instance.
func NewBall(x, y int) *Ball {
	degrees := float64(rand.Intn(360) + 1)
	return &Ball{
		curentSpeed: speed,
		pos:         Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(degrees * (math.Pi / 180)),
			y: math.Sin(degrees * (math.Pi / 180)),
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
	// if b.
	if b.curentSpeed > 0 {
		b.pos.x += b.vel.x * dtMs * b.curentSpeed
		b.pos.y += b.vel.y * dtMs * b.curentSpeed
		fmt.Println(b.curentSpeed)
		b.curentSpeed -= 0.01 * dtMs / 1000
	}
	switch {
	case b.pos.x+radius >= float64(fieldWidth) && b.vel.x > 0:
		b.vel.x *= -1
		//b.vel.y = math.Sin(-math.Pi/4) * speed
	case b.pos.x-radius < 0 && b.vel.x < 0:
		b.vel.x *= -1

	case b.pos.y+radius >= float64(fieldHeight) && b.vel.y > 0:
		b.vel.y *= -1

	case b.pos.y-radius < 0 && b.vel.y < 0:
		b.vel.y *= -1

	}
}

// Draw renders a ball on a screen.
func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
	// screen.Fill(color.RGBA64{0, 0, 0, 255})
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
		// A new ball is created at the center of the screen.
		ball: []*Ball{},
		last: time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.ball = append(g.ball, NewBall(x, y))
	}
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for i := range g.ball {
		g.ball[i].Update(dt, g.width, g.height)
	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.ball {
		g.ball[i].Draw(screen)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
