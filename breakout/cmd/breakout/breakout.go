package main

import (
	"github.com/jakecoffman/learnopengl/breakout"
	"github.com/jakecoffman/learnopengl/breakout/eng"
)

const (
	width  = 800
	height = 600
)

func main() {
	Breakout := &breakout.Game{}
	eng.Run(Breakout, width, height)
}
