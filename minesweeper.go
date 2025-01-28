package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	Lose = 1 << iota
	Win
)

// Represent a one minesweeper board cell.
type Cell struct {
	// Cell is empty when it don't contain the mine or number.
	Empty bool

	// Count of the mines around of this cell. If it 0 the Empty field should be true.
	MinesAround int

	// Flagged field represent is cell flagged of not.
	// I think this field have a right to exists.
	// For example if in the future I will want to add hints in the game.
	Flagged bool

	IsMine   bool // If the cell contain mine it true.
	Revealed bool // If user reveal this cell it true.
}

// board is alias for [][]Cell type. I create it for linking functions to type.
type board [][]Cell

// MineBoard struct defines game map.
type MineBoard struct {
	BoardConfig

	// Real represent the game board with mines.
	Real board
	// User represent user's board with flags and not revealed mines.
	User board
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

// NewSeed function generate int64 seed from time.Now().Unix()
func NewSeed() int64 {
	return rand.New(rand.NewSource(time.Now().Unix())).Int63()
}

// The placeMines function place mines in Board.
// Like params function accept 1D board(you can use To1D() function) and count of mines
// It panic if error occurred.
func (b board) placeMines(minesCount int, seed int64) {
	// Create random generator with given seed.
	r := rand.New(rand.NewSource(seed))

	// Create map for saving unique places for mines.
	mines := make(map[int]struct{}, minesCount)
	for len(mines) < minesCount {
		rand := r.Intn(len(b) * len(b[0]))
		// If place not in the map:
		if _, ok := mines[rand]; !ok {
			mines[rand] = struct{}{}
		}
	}

	// Every iteration increase c variable for count when c will equal number of place in the map.
	var c int
	for i := range b {
		for j := range b[i] {
			// If c is equals number in the map:
			if _, ok := mines[c]; ok {
				b[i][j].IsMine = true
			}
			c++
		}
	}
}

// The PrintBroadGracefully is debug function to see how the board look.
func (b board) Print() {
	for _, row := range b {
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
func NewMineBoard(config BoardConfig) (*MineBoard, error) {
	switch {
	case config.Height <= 1:
		return nil, errors.New("height can't be less than 2")
	case config.Width <= 1:
		return nil, errors.New("width can't be less than 2")
	case config.Mines <= 0:
		return nil, errors.New("count of mines can't be less than 1")
	case config.Seed <= 0:
		err := errors.New("Seed can't be <= 0. Use NewSeed() function to generate random seed.")
		return nil, err
	case config.Height*config.Width <= config.Mines:
		return nil, errors.New("Board can't contain mines more than it length")
	}

	u, err := createBoard(config)
	if err != nil {
		return nil, err
	}
	r, err := createBoard(config)
	if err != nil {
		return nil, err
	}
	err = r.revealAll()
	if err != nil {
		return nil, err
	}
	return &MineBoard{
		BoardConfig: config,
		Real:        r,
		User:        u,
		MinesRemain: config.Mines,
	}, nil
}

// The createBoard create matrix of Cell and fill it with mines and numbers.
func createBoard(config BoardConfig) (board, error) {
	board := make(board, config.Height)
	for i := range board {
		board[i] = make([]Cell, config.Width)
	}
	board.placeMines(config.Mines, config.Seed)
	board.placeNumbers()

	return board, nil
}

// To1D transform given 2D board to a 1D board.
// First purpose of using the 1D board is to put a mines simpler.
// Function accept the 2D board and return 1D board.
// If the error occurred due the process of a board function is panic.
func (b board) To1D() ([]Cell, error) {
	if len(b) == 0 {
		return make([]Cell, 0), errors.New("board is empty")
	}
	// Create result 1D board with capacity of left side * top side
	board1D := make([]Cell, len(b)*len(b[0]))

	for i := range b {
		for j := range b[i] {
			board1D[i+j] = b[i][j]
		}
	}
	return board1D, nil
}

// RevealAll function set all Cell.Revealed to true
func (b board) revealAll() error {
	for i := range b {
		for j := range b[i] {
			b[i][j].Revealed = true
		}
	}
	return nil
}

// placeNumbers function calculate count of mines around every not mine Cell
func (b board) placeNumbers() {
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
	for i, row := range b {
		for j, cell := range row {
			if cell.IsMine {
				// check position [i : j+1]
				/*	■ ■ ■
					■ X 1
					■ ■ ■ */
				if j != len(row)-1 {
					b[i][j+1].MinesAround++
				}

				// check position [i : j-1]
				/*	■ ■ ■
					1 X ■
					■ ■ ■ */
				if j != 0 {
					b[i][j-1].MinesAround++
				}

				// check position [i-1 : j]
				/*	■ 1 ■
					■ X ■
					■ ■ ■ */
				if i != 0 {
					b[i-1][j].MinesAround++

					// check position [i-1 : j+1]
					/*	■ ■ 1
						■ X ■
						■ ■ ■ */
					if j != len(row)-1 {
						b[i-1][j+1].MinesAround++
					}

					// check position [i-1 : j-1]
					/*	1 ■ ■
						■ X ■
						■ ■ ■ */
					if j != 0 {
						b[i-1][j-1].MinesAround++
					}
				}

				// check position [i+1 : j]
				/*	■ ■ ■
					■ X ■
					■ 1 ■ */
				if i != len(b)-1 {
					b[i+1][j].MinesAround++

					// check position [i+1 : j+1]
					/*	■ ■ ■
						■ X ■
						■ ■ 1 */
					if j != len(row)-1 {
						b[i+1][j+1].MinesAround++
					}

					// check position [i+1 : j-1]
					/*	■ ■ ■
						■ X ■
						1 ■ ■ */
					if j != 0 {
						b[i+1][j-1].MinesAround++
					}
				}

			}
		}
	}
	for i, row := range b {
		for j, cell := range row {
			if !cell.IsMine && cell.MinesAround == 0 {
				b[i][j].Empty = true
			}
		}
	}
}

// OpenCell function open cell in user board[y][x]
// If it mine: return Lose const
// If nothing happend: return 0
// Return -1 and error if x or y out of range
func (board *MineBoard) OpenCell(x, y int) (int, error) {
	if board == nil {
		return -1, errors.New("board is nil")
	}
	if len(board.User) <= 1 {
		return -1, errors.New("board is too small, less than 2 row")
	}
	if len(board.User[0]) <= 1 {
		return -1, errors.New("board is too small, less than 2 column")
	}
	if x > len(board.User[0]) || x <= 0 {
		return -1, errors.New("x is out of the row")
	}
	if y > len(board.User) || y <= 0 {
		return -1, errors.New("y is out of the row")
	}
	if board.User[y-1][x-1].IsMine {
		return Lose, nil
	}
	board.User[y-1][x-1].Revealed = true
	return 0, nil
}

// PlaceFlag function mark cell in user board[y][x].
// If it last mine: return Win const.
// Return -1 and error if x or y out of range.
func (board *MineBoard) PlaceFlag(x, y int) (int, error) {
	if x > len(board.User[0]) || x <= 0 {
		return -1, errors.New("x is out of the row")
	}
	if y > len(board.User) || y <= 0 {
		return -1, errors.New("y is out of the row")
	}

	if board.MinesRemain <= 1 {
		return Win, nil
	}
	board.User[y][x].Flagged = true
	board.MinesRemain--
	return 0, nil
}
