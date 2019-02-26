package fourplay

func init() {
	BoardHeight = 6
	BoardWidth = 7
	for i := uint8(0); i < BoardWidth; i++ {
		BottomRowMask = BottomRowMask | BottomMask(i)
	}
}
