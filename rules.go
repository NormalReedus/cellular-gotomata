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

func NewConwaysGameOfLife() *ConwaysGameOfLife {
	return &ConwaysGameOfLife{Kernel{size: 3}}
}

//* CUSTOM GAME
type MyGameOfLife struct {
	Kernel
}

func (mgol MyGameOfLife) ApplyKernel(win *Window) *Dot {
	//* CUSTOM RULES HERE

	return NewDot(win.GridCoords(), nil)
}

func (mgol MyGameOfLife) Size() int {
	return mgol.size
}

func NewMyGameOfLife() *MyGameOfLife {
	return &MyGameOfLife{Kernel{size: 5}} //* CUSTOM KERNEL SIZE HERE
}
