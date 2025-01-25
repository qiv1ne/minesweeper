package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

// Represent a one minesweeper board cell.
type Cell struct {
	Empty       bool // Cell is empty when it don't contain the mine or number.
	MinesAround int  // Count of the mines around of this cell. If it 0 the Empty field should be true.

	// Flagged field represent is cell flagged of not.
	// I think this field have a right to exists. For example if in the future I will want to add hints in the game.
	Flagged bool

	IsMine   bool // If the cell contain mine it true.
	Revealed bool // If user reveal this cell it true.
}

// MineBoard struct defines game map.
type MineBoard struct {
	BoardConfig

	// RealBoard represent the game board with mines.
	RealBoard [][]Cell
	// UserBoard represent user's board with flags and not revealed mines.
	UserBoard [][]Cell
	// MinesRemain represent how many unflaged mines on the user board.
	MinesRemain int
}

// Struct BoardConfig contain essential options for create board.
type BoardConfig struct {
	Width  int   // Width of the board.
	Height int   // Height of the board.
	Mines  int   // Count of the mines.
	Seed   int64 // Seed for generating board
}

func NewSeed() int64 {
	return rand.New(rand.NewSource(time.Now().Unix())).Int63()
}

// The placeMines function accept *MineBoard pointer and place mines in the RealBoard.
// Like params function accept 1D board(you can use To1D() function) and count of mines
// It panic if error occurred.
func placeMines(oneDboard []Cell, minesCount int, seed int64) ([]Cell, error) {
	log.Info().Str("mines count", fmt.Sprintf("%d", minesCount)).Msg("Placing mines to board.")

	if len(oneDboard) < minesCount {
		return nil, fmt.Errorf("Board of size %d can't contain %d mines", len(oneDboard), minesCount)
	}
	if seed <= 0 {
		return nil, fmt.Errorf("Seed can't be <= 0. Use NewSeed() function to generate new seed.", len(oneDboard), minesCount)
	}

	var minesPlaced int

	// Create random generator with given seed.
	r := rand.New(rand.NewSource(seed))

	for minesPlaced < minesCount {
		// Transform board to 1D for using the length of it as limit for random function
		// Randomize the coordinate of the mine.
		x := r.Intn(len(oneDboard))

		if !oneDboard[x].IsMine {
			oneDboard[x].IsMine = true
			minesPlaced++
		}
	}
	log.Info().
		Str("mines count", fmt.Sprintf("%d", minesCount)).
		Msg("Successfully place all mines on the board.")

	return oneDboard, nil
}

// The PrintBroadGracefully is debug function to see how the board look.
func PrintBroadGracefully(board [][]Cell) {
	for _, row := range board {
		for _, cell := range row {
			if !cell.Revealed {
				fmt.Print("■ ")
				continue
			}
			if cell.Flagged {
				fmt.Print("⚑ ")
				continue
			}
			if cell.IsMine {
				fmt.Print("X ")
				continue
			}
			if cell.Empty {
				fmt.Print("□ ")
				continue
			}
			if cell.MinesAround != 0 {
				fmt.Print(cell.MinesAround)
				fmt.Print(" ")
				continue
			}
		}
		fmt.Println()
	}
}

// NewMineBoard function create filled MineBoard struct
func NewMineBoard(opts BoardConfig) (*MineBoard, error) {
	board, err := createBoard(opts)
	if err != nil {
		return nil, err
	}
	realBoard, err := createBoard(opts)
	realBoard, err = RevealAll(realBoard)
	if err != nil {
		return nil, err
	}
	return &MineBoard{
		BoardConfig: opts,
		RealBoard:   realBoard,
		UserBoard:   board,
		MinesRemain: opts.Mines,
	}, nil
}

// The CreateEmptyBoard create empty matrix of Cell.
func createBoard(opts BoardConfig) ([][]Cell, error) {
	log.Info().
		Str("Config", fmt.Sprintf("%v", opts)).
		Msg("Creating empty board.")
	// Creating empty board
	board := make([][]Cell, opts.Height)
	for i := range board {
		board[i] = make([]Cell, opts.Width)
	}

	if opts.Mines != 0 {
		board1d, err := To1D(board)
		if err != nil {
			return [][]Cell{}, err
		}
		boardWithMines, err := placeMines(board1d, opts.Mines, opts.Seed)
		if err != nil {
			return board, err
		}
		board, err = To2D(boardWithMines, opts.Height, opts.Width)
		if err != nil {
			return board, err
		}
		board, err = PlaceNumbers(board)
		if err != nil {
			return board, err
		}
	}
	return board, nil
}

// To1D transform given 2D board to a 1D board.
// First purpose of using the 1D board is to put a mines simpler.
// Function accept the 2D board and return 1D board.
// If the error occurred due the process of a board function is panic.
func To1D(board2D [][]Cell) ([]Cell, error) {
	log.Info().Msg("Converting 2D board to 1D")
	if len(board2D) == 0 {
		return make([]Cell, 0), fmt.Errorf("board is empty(len() = %d)", len(board2D))
	}
	// Create result 1D board with capacity of left side * top side
	board1D := make([]Cell, 0)

	for i := range board2D {
		for j := range board2D[i] {
			board1D = append(board1D, board2D[i][j])
		}
	}
	return board1D, nil
}

// To2D function transform the given 1D board to a 2D board.
// It accept the slice of Cell's, width of board, height of a board and return a matrix of Cell's.
func To2D(board1D []Cell, h, w int) ([][]Cell, error) {
	log.Info().Msg("Converting 1D board to 2D")
	if len(board1D) == 0 {
		return make([][]Cell, 0), fmt.Errorf("Board is empty(len() = %d)", len(board1D))
	}
	if h <= 0 {
		return make([][]Cell, 0), fmt.Errorf("Height( = %d) can't be <= 0", h)
	}
	if w <= 0 {
		return make([][]Cell, 0), fmt.Errorf("Wight( = %d) can't be <= 0", w)
	}
	// Create result 2D board.
	// The count of rows equals to h. The row length is equals to w.
	board2D := make([][]Cell, h)

	// This explanation is more for my better understanding.
	/*
		We have slice: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12].
		We want to create matrix with height 4 and width 3.
		For this we need to cut our slice in to 4 small slice's with 3 elements in it

		We iterate over empty matrix:


				i = 0. First iteration.
				Our empty row(i) accept three elements from slice:
				Elements starts on slice[i * w] and ends on slice[(i+1) * w]
				In out example: from slice[0 * 3 = 0] to slice[1 * 3 = 3]
				We have slice: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12].
								^-----^
				And like this our result now contain [1,2,3] in 0 row


				i = 1. Second iteration
				Get elements from slice[1 * 3 = 3] to slice[2 * 3 = 6]
				matrix[i] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12].
									  ^-----^
				Now matrix look like this:
					[
						[1, 2, 3],
						[4, 5, 6],
					]

				i = 2.
				matrix[2] = slice[2 * 3 = 6] : slice[3 * 3 = 9]
				matrix[2] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12].
											   ^-----^
	*/
	for i := range h {
		board2D[i] = board1D[i*w : (i+1)*w]
	}
	return board2D, nil
}

func RevealAll(board [][]Cell) ([][]Cell, error) {
	if len(board) == 0 {
		return make([][]Cell, 0), fmt.Errorf("Board is empty(len() = %d)", len(board))
	}
	log.Info().Msg("Revealing all cells")
	for i, row := range board {
		for j, cell := range row {
			cell.Revealed = true
			board[i][j] = cell
		}
	}
	return board, nil
}

func PlaceNumbers(board [][]Cell) ([][]Cell, error) {
	if len(board) == 0 {
		return [][]Cell{}, fmt.Errorf("Board can't be empty(len() = %d)", len(board))
	}
	// Another big comment for my understanding
	/*
			We have board:
		i  -------------------------
		0 | [X] [0] [0] [X] [0] [0] |
		1 | [0] [X] [0] [0] [X] [0] |
		2 | [0] [0] [X] [0] [0] [0] |
		3 | [0] [0] [0] [0] [0] [0] |
		4 | [X] [0] [0] [X] [0] [0] |
		   ------------------------
		j:   0	 1	 2	 3	 4	 5

			Every cell have this structure:
		i  -------------
		0 | [i-1 : j-1] [i-1 :   j] [i-1 : j+1] |
		1 | [i   : j-1]     {X}     [i   : j+1] |
		2 | [i+1 : j-1] [i+1 :   j] [i+1 : j+1] |
		   -----------------------------------
		j:       0  	     1		 	 2

			On every check of the neighbor cell we need to check are we going out of the board or not.

			Or we can do it simpler(I get it right now). We can check: if i == 0 then it's top corner and we don't need to check it.

			Ok, I was thinking about write functions that will check is it corner or not.
			But now I think function which will calculate sides will be better.

			I think about increasing count of mines of cells around mine,
			not to check how many mines around of cell.

			Yeah, i don't know how to do this better, I just check all direction.
	*/
	for i, row := range board {
		for j, cell := range row {
			if cell.IsMine {
				// check position [i : j+1]
				/*	■ ■ ■
					■ X 1
					■ ■ ■ */
				if j != len(row)-1 {
					board[i][j+1].MinesAround++
				}

				// check position [i : j-1]
				/*	■ ■ ■
					1 X ■
					■ ■ ■ */
				if j != 0 {
					board[i][j-1].MinesAround++
				}

				// check position [i-1 : j]
				/*	■ 1 ■
					■ X ■
					■ ■ ■ */
				if i != 0 {
					board[i-1][j].MinesAround++

					// check position [i-1 : j+1]
					/*	■ ■ 1
						■ X ■
						■ ■ ■ */
					if j != len(row)-1 {
						board[i-1][j+1].MinesAround++
					}

					// check position [i-1 : j-1]
					/*	1 ■ ■
						■ X ■
						■ ■ ■ */
					if j != 0 {
						board[i-1][j-1].MinesAround++
					}
				}

				// check position [i+1 : j]
				/*	■ ■ ■
					■ X ■
					■ 1 ■ */
				if i != len(board)-1 {
					board[i+1][j].MinesAround++

					// check position [i+1 : j+1]
					/*	■ ■ ■
						■ X ■
						■ ■ 1 */
					if j != len(row)-1 {
						board[i+1][j+1].MinesAround++
					}

					// check position [i+1 : j-1]
					/*	■ ■ ■
						■ X ■
						1 ■ ■ */
					if j != 0 {
						board[i+1][j-1].MinesAround++
					}
				}

			}
		}
	}
	for i, row := range board {
		for j, cell := range row {
			if !cell.IsMine && cell.MinesAround == 0 {
				board[i][j].Empty = true
			}
		}
	}
	return board, nil
}
