package main

import (
	"github.com/jakecoffman/learnopengl/breakout"
	"github.com/jakecoffman/learnopengl/breakout/eng"
)

var resolutions = []struct{x, y int} {
	{800, 600},
	{1024, 768},
	{1680, 1050},
	{1920, 1080},
}

func main() {
	Breakout := &breakout.Game{}
	const choice = 3
	x, y := resolutions[choice].x, resolutions[choice].y
	eng.Run(Breakout, x, y)
}
