package main

import (
	"github.com/gdamore/tcell/v2"
)

func (gd *gameData) draw(s tcell.Screen) {
	s.Clear()
	s.SetContent(gd.snake.head.x, gd.snake.head.y, '@', nil, tcell.StyleDefault)
	for _, p := range gd.snake.tail {
		s.SetContent(p.x, p.y, '·', nil, tcell.StyleDefault)
	}
	s.SetContent(gd.food.x, gd.food.y, '+', nil, tcell.StyleDefault)
	drawBorder(s)
	drawHelp(s)
	if gd.paused {
		drawPause(s)
	}
	if gd.gameOver {
		drawGameOver(s)
	}
	s.Show()
}

func drawHelp(s tcell.Screen) {
	helpStr := "Press Esc to exit."
	w, _ := s.Size()
	for i, c := range helpStr {
		s.SetContent(w-len([]rune(helpStr))+i-2, 1, c, nil, tcell.StyleDefault)
	}
}

func drawBorder(s tcell.Screen) {
	w, h := s.Size()
	drawBox(s, 0, 0, w, h)
}

func drawMessage(s tcell.Screen, messages []string) {
	w, h := s.Size()
	var maxLength int
	for _, m := range messages {
		length := len([]rune(m))
		if length > maxLength {
			maxLength = length
		}
	}
	boxLeft := w/2 - maxLength/2 - 1
	boxTop := h/2 - 1 - len(messages)
	drawBox(s, boxLeft, boxTop, maxLength+4, len(messages)+2)
	for i, message := range messages {
		for j, c := range message {
			s.SetContent(boxLeft+2+j, boxTop+1+i, c, nil, tcell.StyleDefault)
		}
	}
}

func drawPause(s tcell.Screen) {
	drawMessage(s, []string{
		"Game paused.",
		"Press P to unpause.",
	})
}

func drawGameOver(s tcell.Screen) {
	drawMessage(s, []string{
		"Game over!",
		"Press R to restart.",
	})
}

func drawBox(s tcell.Screen, startX, startY, width, height int) {
	s.SetContent(startX, startY, tcell.RuneULCorner, nil, tcell.StyleDefault)
	s.SetContent(startX, startY+height-1, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	s.SetContent(startX+width-1, startY, tcell.RuneURCorner, nil, tcell.StyleDefault)
	s.SetContent(startX+width-1, startY+height-1, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	for y := startY; y <= startY+height-1; y += height - 1 {
		for x := startX + 1; x < startX+width-1; x++ {
			s.SetContent(x, y, tcell.RuneHLine, nil, tcell.StyleDefault)
		}
	}
	for x := startX; x <= startX+width-1; x += width - 1 {
		for y := startY + 1; y < startY+height-1; y++ {
			s.SetContent(x, y, tcell.RuneVLine, nil, tcell.StyleDefault)
		}
	}
}
