package fourplay

func init() {
	BoardHeight = 6
	uBoardHeight = 6
	BoardWidth = 7
	uBoardWidth = 7
	for i := uint8(0); i < uBoardWidth; i++ {
		BottomRowMask = BottomRowMask | BottomMask(i)
	}
}
