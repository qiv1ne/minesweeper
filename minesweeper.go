package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/qiv1ne/log"
)

var (
	logger, _ = log.New(os.Stdin, &log.Opts{
		FuncName:   true,
		FileName:   true,
		LineNumber: true,
		Date:       false,
		Time:       false,
	})
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

// Board is alias for [][]Cell type. I create it for linking functions to type.
type Board [][]Cell

// MineBoard struct defines game map.
type MineBoard struct {
	BoardConfig

	// Real represent the game board with mines.
	Real Board
	// User represent user's board with flags and not revealed mines.
	User Board
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

// The placeMines function place mines in Board.
// Like params function accept 1D board(you can use To1D() function) and count of mines
// It panic if error occurred.
func (b Board) placeMines(minesCount int, seed int64) error {
	logger.Print(log.Info("Placing mines"))
	if seed <= 0 {
		err := errors.New("Seed can't be <= 0. Use NewSeed() function to generate new seed.")
		logger.Print(log.Error(err))
		return err
	}
	logger.Print(log.Info("check count of mines"), log.Str(strconv.Itoa(len(b)), strconv.Itoa(len(b[0]))))
	if len(b)*len(b[0]) <= minesCount {
		err := errors.New("Board can't contain mines more than it length")
		logger.Print(log.Error(err))
		return err
	}

	// Create random generator with given seed.
	r := rand.New(rand.NewSource(seed))

	mines := make(map[int]struct{}, minesCount)
	for len(mines) < minesCount {
		rand := r.Intn(len(b) * len(b[0]))
		if _, ok := mines[rand]; !ok {
			mines[rand] = struct{}{}
		}
	}

	logger.Print(log.Info("Mines calculated"),
		log.Str("placed mines", fmt.Sprintf("%v", mines)))
	var c int
	for i := range b {
		for j := range b[i] {
			// logger.Print(log.Info("Check cell"), log.Str("ok", fmt.Sprintf("%v", mines[c])))
			if _, ok := mines[c]; ok {
				b[i][j].IsMine = true
				c++
				logger.Print(log.Info("Mine placed"),
					log.Str("row", fmt.Sprintf("%d", i)),
					log.Str("column", fmt.Sprintf("%d", j)))
			}
		}
	}
	return nil
}

// The PrintBroadGracefully is debug function to see how the board look.
func (b Board) Print() {
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
func NewMineBoard(opts BoardConfig) (*MineBoard, error) {
	user, err := createBoard(opts)
	if err != nil {
		return nil, err
	}
	real, err := createBoard(opts)
	if err != nil {
		return nil, err
	}
	err = real.RevealAll()
	if err != nil {
		return nil, err
	}
	return &MineBoard{
		BoardConfig: opts,
		Real:        real,
		User:        user,
		MinesRemain: opts.Mines,
	}, nil
}

// The CreateEmptyBoard create empty matrix of Cell.
func createBoard(opts BoardConfig) (Board, error) {
	logger.Print(log.Info("creating new board"))
	// Creating empty board
	board := make(Board, opts.Height)
	for i := range board {
		board[i] = make([]Cell, opts.Width)
	}

	if opts.Mines != 0 {
		err := board.placeMines(opts.Mines, opts.Seed)
		if err != nil {
			return board, err
		}
		err = board.placeNumbers()
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
func (b Board) To1D() ([]Cell, error) {
	logger.Print(log.Info("Converting 2D board to 1D"))
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

func (b Board) RevealAll() error {
	if len(b) == 0 {
		return errors.New("board is empty")
	}
	for i := range b {
		for j := range b[i] {
			b[i][j].Revealed = true
		}
	}
	return nil
}

func (b Board) placeNumbers() error {
	if len(b) == 0 {
		return errors.New("board can't be empty")
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
	return nil
}
