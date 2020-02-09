package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

var (
	normalStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	activeStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
)

func runUI(items []*commandConfig) (*commandConfig, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	defer screen.Fini()

	selectedIndex := 0

	for {
		// Render UI
		screen.Clear()
		writeUILine(screen, activeStyle, 0, "Select a command to run:")
		for i, item := range items {
			item.render(screen, i+2, i == selectedIndex)
		}
		_, height := screen.Size()
		writeUILine(screen, activeStyle, height-1, "Use arrow keys to navigate, <Enter> to select and <Esc> to exit")
		screen.Show()

		// Handle an event
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {

			// Arrow up
			case tcell.KeyUp:
				selectedIndex--
				if selectedIndex < 0 {
					selectedIndex = 0
				}
				break

			// Arrow down
			case tcell.KeyDown:
				selectedIndex++
				if selectedIndex >= len(items) {
					selectedIndex = len(items) - 1
				}
				break

			// Home
			case tcell.KeyHome:
				selectedIndex = 0
				break

			// End
			case tcell.KeyEnd:
				selectedIndex = len(items) - 1
				break

			// Enter
			case tcell.KeyEnter:
				return items[selectedIndex], nil

			// Esc/Q/CtrlC
			case tcell.KeyEscape:
				return nil, nil
			case tcell.KeyCtrlQ:
				return nil, nil
			case tcell.KeyCtrlC:
				return nil, nil
			}

			if ev.Rune() == 'q' || ev.Rune() == 'Q' {
				return nil, nil
			}
		}
	}
}

func (t *commandConfig) render(screen tcell.Screen, line int, selected bool) {
	var style tcell.Style
	var text string
	if selected {
		style = activeStyle
		text = fmt.Sprintf("[x] %s", t.Name)
	} else {
		style = normalStyle
		text = fmt.Sprintf("[ ] %s", t.Name)
	}

	writeUILine(screen, style, line, text)
}

func writeUILine(screen tcell.Screen, style tcell.Style, line int, text string) {
	for i, c := range text {
		screen.SetContent(i, line, c, nil, style)
	}

	width, _ := screen.Size()
	for i := len(text); i < width; i++ {
		screen.SetContent(i, line, ' ', nil, style)
	}
}
