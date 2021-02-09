package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	gameloop "github.com/kutase/go-gameloop"
)

type point struct {
	x int
	y int
}

func randomFoodLocation(w, h int) *point {
	return &point{
		// Make sure x coords are always a multiple of 2
		x: 1 + 2*rand.Intn((w-2)/2),
		y: 3 + rand.Intn(h-4),
	}
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
	screen   tcell.Screen
	gameloop *gameloop.GameLoop
	paused   bool
}

func newSnake(x, y int) *snake {
	const snakeLength = 5
	head := &point{x: x, y: y}
	tail := []*point{}
	for i := 1; i <= snakeLength; i++ {
		tail = append(tail, &point{x: x - (2 * i), y: y})
	}
	return &snake{head, tail}
}

func newGameData(s tcell.Screen) *gameData {
	w, h := s.Size()
	gl := gameloop.New(10, nil)
	gd := &gameData{
		// The x-coord is gross to make it line up on 2n+1 grid lines that food generates on.
		snake:    newSnake(w/2+(w/2)%2-1, h/2),
		dir:      &point{x: 2, y: 0},
		food:     randomFoodLocation(w, h),
		gameOver: false,
		screen:   s,
		gameloop: gl,
	}
	gd.setUpdateFunc()
	return gd
}

func (gd *gameData) setUpdateFunc() {
	gd.gameloop.SetOnUpdate(func(delta float64) {
		if gd.paused {
			return
		}
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
			gd.food = randomFoodLocation(gd.screen.Size())
		}
		gd.snake.head.x += gd.dir.x
		gd.snake.head.y += gd.dir.y
		w, h := gd.screen.Size()
		if gd.snake.head.x < 1 || gd.snake.head.x >= w-1 || gd.snake.head.y < 2 || gd.snake.head.y >= h-1 {
			gd.gameOver = true
			gd.draw(gd.screen)
			gd.gameloop.Stop()
		}
		for i := 0; i < len(gd.snake.tail); i++ {
			if gd.snake.head.x == gd.snake.tail[i].x && gd.snake.head.y == gd.snake.tail[i].y {
				gd.gameOver = true
				gd.draw(gd.screen)
				gd.gameloop.Stop()
			}
		}
		gd.draw(gd.screen)
	})
}

func (gd *gameData) start() {
	gd.gameloop.Start()
}

func (gd *gameData) stop() {
	gd.gameloop.Stop()
}
