package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
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

	gd := newGameData(s)
	gd.start()

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
				if gd.dir.y != 1 && !gd.paused {
					gd.dir = &point{x: 0, y: -1}
				}
			case tcell.KeyDown:
				if gd.dir.y != -1 && !gd.paused {
					gd.dir = &point{x: 0, y: 1}
				}
			case tcell.KeyLeft:
				if gd.dir.x != 2 && !gd.paused {
					gd.dir = &point{x: -2, y: 0}
				}
			case tcell.KeyRight:
				if gd.dir.x != -2 && !gd.paused {
					gd.dir = &point{x: 2, y: 0}
				}
			}
			if ev.Rune() == 'p' || ev.Rune() == 'P' {
				gd.paused = !gd.paused
			}
			if (ev.Rune() == 'r' || ev.Rune() == 'R') && gd.gameOver {
				w, h := s.Size()
				gd.snake = newSnake(w/2, h/2)
				gd.food = randomPoint(w, h)
				gd.dir = &point{x: 2, y: 0}
				gd = newGameData(s)
				gd.start()
			}
			if ev.Rune() == 'q' || ev.Rune() == 'Q' {
				s.Fini()
				os.Exit(0)
			}
			gd.draw(s)
		}
	}
}
