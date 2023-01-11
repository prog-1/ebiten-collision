# Exercises

1. Add collisions between a ball and walls. The ball must bounce off the screen
   border.
2. Make an initial speed vector random.
3. Add a ball only when a mouse button is clicked. An initial position is
   determined by a mouse position at a click. You may use the following
   functions:

    * `[inpututil.IsMouseButtonJustPressed]`(https://pkg.go.dev/github.com/hajimehoshi/ebiten/inpututil#IsMouseButtonJustPressed).
    * `[ebiten.CursorPosition()]`(https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#CursorPosition).
4. Add friction. A ball slows down until it completely stops.
5. Add more balls to the field with a mouse click at a cursor position.
6. Add [elastic collisions](https://en.wikipedia.org/wiki/Elastic_collision)
   between balls.
7. Add a ball track of a fixed length.
8. Add a fading ball track.