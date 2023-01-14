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
	// speed = 0.4
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

func sign() float64 {
	if rand.Intn(2) == 0 {
		return -1.0
	}
	return 1.0
}

// NewBall initializes and returns a new Ball instance.
func NewBall() *Ball {
	x, y := ebiten.CursorPosition()
	return &Ball{
		pos: Point{float64(x), float64(y)},
		vel: Point{
			x: math.Cos(math.Pi/4) * rand.Float64() * sign(),
			y: math.Sin(math.Pi/4) * rand.Float64() * sign(),
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
		b.vel.x = -b.vel.x
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	case b.pos.x-radius <= 0:
		b.pos.x = 1 + radius
		b.vel.x = -b.vel.x
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	case b.pos.y+radius >= float64(fieldHeight):
		b.pos.y = float64(fieldHeight-1) - radius
		b.vel.y = -b.vel.y
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	case b.pos.y-radius <= 0:
		b.pos.y = 1 + radius
		b.vel.y = -b.vel.y
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	}
}

// Draw renders a ball on a screen.
func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
	tmp := b.color
	for i := len(b.track) - 1; i >= 0; i-- {
		ebitenutil.DrawCircle(screen, b.track[i].x, b.track[i].y, 2, tmp)
		tmp.A -= 10
	}
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
		balls: []*Ball{},
		last:  time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.balls = append(g.balls, NewBall())
	}
	for i := range g.balls {
		g.balls[i].Update(dt, g.width, g.height)
		if len(g.balls[i].track) < 25 {
			g.balls[i].track = append(g.balls[i].track, g.balls[i].pos)
		} else {
			g.balls[i].track = append(g.balls[i].track[1:], g.balls[i].pos)
		}

	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.balls {
		g.balls[i].Draw(screen)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
