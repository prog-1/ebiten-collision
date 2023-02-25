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
)

// Point is a struct for representing 2D vectors.
type Point struct {
	x, y float64
}

type Trail struct {
	pos   Point
	color color.RGBA
}

type Ball struct {
	// Ball position on a screen.
	pos Point
	// Ball speed in px/ms.
	vel   Point
	color color.RGBA
	trail []Trail
}

func startDir() float64 {
	a := rand.Float64()
	b := rand.Intn(2)
	if b == 0 {
		return -a
	}
	return a
}

// NewBall initializes and returns a new Ball instance.
func NewBall() *Ball {
	speed := startDir()
	x, y := ebiten.CursorPosition()
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
	case b.pos.x >= float64(fieldWidth)-radius:
		b.pos.x = float64(fieldWidth) - radius
		b.vel.x *= -0.99
		b.vel.y *= 0.99
	case b.pos.x <= radius:
		b.pos.x = radius
		b.vel.x *= -0.99
		b.vel.y *= 0.99
	case b.pos.y >= float64(fieldHeight)-radius:
		b.pos.y = float64(fieldHeight) - radius
		b.vel.x *= 0.99
		b.vel.y *= -0.99
	case b.pos.y <= radius:
		b.pos.y = radius
		b.vel.x *= 0.99
		b.vel.y *= -0.99

	}
	b.trail = append(b.trail, Trail{pos: b.pos, color: b.color})
	if len(b.trail) > 45 {
		b.trail = b.trail[1:]
	}
}

// Draw renders a ball on a screen.
func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.pos.x, b.pos.y, radius, b.color)
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
		// A new ball is created at the cursor position
		ball: []*Ball{},
		last: time.Now(),
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
		g.ball = append(g.ball, NewBall())
	}
	for i := 0; i < len(g.ball); i++ {
		for j := i + 1; j < len(g.ball); j++ {
			if g.ball[i].colliding(g.ball[j]) {
				g.ball[i].resolving(g.ball[j])
			}
		}
		g.ball[i].Update(dt, g.width, g.height)
	}
	return nil
}

// I understood how methods works yay!
func (a *Ball) colliding(b *Ball) bool {
	distX, distY := a.pos.x-b.pos.x, a.pos.y-b.pos.y
	distSqr := math.Sqrt(distX*distX + distY*distY)
	return distSqr <= float64(radius*2)
}

func (a *Ball) resolving(b *Ball) { // in fact not really, this method is weird junk
	// also i have no idea if it works with fast-moving balls - my laptop is too slow
	// but with slow-moving it sometimes even work right, but mostly weird

	// distX, distY := b.pos.x-a.pos.x, b.pos.y-a.pos.y
	// distSqr := (distX * distX) + (distY * distY)
	// angle := math.Atan2(b.pos.y-a.pos.y, b.pos.x-a.pos.x)
	// distToMove := 2*radius - distSqr
	// b.pos.x += float64(math.Cos(angle) * distToMove)
	// b.pos.y += float64(math.Cos(angle) * distToMove)

	// i read it should be tangent vectors related code here, but it is too complicated - i dont wanna to work with this rn
	//ehhhh vectors in go... what?
	//https://www.netguru.com/blog/vector-operations-in-go
	a.vel, b.vel = b.vel, a.vel
	a.pos.x, a.pos.y = a.pos.x+a.vel.x, a.pos.y+a.vel.y
	//Sources:
	//https://stackoverflow.com/questions/345838/ball-to-ball-collision-detection-and-handling
	//https://flatredball.com/documentation/tutorials/math/circle-collision/
}

// i hate this method so much!

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{45, 45, 45, 45})
	for _, i := range g.ball {
		i.Draw(screen)
		for j := len(i.trail) - 1; j >= 0; j-- {
			i.trail[j].color.A -= 100
			ebitenutil.DrawCircle(screen, i.trail[j].pos.x, i.trail[j].pos.y, radius, i.trail[j].color)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ball masquarade!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
