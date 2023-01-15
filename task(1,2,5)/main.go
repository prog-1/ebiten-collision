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
	//speed = 0.4
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

func random() (x float64, y float64, s float64) {
	r1 := rand.Intn(2)
	r2 := rand.Intn(2)
	r3 := rand.Intn(2)
	s = rand.Float64()
	if r3 == 0 {
		s = -s
	}
	if int(r1)%2 == 0 && int(r2)%2 == 0 {
		x, y = -1, -1
		return float64(y), float64(x), float64(s)
	} else if int(r1)%2 == 0 && int(r2)%2 != 0 {
		x, y = 1, -1
		return float64(x), float64(y), float64(s)
	} else if int(r1)%2 != 0 && int(r2)%2 == 0 {
		x, y = -1, 1
		y = 1
		return float64(x), float64(y), float64(s)
	} else {
		x, y = 1, 1
		return float64(x), float64(y), float64(s)
	}
}

// NewBall initializes and returns a new Ball instance.
func NewBall(x, y int) *Ball {
	x1, y1, speed := random()
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(math.Pi/4) * speed * x1,
			y: math.Sin(math.Pi/4) * speed * y1,
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
		b.vel.x = -b.vel.x
		//b.vel.x = -rand.Float64()
	case b.pos.x-radius < 0:
		b.vel.x = -b.vel.x
		//b.vel.x = rand.Float64()
	case b.pos.y+radius >= float64(fieldHeight):
		b.vel.y = -b.vel.y
		//b.vel.y = -rand.Float64()
	case b.pos.y-radius < 0:
		//b.vel.y = rand.Float64()
		b.vel.y = -b.vel.y
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
	ball          *Ball
	// last is a timestamp when Update was called last time.
	last time.Time
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		// A new ball is created at the center of the screen.
		ball: NewBall(width/2, height/2),
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
	g.ball.Update(dt, g.width, g.height)
	for _, v := range g.balls {
		v.Update(dt, g.width, g.height)
	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.ball.Draw(screen)
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
