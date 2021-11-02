package main

type Convolver interface {
	ApplyKernel(*Window) *Dot
	Size() int
}

type Kernel struct {
	size int
}

//* -------------------------
//* CONWAY'S GAME OF LIFE
//* -------------------------
//? A RULESET LOOKS LIKE THIS:
type ConwaysGameOfLife struct {
	Kernel
}

func (gol ConwaysGameOfLife) ApplyKernel(win *Window) *Dot { //? every kernel returns a dot or no dot for the center value of the window
	//? CUSTOM RULES HERE
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
	return &ConwaysGameOfLife{Kernel{size: 3}} //? CUSTOM KERNEL SIZE HERE
}

//* -------------------------
//* CUSTOM GAME 1
//* -------------------------
type CustomGame1 struct {
	Kernel
}

func (CustomGame1) ApplyKernel(win *Window) *Dot {
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

		return win.Center() // Alive cell defaults to 'stay alive'

	} else {
		// If cell is dead
		if contiguousAlive >= 2 {
			return NewDot(win.GridCoords(), nil)
		}

		return nil // Dead cell default to 'stay dead'
	}
}

func (gol CustomGame1) Size() int {
	return gol.size
}

func NewCustomGame1() Convolver {
	return &CustomGame1{Kernel{size: 3}}
}

//* -------------------------
//* CUSTOM GAME 2
//* -------------------------
type CustomGame2 struct {
	Kernel
}

func (CustomGame2) ApplyKernel(win *Window) *Dot {
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
		if diagonalAlive-contiguousAlive > 1 {
			return nil
		}

		return win.Center() // Alive cell defaults to 'stay alive'

	} else {
		// If cell is dead
		if contiguousAlive-diagonalAlive > 0 {
			return NewDot(win.GridCoords(), nil)
		}

		return nil // Dead cell default to 'stay dead'
	}
}

func (gol CustomGame2) Size() int {
	return gol.size
}

func NewCustomGame2() Convolver {
	return &CustomGame2{Kernel{size: 3}}
}

//* -------------------------
//* CUSTOM GAME MOD 1
//* -------------------------
type CustomGameMod1 struct {
	Kernel
}

func (CustomGameMod1) ApplyKernel(win *Window) *Dot {
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

		if contiguousAlive >= 3 {
			return nil
		}

		return win.Center() // Alive cell defaults to 'stay alive'

	} else {
		// If cell is dead
		if diagonalAlive >= 3 {
			return NewDot(win.GridCoords(), nil)
		}

		return nil // Dead cell default to 'stay dead'
	}
}

func (gol CustomGameMod1) Size() int {
	return gol.size
}

func NewCustomGameMod1() Convolver {
	return &CustomGameMod1{Kernel{size: 3}}
}
