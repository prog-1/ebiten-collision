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
	maxSpeed = 1
	friction = 0.99
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
	rad := rand.Float64() * 2 * math.Pi
	speed := rand.Float64() * maxSpeed
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(rad) * speed,
			y: math.Sin(rad) * speed,
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
// dtMs defines a time interval in milliseconds between now and a previous time
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
	default:
		return
	}
	b.vel.x, b.vel.y = b.vel.x*friction, b.vel.y*friction
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
		last: time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return screenWidth, screenHeight
}

// Update updates a game state.
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.balls = append(g.balls, NewBall(ebiten.CursorPosition()))
	}
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for i, b := range g.balls {
		// Elastic collision - balls swaps their velocities
		// Algo:
		// Check whether distance between one ball and other balls <= radius
		//	1.
		// Swap velocities
		for _, b1 := range g.balls[i+1:] {
			dif := Point{b.pos.x - b1.pos.x, b.pos.y - b1.pos.y}
			dist := math.Sqrt(math.Pow(dif.x, 2) + math.Pow(dif.y, 2))
			if math.Abs(dist) <= 2*radius {
				b.vel, b1.vel = b1.vel, b.vel
			}
		}
		b.Update(dt, g.width, g.height)
	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.balls {
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
