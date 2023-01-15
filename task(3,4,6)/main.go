package main

import (
	"image/color"
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
	radius = 20
	// Ball default speed in px/ms.
	speed = 0.9
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
	color color.RGBA
}

// NewBall initializes and returns a new Ball instance.
func NewBall(x, y int) *Ball {
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(math.Pi/4) * speed,
			y: math.Sin(math.Pi/4) * speed,
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
		b.pos.x = float64(fieldWidth-1) - radius
		b.vel.x = -(b.vel.x/3 + b.vel.x/2)
		b.vel.y = (b.vel.y/3 + b.vel.y/2)

	case b.pos.x-radius <= 0:
		b.pos.x = 1 + radius
		b.vel.x = -(b.vel.x/3 + b.vel.x/2)
		b.vel.y = (b.vel.y/3 + b.vel.y/2)

	case b.pos.y+radius >= float64(fieldHeight):
		b.pos.y = float64(fieldHeight-1) - radius
		b.vel.x = (b.vel.x/3 + b.vel.x/2)
		b.vel.y = -(b.vel.y/3 + b.vel.y/2)

	case b.pos.y-radius <= 0:
		b.pos.y = 1 + radius
		b.vel.x = (b.vel.x/3 + b.vel.x/2)
		b.vel.y = -(b.vel.y/3 + b.vel.y/2)
	case b.vel.x > 0 && b.vel.x <= 0.1:
		b.vel.x = 0
		b.vel.y = 0
	case b.vel.y > 0 && b.vel.y <= 0.1:
		b.vel.y = 0
		b.vel.x = 0
	}

}

// Draw renders a ball on a screen.
func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
}

// Game is a game instance.
type Game struct {
	width, height int
	balls         []*Ball
	// last is a timestamp when Update was called last time.
	last time.Time
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		// A new ball is created at the center of the screen.
		//ball: NewBall(width/2, height/2),
		last: time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.balls = append(g.balls, NewBall(ebiten.CursorPosition()))
	}
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for i, v := range g.balls {
		v.Update(dt, g.width, g.height)
		for _, z := range g.balls[i+1:] {
			col := v.pos.x+2*radius > z.pos.x &&
				v.pos.y < z.pos.y+2*radius &&
				z.pos.x+2*radius > v.pos.x &&
				z.pos.y < v.pos.y+2*radius
			if col {
				v.vel, z.vel = z.vel, v.vel
			}
		}
	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for _, v := range g.balls {
		v.Draw(screen)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
