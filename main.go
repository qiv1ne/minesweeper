package main

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
	board.Real.Print()
	_, err = board.OpenCell(1, 1)
	board.User.Print()

}
