package main

import (
	"fmt"

	"github.com/qiv1ne/log"
)

func main() {
	board, err := NewMineBoard(BoardConfig{
		Mines:  8,
		Width:  3,
		Height: 3,
		Seed:   NewSeed(),
	})
	if err != nil {
		panic(err)
	}
	logger.Print(log.Info(fmt.Sprintf("%v", board.Real)))
	board.Real.Print()
	// board.User.Print()
}
