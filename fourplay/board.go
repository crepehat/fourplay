package fourplay

import (
	"errors"
	"fmt"
	"sync"
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
	CurrentBoard uint64
	Position     uint64
	Moves        int8
}

// BoardHeight represents the height of the board
var BoardHeight int8
var uBoardHeight uint8

// BoardWidth represents the width of the board
var BoardWidth int8
var uBoardWidth uint8

// BottomRowMask is a one in each of the last rows,
// calculated according to size of board
var BottomRowMask uint64

func NewBoard() Board {
	return Board{
		CurrentBoard: 0,
		Position:     0,
		Moves:        0,
	}
}

func XYMask(x, y uint8) uint64 {
	return ((1 << y) << (x * (uBoardHeight + 1)))
}

func PrintGrid(grid uint64) {
	for y := uBoardHeight; y != 255; y-- {
		for x := uint8(0); x < uBoardWidth; x++ {
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
	// fmt.Println(strconv.FormatInt(b.Position, 2))
	// fmt.Println(strconv.FormatInt(b.CurrentBoard, 2))
	// PrintGrid(b.CurrentBoard)
	// PrintGrid(b.Position)
	bothplayersboard := b.CurrentBoard + b.Position + BottomRowMask
	PrintGrid(bothplayersboard)

}

func TopMask(column uint8) uint64 {
	return (1 << (uBoardHeight - 1)) << (column * (uBoardHeight + 1))
}

func BottomMask(column uint8) uint64 {
	return (1 << (column * (uBoardHeight + 1)))
}

func ColumnMask(column uint8) uint64 {
	return (((1 << (uBoardHeight + 1)) - 1) << ((uBoardHeight + 1) * column))
}

// ValidMove is valid
func (b Board) ValidMove(column uint8) bool {
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
func (b Board) MakeMove(column uint8) (newboard Board) {
	newboard.Position = b.CurrentBoard ^ b.Position
	newboard.CurrentBoard = b.CurrentBoard | (b.CurrentBoard + BottomMask(column))
	newboard.Moves = b.Moves + 1
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
		if column >= uBoardWidth {
			return board, ErrInvalidColumn
		}
		if !board.ValidMove(column) {
			return board, ErrColumnFull
		}
		if board.IsWinningMove(column) {
			return board, ErrGameWon
		}
		board = board.MakeMove(column)
	}

	return
}

func (b Board) IsWinningMove(column uint8) bool {
	newposition := b.Position | ((b.CurrentBoard + BottomMask(column)) & ColumnMask(column))
	return IsWinning(newposition)
}

func IsWinning(position uint64) bool {
	// playerPosition := b.CurrentBoard ^ b.Position

	// horizontal
	m := position & (position >> (uBoardHeight + 1))
	if m&(m>>(2*(uBoardHeight+1))) != 0 {
		return true
	}
	// diag1
	m = position & (position >> (uBoardHeight))
	if m&(m>>(2*(uBoardHeight))) != 0 {
		return true
	}
	// diag2
	m = position & (position >> (uBoardHeight + 2))
	if m&(m>>(2*(uBoardHeight+2))) != 0 {
		return true
	}
	// vertical
	m = position & (position >> 1)
	if m&(m>>(2)) != 0 {
		return true
	}
	return false
}

func (b Board) GogaMax(wg *sync.WaitGroup, alphachan, betachan, maxchan chan int8) {
	defer wg.Done()

	if b.Moves == BoardHeight*BoardWidth {
		maxchan <- int8(0)
		return
	}

	// check if we can win next move
	for i := uint8(0); i < uBoardWidth; i++ {
		if b.ValidMove(i) && b.IsWinningMove(i) {
			// b.Print()
			// fmt.Printf("Can win next move. Score: %d\n", (BoardHeight*BoardWidth+1-b.Moves)/2)
			maxchan <- (BoardHeight*BoardWidth + 1 - b.Moves) / 2
		}
	}

	// upper bound of score since cannot win next move
	// max := (BoardHeight*BoardWidth - 1 - b.Moves) / 2

	beta := <-betachan
	alpha := <-alphachan
	max := <-maxchan
	if beta > max {
		beta = max
		if alpha >= beta {
			return beta
		}
	}

	var score int8
	for i := uint8(0); i < uBoardWidth; i++ {
		if b.ValidMove(i) {
			nextBoard := b.MakeMove(i)
			// nextBoard.Print()
			score = -nextBoard.NegaMax(-beta, -alpha)
			if score >= beta {
				return score
			}
			if score > alpha {
				alpha = score
			}
			// fmt.Println(bestScore)
		}
	}

	return alpha
}

func (b Board) NegaMax(alpha, beta int8) int8 {
	if b.Moves == BoardHeight*BoardWidth {
		return int8(0)
	}

	// check if we can win next move
	for i := uint8(0); i < uBoardWidth; i++ {
		if b.ValidMove(i) && b.IsWinningMove(i) {
			// b.Print()
			// fmt.Printf("Can win next move. Score: %d\n", (BoardHeight*BoardWidth+1-b.Moves)/2)
			return (BoardHeight*BoardWidth + 1 - b.Moves) / 2
		}
	}

	// upper bound of score since cannot win next move
	max := (BoardHeight*BoardWidth - 1 - b.Moves) / 2
	if beta > max {
		beta = max
		if alpha >= beta {
			return beta
		}
	}

	var score int8
	for i := uint8(0); i < uBoardWidth; i++ {
		if b.ValidMove(i) {
			nextBoard := b.MakeMove(i)
			// nextBoard.Print()
			score = -nextBoard.NegaMax(-beta, -alpha)
			if score >= beta {
				return score
			}
			if score > alpha {
				alpha = score
			}
			// fmt.Println(bestScore)
		}
	}

	return alpha
}
