package breakout

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type Level struct {
	Bricks []*Object
}

func NewLevel() *Level {
	return &Level{
		Bricks: []*Object{},
	}
}

func (l *Level) Load(file string, lvlWidth, lvlHeight int) error {
	l.Bricks = l.Bricks[:0]

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	tileData := [][]int{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		row := []int{}
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("Failed to parse level: %s", err.Error())
			}
			row = append(row, i)
		}
		tileData = append(tileData, row)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Failed to scan file: %s", err)
	}
	if len(tileData) > 0 {
		return l.init(tileData, lvlWidth, lvlHeight)
	}
	return nil
}

func (l *Level) Draw(renderer *SpriteRenderer) {
	for _, tile := range l.Bricks {
		if tile.IsSolid && !tile.Destroyed {
			tile.Draw(renderer)
		}
	}
}

func (l *Level) IsCompleted() bool {
	for _, tile := range l.Bricks {
		if !tile.IsSolid && !tile.Destroyed {
			return false
		}
	}
	return true
}

func (l *Level) init(tileData [][]int, lvlWidth, lvlHeight int) error {
	height := len(tileData)
	width := len(tileData[0])
	unitWidth := lvlWidth / width
	unitHeight := lvlHeight / height

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if tileData[y][x] == 1 {
				pos := Vec2(unitWidth*x, unitHeight*y)
				size := Vec2(unitWidth, unitHeight)
				obj := NewGameObject2(pos, size, ResourceManager.Texture("block_solid"), mgl32.Vec3{.8, .8, .7}, mgl32.Vec2{})
				obj.IsSolid = true
				l.Bricks = append(l.Bricks, obj)
			} else if tileData[y][x] > 1 {
				color := mgl32.Vec3{1, 1, 1}
				switch tileData[y][x] {
				case 2:
					color = mgl32.Vec3{.2, .6, 1}
				case 3:
					color = mgl32.Vec3{0, .7, 0}
				case 4:
					color = mgl32.Vec3{.8, .8, .4}
				case 5:
					color = mgl32.Vec3{1, .5, 0}
				}

				pos := Vec2(unitWidth*x, unitHeight*y)
				size := Vec2(unitWidth, unitHeight)
				l.Bricks = append(l.Bricks, NewGameObject2(pos, size, ResourceManager.Texture("block"), color, mgl32.Vec2{}))
			}
		}
	}

	return nil
}

func Vec2(x, y int) mgl32.Vec2 {
	return mgl32.Vec2{float32(x), float32(y)}
}
