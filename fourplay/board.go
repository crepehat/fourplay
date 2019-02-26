package fourplay

import (
	"errors"
	"fmt"
)

/*
key = position + mask + bottom

board     position  mask      bottom    key
          0000000   0000000   0000000   0000000
.......   0000000   0000000   0000000   0001000
...o...   0000000   0001000   0000000   0010000
..xx...   0011000   0011000   0000000   0011000
..ox...   0001000   0011000   0000000   0001100
..oox..   0000100   0011100   0000000   0000110
..oxxo.   0001100   0011110   1111111   1101101
*/

// Board represents all the pieces currently on the board
type Board struct {
	CurrentBoard   uint64
	PlayerPosition uint64
	Moves          uint8
}

// BoardHeight represents the height of the board
var BoardHeight uint8

// BoardWidth represents the width of the board
var BoardWidth uint8

// BottomRowMask is a one in each of the last rows,
// calculated according to size of board
var BottomRowMask uint64

func NewBoard() Board {
	return Board{
		CurrentBoard:   0,
		PlayerPosition: 0,
		Moves:          0,
	}
}

func XYMask(x, y uint8) uint64 {
	return (1 << y << (x * (BoardHeight + 1)))
}

func PrintGrid(grid uint64) {
	for y := BoardHeight; y != 255; y-- {
		// fmt.Println(y)
		for x := uint8(0); x < BoardWidth; x++ {
			// fmt.Println(x)
			// fmt.Println(grid, XYMask(x,y))
			// fmt.Printf("%d%d ", x, y)
			if grid&XYMask(x, y) != 0 {
				fmt.Printf("x ")
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Print the board
func (b Board) Print() {
	// fmt.Println(strconv.FormatInt(b.PlayerPosition, 2))
	// fmt.Println(strconv.FormatInt(b.CurrentBoard, 2))
	// PrintGrid(b.CurrentBoard)
	// PrintGrid(b.PlayerPosition)
	bothplayersboard := b.CurrentBoard + b.PlayerPosition + BottomRowMask
	PrintGrid(bothplayersboard)

}

func TopMask(column uint8) uint64 {
	return (1 << (BoardHeight - 1)) << (column * (BoardHeight + 1))
}

func BottomMask(column uint8) uint64 {
	return (1 << (column * (BoardHeight + 1)))
}

// CheckMove is valid
func (b Board) CheckMove(column uint8) bool {
	return TopMask(column)&b.CurrentBoard == 0
}

var ErrRuneNotInt = errors.New("type: rune was not int")

func CharToUint8(r rune) (uint8, error) {
	if '0' <= r && r <= '9' {
		return uint8(r) - '0', nil
	}
	return uint8(0), ErrRuneNotInt
}

var ErrColumnFull = errors.New("move: not enough room in that column")
var ErrInvalidColumn = errors.New("move: not that many columns")
var ErrGameWon = errors.New("move: game won")

// MakeMove on a board.
func (b Board) MakeMove(column uint8) (newboard Board, err error) {
	if !b.CheckMove(column) {
		return b, ErrColumnFull
	}
	if column >= BoardWidth {
		return b, ErrInvalidColumn
	}
	newboard.PlayerPosition = b.CurrentBoard ^ b.PlayerPosition
	newboard.CurrentBoard = b.CurrentBoard | (b.CurrentBoard + BottomMask(column))
	newboard.Moves = b.Moves + 1
	newboard.Print()
	if newboard.Winning() {
		return b, ErrGameWon
	}
	return
}

// CreateFromSequence takes a board sequence and returns a board
// White goes first
func CreateFromSequence(sequence string) (board Board, err error) {
	board = NewBoard()
	for _, rune := range sequence {
		column, errchar := CharToUint8(rune)
		column--
		if errchar != nil {
			fmt.Println(err)
		}
		fmt.Println(board.Moves, column+1)
		board, err = board.MakeMove(column)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

// func

func (b Board) Winning() bool {
	playerPosition := b.CurrentBoard ^ b.PlayerPosition

	// horizontal
	m := playerPosition & (playerPosition >> (BoardHeight + 1))
	if m&(m>>(2*(BoardHeight+1))) != 0 {
		return true
	}
	// diag1
	m = playerPosition & (playerPosition >> (BoardHeight))
	if m&(m>>(2*(BoardHeight))) != 0 {
		return true
	}
	// diag2
	m = playerPosition & (playerPosition >> (BoardHeight + 2))
	if m&(m>>(2*(BoardHeight+2))) != 0 {
		return true
	}
	// vertical
	m = playerPosition & (playerPosition >> 1)
	if m&(m>>(2)) != 0 {
		return true
	}
	return false
}
