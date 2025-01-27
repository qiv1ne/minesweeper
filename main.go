package main

import (
	"fmt"

	"github.com/qiv1ne/log"
)

func main() {
	seed := 1
	board, err := NewMineBoard(BoardConfig{
		Mines:  8,
		Width:  4,
		Height: 4,
		Seed:   int64(seed),
	})
	if err != nil {
		panic(err)
	}
	logger.Print(log.Info(fmt.Sprintf("%v", board.Real)))
	board.Real.Print()
	result, err := board.OpenCell(1, 1)
	if err != nil {
		logger.Print(log.Error(err))
	}
	if result == Lose {
		logger.Print(log.Info("PROEBALI"))
	}
	board.User.Print()

}
