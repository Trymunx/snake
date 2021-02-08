package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	gameloop "github.com/kutase/go-gameloop"
)

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type point struct {
	x int
	y int
}

type snake struct {
	head *point
	tail []*point
}

type gameData struct {
	snake    *snake
	dir      *point
	food     *point
	gameOver bool
}

func newSnake(x, y int) *snake {
	head := &point{x: x, y: y}
	tail := []*point{}
	for i := 1; i < 6; i++ {
		tail = append(tail, &point{x: x - (2 * i), y: y})
	}
	return &snake{head, tail}
}

func newGameData(w, h int) *gameData {
	return &gameData{
		snake: newSnake(w/2, h/2),
		dir:   &point{x: 2, y: 0},
		food:  randomPoint(w, h),
	}
}

func randomPoint(w, h int) *point {
	return &point{
		// Make sure x coords are always a multiple of 2
		x: 2 * (rand.Intn((w-2)/2) + 1),
		y: 1 + rand.Intn(h-2),
	}
}

func main() {
	encoding.Register()
	s, err := tcell.NewScreen()
	handleErr(err)
	if err = s.Init(); err != nil {
		handleErr(err)
	}

	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()

	gd := newGameData(s.Size())
	gl := gameloop.New(10, nil)
	gl.SetOnUpdate(func(delta float64) {
		var last *point
		for i := len(gd.snake.tail) - 1; i >= 0; i-- {
			if last == nil {
				last = &point{x: gd.snake.tail[i].x, y: gd.snake.tail[i].y}
			}
			if i > 0 {
				gd.snake.tail[i].x = gd.snake.tail[i-1].x
				gd.snake.tail[i].y = gd.snake.tail[i-1].y
			} else {
				gd.snake.tail[i].x = gd.snake.head.x
				gd.snake.tail[i].y = gd.snake.head.y
			}
		}
		if gd.snake.head.x == gd.food.x && gd.snake.head.y == gd.food.y {
			gd.snake.tail = append(gd.snake.tail, last)
			gd.food = randomPoint(s.Size())
		}
		gd.snake.head.x += gd.dir.x
		gd.snake.head.y += gd.dir.y
		w, h := s.Size()
		if gd.snake.head.x < 1 || gd.snake.head.x >= w-1 || gd.snake.head.y < 1 || gd.snake.head.y >= h-1 {
			gd.gameOver = true
			gd.draw(s)
			gl.Stop()
		}
		for i := 0; i < len(gd.snake.tail); i++ {
			if gd.snake.head.x == gd.snake.tail[i].x && gd.snake.head.y == gd.snake.tail[i].y {
				gd.gameOver = true
				gd.draw(s)
				gl.Stop()
			}
		}
		gd.draw(s)
	})
	gl.Start()

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			gd.draw(s)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			case tcell.KeyUp:
				gd.dir = &point{x: 0, y: -1}
			case tcell.KeyDown:
				gd.dir = &point{x: 0, y: 1}
			case tcell.KeyLeft:
				gd.dir = &point{x: -2, y: 0}
			case tcell.KeyRight:
				gd.dir = &point{x: 2, y: 0}
			}
			if (ev.Rune() == 'r' || ev.Rune() == 'R') && gd.gameOver {
				w, h := s.Size()
				gd.snake = newSnake(w/2, h/2)
				gd.food = randomPoint(w, h)
				gd.dir = &point{x: 2, y: 0}
				gd.gameOver = false
				gl.Start()
			}
			gd.draw(s)
		}
	}
}

func (gd *gameData) draw(s tcell.Screen) {
	s.Clear()
	s.SetContent(gd.snake.head.x, gd.snake.head.y, '@', nil, tcell.StyleDefault)
	for _, p := range gd.snake.tail {
		s.SetContent(p.x, p.y, '·', nil, tcell.StyleDefault)
	}
	s.SetContent(gd.food.x, gd.food.y, '+', nil, tcell.StyleDefault)
	drawBorder(s)
	drawHelp(s)
	if gd.gameOver {
		drawGameOver(s)
	}
	s.Show()
}

func drawHelp(s tcell.Screen) {
	helpStr := "Press Esc to exit."
	w, _ := s.Size()
	for i, c := range helpStr {
		s.SetContent(w-len(helpStr)+i-2, 1, c, nil, tcell.StyleDefault)
	}
}

func drawBorder(s tcell.Screen) {
	w, h := s.Size()
	drawBox(s, 0, 0, w-1, h-1)
}

func drawGameOver(s tcell.Screen) {
	w, h := s.Size()
	message := "Game over! Press R to restart."
	boxLeft := w/2 - len(message)/2 - 3
	boxTop := h/2 - 1
	drawBox(s, boxLeft, boxTop, len(message)+4, 2)
	for i, c := range message {
		s.SetContent(boxLeft+2+i, boxTop+1, c, nil, tcell.StyleDefault)
	}
}

func drawBox(s tcell.Screen, startX, startY, width, height int) {
	s.SetContent(startX, startY, tcell.RuneULCorner, nil, tcell.StyleDefault)
	s.SetContent(startX, startY+height, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	s.SetContent(startX+width, startY, tcell.RuneURCorner, nil, tcell.StyleDefault)
	s.SetContent(startX+width, startY+height, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	for y := startY; y <= startY+height; y += height {
		for x := startX + 1; x < width; x++ {
			s.SetContent(x, y, tcell.RuneHLine, nil, tcell.StyleDefault)
		}
	}
	for x := startX; x <= startX+width; x += width {
		for y := startY + 1; y < height; y++ {
			s.SetContent(x, y, tcell.RuneVLine, nil, tcell.StyleDefault)
		}
	}
}
