package main

func main() {
	board, err := NewMineBoard(BoardConfig{
		Mines:  10,
		Width:  10,
		Height: 10,
		Seed:   NewSeed(),
	})
	if err != nil {
		panic(err)
	}
	PrintBroadGracefully(board.RealBoard)
	PrintBroadGracefully(board.UserBoard)
	b, _ := RevealAll(board.UserBoard)
	PrintBroadGracefully(b)
}
