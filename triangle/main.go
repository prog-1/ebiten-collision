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
)

type Point struct {
	x, y       float64
	velX, velY float64
}

type Line struct {
	point1 Point
	point2 Point
	color  color.RGBA
}
type Game struct {
	width, height int
	point         []*Point
	line          []*Line
	// last is a timestamp when Update was called last time.
	last time.Time
}

func (l *Line) Update(dtMs float64, a, b *Point) {
	l.point1.x = a.x
	l.point1.y = a.y
	l.point2.x = b.x
	l.point2.y = b.y
}

func (p *Point) Update(dtMs float64, fieldWidth, fieldHeight int) {
	switch {
	case p.x >= float64(fieldWidth):
		p.x = float64(fieldWidth)
		p.velX *= -1
	case p.x <= 0:
		p.x = 0
		p.x *= -1
	case p.y >= float64(fieldHeight):
		p.y = float64(fieldHeight)
		p.velY *= -1
	case p.y <= 0:
		p.y = 0
		p.velY *= -0.99

	}
	p.x += p.velX * dtMs
	p.y += p.velY * dtMs

}

func NewGame(width, height int) *Game {
	var points [3]Point
	for i := range points {
		points[i] = *NewPoint()
	}
	return &Game{
		width:  width,
		height: height,
		// A new ball is created at the cursor position
		point: []*Point{},
		line:  []*Line{},
		last:  time.Now(),
	}
}

func startDir() float64 {
	a := rand.Float64()
	b := rand.Intn(2)
	if b == 0 {
		return -a
	}
	return a
}

func NewLine(a, b *Point) *Line {
	return &Line{
		point1: Point{x: a.x, y: a.y},
		point2: Point{x: b.x, y: b.y},
		color:  color.RGBA{255, 255, 255, 0},
	}
}

func NewPoint() *Point {
	speed := startDir()
	x, y := rand.Int()*int(startDir()), rand.Int()*int(startDir())
	return &Point{
		x:    float64(x),
		y:    float64(y),
		velX: math.Cos(math.Pi/4) * speed,
		velY: math.Cos(math.Pi/4) * speed,
	}
}

func (l *Line) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, l.point1.x, l.point1.y, l.point2.x, l.point2.y, l.color)
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.line[0].Update(dt, g.point[0], g.point[1])
	g.line[1].Update(dt, g.point[1], g.point[2])
	g.line[2].Update(dt, g.point[0], g.point[2])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{45, 45, 45, 45})
	for j := range g.line {
		g.line[j].Draw(screen)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Triangle")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
