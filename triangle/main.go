package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// Ball radius.
	//radius = 20
	// Ball default speed in px/ms.
	speed = 0.4
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
	track []Point
}

// NewBall initializes and returns a new Ball instance.
func NewBall(x, y int) *Ball {
	r := 2 * math.Pi * rand.Float64() // from 0 to 2Pi (rand.Float64() is from 0 to 1)
	return &Ball{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(r) * speed,
			y: math.Sin(r) * speed,
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
	switch {
	case b.pos.x+b.vel.x*dtMs >= float64(fieldWidth-1) || b.pos.x+b.vel.x*dtMs <= 0:
		b.vel.x = -b.vel.x
	case b.pos.y+b.vel.y*dtMs >= float64(fieldHeight-1) || b.pos.y+b.vel.y*dtMs <= 0:
		b.vel.y = -b.vel.y
	}
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
	if len(b.track) > 10 {
		b.track = b.track[1:]
	}
	b.track = append(b.track, b.pos)
	// if b.vel.x > 0 {
	// 	b.vel.x -= 0.001
	// }
	// if b.vel.x < 0 {
	// 	b.vel.x += 0.001
	// }
	// if b.vel.y > 0 {
	// 	b.vel.y -= 0.001
	// }
	// if b.vel.y < 0 {
	// 	b.vel.y += 0.001
	// }
}

var transparency = []byte{0x00, 0x1a, 0x33, 0x4d, 0x66, 0x80, 0x99, 0xb3, 0xcc, 0xe6, 0xff}

// Draw renders a ball on a screen.
// func (b *Ball) DrawCircle(screen *ebiten.Image) {
// 	for i, t := range b.track {
// 		b.color.A = transparency[i]
// 		ebitenutil.DrawCircle(screen, t.x, t.y, radius, b.color)
// 	}
// }

func (g *Game) DrawTriangle(screen *ebiten.Image) {
	color := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	for i := 0; i < len(g.balls[0].track); i++ {
		color.A = transparency[i]
		ebitenutil.DrawLine(screen, g.balls[0].track[i].x, g.balls[0].track[i].y, g.balls[1].track[i].x, g.balls[1].track[i].y, color)
		ebitenutil.DrawLine(screen, g.balls[1].track[i].x, g.balls[1].track[i].y, g.balls[2].track[i].x, g.balls[2].track[i].y, color)
		ebitenutil.DrawLine(screen, g.balls[2].track[i].x, g.balls[2].track[i].y, g.balls[0].track[i].x, g.balls[0].track[i].y, color)
	}
}

// Game is a game instance.
type Game struct {
	width, height int
	balls         [3]Ball
	// last is a timestamp when Update was called last time.
	last time.Time
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	var balls [3]Ball
	for i := range balls {
		balls[i] = *NewBall(rand.Intn(width), rand.Intn(height))
	}
	return &Game{
		width:  width,
		height: height,
		balls:  balls,
		last:   time.Now(),
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
	for i := range g.balls {
		//for _, anotherball := range g.balls[i+1:] {
		//
		//}
		// if len(g.balls[i].track) > 10 {
		// 	g.balls[i].track = g.balls[i].track[1:]
		// }
		// g.balls[i].track = append(g.balls[i].track, g.balls[i].pos)
		g.balls[i].Update(dt, g.width, g.height)
	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// for _, ball := range g.balls {
	// 	ball.DrawCircle(screen)
	// }
	g.DrawTriangle(screen)
}

func main() {
	t := time.Now().UnixNano()
	rand.Seed(t)
	// fmt.Println(t)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
