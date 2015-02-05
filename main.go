package main

import (
	"fmt"
	"sudoku_solver/sudoku_solver"
)

// Quick manual test
func main() {

	// This is apparantly the hardest Sudoku puzzle there is?
	// Easy!
	// LOL :D
	board := [][]int{
		{8, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 3, 6, 0, 0, 0, 0, 0},
		{0, 7, 0, 0, 9, 0, 2, 0, 0},
		{0, 5, 0, 0, 0, 7, 0, 0, 0},
		{0, 0, 0, 0, 4, 5, 7, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 3, 0},
		{0, 0, 1, 0, 0, 0, 0, 6, 8},
		{0, 0, 8, 5, 0, 0, 0, 1, 0},
		{0, 9, 0, 0, 0, 0, 4, 0, 0}}

	sp := sudoku_solver.SudokuPuzzle{Board: board}
	sp.Solve()

	fmt.Println("\nThe solution to this board is:\n")
	fmt.Println(sp.DrawAsciiBoard())
}
