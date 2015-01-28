package sudoku_solver

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Our basic Sudoku puzzle construction, went for an OO approach
type SudokuPuzzle struct {
	Board [][]int
}

// Public method to solve, generates the start point of possible values and then
// passes it off to the search member to be solved.
func (sp *SudokuPuzzle) Solve() error {

	_, possibleValues := sp.generateGridOfPossibleValues(sp.Board)

	// Reduce the values early before running
	possibleValues, contradiction := sp.search(sp.runEliminations(possibleValues))

	if contradiction {
		return errors.New("Board has no possible solutions")
	}

	sp.Board = sp.flattenPossibleValues(possibleValues)

	return nil
}

// Prints out the Sudoku board
func (sp *SudokuPuzzle) DrawAsciiBoard() string {

	boardPrintout := ""
	line := ""

	for i := 0; i < 9; i++ {
		if i != 0 && i%3 == 0 {
			line += fmt.Sprintf("-+-")
		}
		line += strings.Repeat("-", 3)
	}
	fmt.Sprintf(line, "\n")

	for i := 0; i < 9; i++ {

		if i != 0 && i%3 == 0 {
			boardPrintout += fmt.Sprintf(line)
			boardPrintout += fmt.Sprintf("\n")
		}

		for j := 0; j < 9; j++ {
			if j != 0 && j%3 == 0 {
				boardPrintout += fmt.Sprintf(" | ")
			}
			boardPrintout += fmt.Sprintf(" %d ", sp.Board[i][j])
		}
		boardPrintout += fmt.Sprintf("\n")
	}

	return boardPrintout
}

// Creates a start point, will be a 3D array where the third dimension holds
// all possible values. For numbers that are already placed, that is the only
// possible value - for all others they have possible values 1-9.
// We will hand this off for elimination and assignment to attempt to get
// our 3D array into a 2D array (where the third dimension is 1 value).
func (sp *SudokuPuzzle) generateGridOfPossibleValues(board [][]int) (error, [][][]int) {

	// Generate a 3D array
	possibleValues := make([][][]int, len(board))

	for i := range possibleValues {
		possibleValues[i] = make([][]int, len(board))
		for j := range possibleValues[i] {

			if board[i][j] > 0 {
				// If a value exists in the place already, that is of course the
				// only possible value
				possibleValues[i][j] = []int{board[i][j]}
			} else {
				// Otherwise it's possible values range from 1 to 9
				for n := 1; n < 10; n++ {

					// If n exists in the current row, column or square
					possibleValues[i][j] = append(possibleValues[i][j], n)
				}
			}
		}
	}

	return nil, possibleValues
}

// Flatten grid of possible values to a standard board
func (sp *SudokuPuzzle) flattenPossibleValues(possibleValues [][][]int) [][]int {

	board := make([][]int, len(possibleValues))

	// Create a clone
	for i := 0; i < 9; i++ {
		board[i] = make([]int, len(possibleValues[i]))
		for j := 0; j < 9; j++ {
			board[i][j] = possibleValues[i][j][0]
		}
	}

	return board
}

// Search is the magical recursive function. The process is something like this:
// If the board is solved or has an immediate contradiction, return out
// Otherwise, create a copy of our board. We find the easiest starting point
// which is the square with the lowest number of possibilities.
//
// We loop the possibilities, and with each of them we try eliminating based
// on the rules of Sudoku.
//
// If the elimination resulted in any square being empty, then of course that
// breaks the rules and it cannot possibly be correct. If there is no
// contradiction then there's a chance it could be correct.
//
// If it is incorrect, we remove that from all possibilities. We then run a
// new search on that and look for any other contradictions. If there are any
// contradictions at all in the tree, they all return true and we know to try
// another branch.
//
// If it is correct, we leave it as is and run a new search
func (sp *SudokuPuzzle) search(board [][][]int, contradiction bool) ([][][]int, bool) {

	// Find the minimum variable (the start point) and try eliminating some
	// of the values...
	minI, minJ := sp.findMinimum(board)

	// Board has contradiction
	if sp.hasContradiction(board) {
		return board, true
	}

	// If it is solved then return the board
	if sp.isSolved(board) {
		return board, false
	}

	// Create a copy of the board (this is a deep copy)
	boardCopy := sp.deepCopyBoard(board)

	// Set up our test value
	for _, testValue := range board[minI][minJ] {

		// Create a clone of the board, with the value we're testing in place
		boardCopy[minI][minJ] = []int{testValue}

		// Run our elimination loop
		boardCopy, contradiction = sp.runEliminations(boardCopy)

		if contradiction {

			// In the case of an assignment contradiction, remove it from the
			// possibilities and set the board to continue on
			board[minI][minJ] =
				deleteIntFromSlice(board[minI][minJ], testValue)

			boardCopy = board
		} else {

			// Because we have eliminated a value, we need to look down at the
			// possibilities. If there is a contradiction then we must rewind
			boardCopy, contradiction = sp.search(boardCopy, false)

			if !contradiction {
				return boardCopy, false
			}
			boardCopy = board
		}
	}

	// Branch out!
	return sp.search(boardCopy, false)
}

// Finds what will be the easiest start point in the board, which is that with
// the lowest number of possibilities
func (sp *SudokuPuzzle) findMinimum(board [][][]int) (int, int) {

	minI := 0
	minJ := 0

	minLen := 10

	// Find what has the fewest possibilities to be the start point
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if len(board[i][j]) > 1 && len(board[i][j]) < minLen {
				minLen = len(board[i][j])
				minI = i
				minJ = j
			}
		}
	}
	return minI, minJ
}

// If there is a contradiction then we return true, else false.
// A contradiction means that there is a tile which has no possible values
func (sp *SudokuPuzzle) hasContradiction(board [][][]int) bool {

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if len(board[i][j]) == 0 {
				return true
			}
		}
	}
	return false
}

// Perform a really basic check to see if the game is solved. A game is solved
// if all squares have been reduced to a length of 1
func (sp *SudokuPuzzle) isSolved(board [][][]int) bool {

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if len(board[i][j]) != 1 {
				return false
			}
		}
	}
	return true
}

// Creates a deep copy of the board, important to make and copy at
// every level!
func (sp *SudokuPuzzle) deepCopyBoard(board [][][]int) [][][]int {

	// Create a clone
	newBoard := make([][][]int, len(board))
	copy(newBoard, board)
	for i := 0; i < 9; i++ {
		newBoard[i] = make([][]int, len(board[i]))
		copy(newBoard[i], board[i])
		for j := 0; j < 9; j++ {
			newBoard[i][j] = make([]int, len(board[i][j]))
			copy(newBoard[i][j], board[i][j])
		}
	}

	return newBoard
}

// Loops through each tile and eliminates impossible solutions
func (sp *SudokuPuzzle) runEliminations(board [][][]int) ([][][]int, bool) {

	var contradiction bool

	// Reduce possibilities
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {

			// Eliminate values for each square on the board
			board, contradiction = sp.eliminate(board, y, x)

			// If any contradictions came back then break out and
			// try the next value
			if contradiction {
				return make([][][]int, 0), true
			}
		}
	}
	return board, false
}

// Performs an elimination on a square. It will loop through the column, row and
// square of which it resides.
func (sp *SudokuPuzzle) eliminate(board [][][]int, y int, x int) ([][][]int, bool) {

	for i := 0; i < 9; i++ {

		// Eliminate values in the current column
		if len(board[i][x]) == 1 &&
			i != y && intInSlice(board[y][x], board[i][x][0]) {

			board[y][x] = deleteIntFromSlice(board[y][x], board[i][x][0])

			// Check if it has caused a contradiction and return early
			if len(board[y][x]) == 0 {
				// Contradiction, tis not possible!
				return make([][][]int, 0), true
			}
		}

		// Eliminate values in the current row
		if len(board[y][i]) == 1 &&
			i != x && intInSlice(board[y][x], board[y][i][0]) {

			board[y][x] = deleteIntFromSlice(board[y][x], board[y][i][0])

			// Check if it has caused a contradiction and return early
			if len(board[y][x]) == 0 {
				// Contradiction, tis not possible!
				return make([][][]int, 0), true
			}
		}

	}

	// Calculate values for our square search
	sqX, sqY := sp.getSquare(x, y)

	deltaX := sqX * 3
	deltaY := sqY * 3

	// Run deletion on the square
	for i := deltaY; i < deltaY+3; i++ {
		for j := deltaX; j < (deltaX)+3; j++ {

			if len(board[i][j]) == 1 && i != y && x != j &&
				intInSlice(board[y][x], board[i][j][0]) {

				board[y][x] =
					deleteIntFromSlice(board[y][x], board[i][j][0])

				// Check if it has caused a contradiction and return early
				if len(board[y][x]) == 0 {
					// Contradiction, tis not possible!
					return make([][][]int, 0), true
				}
			}
		}
	}

	return board, false
}

// Gets the square position for a given x and y
func (sp *SudokuPuzzle) getSquare(x int, y int) (int, int) {

	return int(math.Floor(float64(x / 3))), int(math.Floor(float64(y / 3)))
}

//---------------------------------------------------------------------------//

// A helpful function which will remove an integer from a slice by recreating
// the slice. This means we always have the 0 key for convenience.
func deleteIntFromSlice(slice []int, elim int) []int {

	for key, value := range slice {
		if value == elim {
			return append(slice[:key], slice[key+1:]...)
		}
	}

	return slice
}

// We check if the int exists in the slice
func intInSlice(slice []int, value int) bool {

	for _, sliceValue := range slice {
		if sliceValue == value {
			return true
		}
	}
	return false
}
