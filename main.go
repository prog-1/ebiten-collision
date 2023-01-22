package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//---------------------------Declaration--------------------------------

const (
	screenWidth  = 640
	screenHeight = 480

	// Ball radius.
	radius = 20
	// Ball default speed in px/ms.
	speed = 0.6
)

type Game struct {
	width, height int
	balls         []*Ball
	// last is a timestamp when Update was called last time.
	last time.Time
}

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

//---------------------------Update-------------------------------------

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.balls = append(g.balls, NewBall(ebiten.CursorPosition()))
	}
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for _, ball := range g.balls {
		ball.Update(dt, g.width, g.height)
	}
	return nil
}

// Update changes a ball state.
// dtMs defines a time interval in microseconds between now and a previous time
// when Update was called.
func (b *Ball) Update(dtMs float64, fieldWidth, fieldHeight int) {
	switch {
	case b.pos.x+radius >= float64(fieldWidth): //right
		b.vel.x = -b.vel.x
	case b.pos.x-radius < 0: //left
		b.vel.x = -b.vel.x
	case b.pos.y+radius >= float64(fieldHeight): //top
		b.vel.y = -b.vel.y
	case b.pos.y-radius < 0: //bottom
		b.vel.y = -b.vel.y
	}
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
}

//---------------------------Draw-------------------------------------

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for _, ball := range g.balls {
		ball.Draw(screen)
	}
}

// Draw renders a ball on a screen.
func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
}

//-------------------------Functions----------------------------------

// NewBall initializes and returns a new Ball instance.
func NewBall(x, y int) *Ball {
	//sp := rand.Float64()
	fmt.Println("new ball")
	rad := rand.Float64() * math.Pi * 2 //random ball velocity
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(rad) * speed, // Cos(math.Pi/4)
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

//---------------------------Main-------------------------------------

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:    width,
		height:   height,
		last:     time.Now(),
	}
}
