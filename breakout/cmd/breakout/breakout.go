package main

import (
	"github.com/jakecoffman/learnopengl/breakout"
	"github.com/jakecoffman/learnopengl/breakout/eng"
)

func main() {
	Breakout := &breakout.Game{}
	eng.Run(Breakout, 800, 600)
}
