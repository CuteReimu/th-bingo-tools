package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
)

type gui struct {
	mainWindow *walk.MainWindow
}

func (g *gui) init() {
	screenX, screenY := getSystemMetrics(0), getSystemMetrics(1)
	width, height := screenX/2, screenY*3/4
	window := MainWindow{
		Title:    "th-bingo-tools",
		Font:     Font{PointSize: 12},
		Bounds:   Rectangle{X: (screenX - width) / 2, Y: (screenY - height) / 2, Width: width, Height: height},
		AssignTo: &g.mainWindow,
		Layout:   VBox{},
		Children: []Widget{
			Slider{
				MinValue: 1,
				MaxValue: 1000,
				Value:    100,
			},
		},
	}
	ch := make(chan struct{})
	go func() {
		err := window.Create()
		if err != nil {
			panic(err)
		}
		close(ch)
		code := g.mainWindow.Run()
		os.Exit(code)
	}()
	<-ch
}
