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
	radius       = 20
)

type Game struct {
	width, height int
	balls         []*Ball
	last          time.Time
}

type Point struct {
	x, y float64
}

type Ball struct {
	pos Point
	// Ball speed in px/ms.
	vel   Point
	color color.RGBA
	trail *Ball
}

func NewBall(x, y int) *Ball {
	mul := float64(rand.Intn(5) + 1)
	speed := rand.Float64()
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(math.Pi*mul/3) * speed,
			y: math.Sin(math.Pi*mul/3) * speed,
		},
		color: color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		},
		trail: &Ball{pos: Point{x: float64(x + radius - 1), y: float64(y)},
			color: color.RGBA{
				R: uint8(0),
				G: uint8(0),
				B: uint8(0),
				A: 255,
			},
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
	b.trail.pos.x += b.vel.x * dtMs
	b.trail.pos.y += b.vel.y * dtMs
	if b.vel.x > 0 {
		b.vel.x -= 0.0005
	} else if b.vel.x < 0 {
		b.vel.x += 0.0005
	}
	if b.vel.y > 0 {
		b.vel.y -= 0.0005
	} else if b.vel.y < 0 {
		b.vel.y += 0.0005
	}
	switch {
	case b.pos.x+radius >= float64(fieldWidth) && b.pos.y+radius >= float64(fieldHeight):
		b.vel.x = -b.vel.x
		b.vel.y = -b.vel.y
	case b.pos.x+radius >= float64(fieldWidth) && b.pos.y+radius <= 0:
		b.vel.x = -b.vel.x
		b.vel.y = -b.vel.y
	case b.pos.x+radius <= 0 && b.pos.y+radius >= float64(fieldHeight):
		b.vel.x = -b.vel.x
		b.vel.y = -b.vel.y
	case b.pos.x+radius <= 0 && b.pos.y+radius <= 0:
		b.vel.x = -b.vel.x
		b.vel.y = -b.vel.y
	case b.pos.x+radius-1 >= float64(fieldWidth):
		b.vel.x = -b.vel.x
	case b.pos.x-radius+1 <= 0:
		b.vel.x = -b.vel.x
	case b.pos.y+radius-1 >= float64(fieldHeight):
		b.vel.y = -b.vel.y
	case b.pos.y-radius+1 <= 0:
		b.vel.y = -b.vel.y
	}
	b.trail.vel.x = b.vel.x
	b.trail.vel.y = b.vel.y
}

func (b *Ball) DrawTrail(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius/5, b.color)
}

func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
}

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

func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.balls = append(g.balls, NewBall(ebiten.CursorPosition()))
	}
	for _, b := range g.balls {
		b.Update(dt, g.width, g.height)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.balls {
		b.Draw(screen)
		b.DrawTrail(screen)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
