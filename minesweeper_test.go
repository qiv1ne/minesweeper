package main

import (
	"testing"
)

var (
	b = board(
		[][]Cell{
			{
				{Empty: false, MinesAround: 1, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 1, Flagged: false, IsMine: true, Revealed: true},
				{Empty: false, MinesAround: 3, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 2, Flagged: false, IsMine: false, Revealed: true},
			},
			{
				{Empty: false, MinesAround: 3, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 4, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 4, Flagged: false, IsMine: true, Revealed: true},
				{Empty: false, MinesAround: 2, Flagged: false, IsMine: true, Revealed: true},
			},
			{
				{Empty: false, MinesAround: 2, Flagged: false, IsMine: true, Revealed: true},
				{Empty: false, MinesAround: 3, Flagged: false, IsMine: true, Revealed: true},
				{Empty: false, MinesAround: 5, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 3, Flagged: false, IsMine: true, Revealed: true},
			},
			{
				{Empty: false, MinesAround: 2, Flagged: false, IsMine: true, Revealed: true},
				{Empty: false, MinesAround: 3, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 3, Flagged: false, IsMine: false, Revealed: true},
				{Empty: false, MinesAround: 1, Flagged: false, IsMine: true, Revealed: true},
			},
		})
)

func TestOpenCell(t *testing.T) {
	b, _ := NewMineBoard(BoardConfig{
		Width:  3,
		Height: 3,
		Mines:  4,
		Seed:   1,
	})
	b.OpenCell(1, 1)
	if b.User[0][0].Revealed == false {
		t.Fatalf("Want user board[0][0].Revealed = true, get: false")
	}
}
func TestOpenCellWithEmptyBoard(t *testing.T) {
	b := &MineBoard{}
	_, err := b.OpenCell(1, 1)
	if err == nil {
		t.Fatalf("Want error, get: %v", err)
	}
}

func TestOpenCellOutsideOfBoard(t *testing.T) {
	b, _ := NewMineBoard(BoardConfig{
		Width:  3,
		Height: 3,
		Mines:  4,
		Seed:   1,
	})
	_, err := b.OpenCell(100, 100)
	if err == nil {
		t.Fatalf("Want error, get: %v", err)
	}
}
func TestOpenCellOutsideOfBoardReversed(t *testing.T) {
	b, _ := NewMineBoard(BoardConfig{
		Width:  3,
		Height: 3,
		Mines:  4,
		Seed:   1,
	})
	_, err := b.OpenCell(-100, -100)
	if err == nil {
		t.Fatalf("Want error, get: %v", err)
	}
}

func TestPlaceFlag(t *testing.T) {
	b, _ := NewMineBoard(BoardConfig{
		Width:  3,
		Height: 3,
		Mines:  4,
		Seed:   1,
	})
	b.PlaceFlag(1, 1)
	if b.User[0][0].Flagged == false {
		t.Fatalf("Want user board[0][0].Flagged = true, get: %v", b.User[0][0].Flagged)
	}
}

func TestPlaceFlagWithEmptyBoard(t *testing.T) {
	b := &MineBoard{}
	_, err := b.PlaceFlag(1, 1)
	if err == nil {
		t.Fatalf("Want error, get: %v", err)
	}
}
func TestPlaceFlagOutsideOfBoard(t *testing.T) {
	b, _ := NewMineBoard(BoardConfig{
		Width:  3,
		Height: 3,
		Mines:  4,
		Seed:   1,
	})
	_, err := b.PlaceFlag(1, 1)
	if err == nil {
		t.Fatalf("Want error, get: %v", err)
	}
}
