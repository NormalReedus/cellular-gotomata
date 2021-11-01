package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func leftClick() *Point {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return &Point{X: x / CELL_SIZE, Y: y / CELL_SIZE}
	}

	return nil
}

func rightClick() *Point {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		return &Point{X: x / CELL_SIZE, Y: y / CELL_SIZE}
	}

	return nil
}

func spaceKey() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func cKey() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyC)
}
