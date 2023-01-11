// Exercises:
//  1. Add collisons: a ball bounces off the screen borders.
//  2. Randomize a speed vector.
//  3. An initial ball position is determined by a mouse position during a
//     click.
//     Use `inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)`
//     function to determine a mouse click.
//     Use `ebiten.CursorPosition()` to get a cursor's position.
//  4. Add friction: a ball slows down.
//  5. Add more balls to the field while clicking a mouse (with the cursor's
//     position). No balls are added, if the new ball's position is overlapping
//     any other ball.
//  6. Add collisisions between balls.
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

	// Default radius of the ball.
	radius = 10
	// Ball speed in pixels/ms.
	speed = 5
)

type Point struct {
	x, y float64
}

type Ball struct {
	pos, vel Point
	color    color.RGBA
}

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

func (b *Ball) Update(dtMs float64, fieldWidth, fieldHeight int) {
	dtMs /= 20
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
	switch {
	case b.pos.x >= float64(fieldWidth):
		b.pos.x = 0
	case b.pos.x < 0:
		b.pos.x = float64(fieldWidth)
	case b.pos.y >= float64(fieldHeight):
		b.pos.y = 0
	case b.pos.y < 0:
		b.pos.y = float64(fieldHeight)
	}
}

func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
}

type Game struct {
	width, height int
	circle        *Ball
	last          time.Time
}

func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		circle: NewBall(width/2, height/2),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t

	g.circle.Update(dt, g.width, g.height)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.circle.Draw(screen)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
