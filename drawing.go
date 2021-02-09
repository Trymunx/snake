package main

import (
	"github.com/gdamore/tcell/v2"
)

var snakeColour = tcell.StyleDefault

func (gd *gameData) draw(s tcell.Screen) {
	s.Clear()
	s.SetContent(gd.snake.head.x, gd.snake.head.y, '@', nil, snakeColour)
	for _, p := range gd.snake.tail {
		s.SetContent(p.x, p.y, '·', nil, snakeColour)
	}
	s.SetContent(gd.food.x, gd.food.y, '+', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
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
	helpStr := "Press P to pause/unpause. Press Esc to exit."
	w, _ := s.Size()
	for i, c := range helpStr {
		s.SetContent(w-len([]rune(helpStr))+i-2, 0, c, nil, tcell.StyleDefault)
	}
}

func drawBorder(s tcell.Screen) {
	w, h := s.Size()
	style := tcell.StyleDefault.Foreground(tcell.ColorDarkGrey.TrueColor())
	drawBox(s, 0, 1, w, h-1, style)
}

func drawMessage(s tcell.Screen, messages []string, messageStyle, boxStyle tcell.Style) {
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
	drawBox(s, boxLeft, boxTop, maxLength+4, len(messages)+2, boxStyle)
	for i, message := range messages {
		for j, c := range message {
			s.SetContent(boxLeft+2+j, boxTop+1+i, c, nil, messageStyle)
		}
	}
}

func drawPause(s tcell.Screen) {
	drawMessage(s, []string{
		"Game paused.",
		"Press P to unpause.",
	}, tcell.StyleDefault, tcell.StyleDefault.Foreground(tcell.ColorBlue))
}

func drawGameOver(s tcell.Screen) {
	drawMessage(s, []string{
		"Game over!",
		"Press R to restart.",
	}, tcell.StyleDefault, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func drawBox(s tcell.Screen, startX, startY, width, height int, style tcell.Style) {
	s.SetContent(startX, startY, tcell.RuneULCorner, nil, style)
	s.SetContent(startX, startY+height-1, tcell.RuneLLCorner, nil, style)
	s.SetContent(startX+width-1, startY, tcell.RuneURCorner, nil, style)
	s.SetContent(startX+width-1, startY+height-1, tcell.RuneLRCorner, nil, style)
	for y := startY; y <= startY+height-1; y += height - 1 {
		for x := startX + 1; x < startX+width-1; x++ {
			s.SetContent(x, y, tcell.RuneHLine, nil, style)
		}
	}
	for x := startX; x <= startX+width-1; x += width - 1 {
		for y := startY + 1; y < startY+height-1; y++ {
			s.SetContent(x, y, tcell.RuneVLine, nil, style)
		}
	}
}
