package main

import "sudoku_solver/sudoku_solver"

// Quick manual test
func main() {

	board := [][]int{
		{8, 5, 0, 0, 0, 2, 4, 0, 0},
		{7, 2, 0, 0, 0, 0, 0, 0, 9},
		{0, 0, 4, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 7, 0, 0, 2},
		{3, 0, 5, 0, 0, 0, 9, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 8, 0, 0, 7, 0},
		{0, 1, 7, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 3, 6, 0, 4, 0}}

	sp := sudoku_solver.SudokuPuzzle{Board: board}
	sp.Solve()
	sp.PrintBoard()
}
