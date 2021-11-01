package main

type Convolver interface {
	ApplyKernel(*Window) *Dot
	Size() int
}

type Kernel struct {
	size int
}

//* A RULESET LOOKS LIKE THIS:
type ConwaysGameOfLife struct {
	Kernel
}

func (gol ConwaysGameOfLife) ApplyKernel(win *Window) *Dot {
	//* https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
	dead := win.NumEmptyNeighbors()
	alive := 8 - dead

	currentVal := win.Center()

	// "Any live cell with two or three live neighbours survives."
	if currentVal != nil && between(alive, 2, 3) {
		return currentVal
	}

	// "Any dead cell with three live neighbours becomes a live cell."
	if currentVal == nil && alive == 3 {
		return NewDot(win.GridCoords(), nil) // parentGrid is set in grid.Convolve instead
	}

	// "All other live cells die in the next generation." / "...all other dead cells stay dead."
	return nil
}

func (gol ConwaysGameOfLife) Size() int {
	return gol.size
}

func NewConwaysGameOfLife() Convolver {
	return &ConwaysGameOfLife{Kernel{size: 3}}
}

//* CUSTOM GAME
type MyGameOfLife struct {
	Kernel
}

func (mgol MyGameOfLife) ApplyKernel(win *Window) *Dot {
	//* CUSTOM RULES HERE
	diagonal := [4]*Dot{win.Get(*NewPoint(0, 0)), win.Get(*NewPoint(2, 0)), win.Get(*NewPoint(0, 2)), win.Get(*NewPoint(2, 2))}
	contiguous := [4]*Dot{win.Get(*NewPoint(1, 0)), win.Get(*NewPoint(2, 1)), win.Get(*NewPoint(1, 2)), win.Get(*NewPoint(0, 1))}

	var diagonalAlive int
	for _, dot := range diagonal {
		if dot != nil {
			diagonalAlive++
		}
	}

	var contiguousAlive int
	for _, dot := range contiguous {
		if dot != nil {
			contiguousAlive++
		}
	}

	// If cell is alive
	if win.Center() != nil {
		if contiguousAlive == 4 {
			return nil
		}

		if diagonalAlive >= 3 {
			return nil
		}

		return win.Center()
	} else {
		// If cell is dead
		if contiguousAlive >= 2 {
			return NewDot(win.GridCoords(), nil)
		}

	}

	return nil
}

func (mgol MyGameOfLife) Size() int {
	return mgol.size
}

func NewMyGameOfLife() Convolver {
	return &MyGameOfLife{Kernel{size: 3}} //* CUSTOM KERNEL SIZE HERE
}
