package request

type WordsLenType int

const (
	defaultBitSizeOfEntropy = 128
	interval                = 32
)

const (
	WordsLenTwelve    WordsLenType = 12
	WordsLenFifteen   WordsLenType = 15
	WordsLenEighteen  WordsLenType = 18
	WordsTwentyOne    WordsLenType = 21
	WordsLenTwentyTwo WordsLenType = 24
)

func (w *WordsLenType) GetBitSize() int {
	if w == nil {
		return 0
	}
	bitSize := defaultBitSizeOfEntropy
	var increment int
	switch *w {
	case WordsLenTwelve:
	case WordsLenFifteen:
		increment = interval
	case WordsLenEighteen:
		increment = interval * 2
	case WordsTwentyOne:
		increment = interval * 3
	case WordsLenTwentyTwo:
		increment = interval * 4
	default:
		//the number length of words should be 12, 15, 18, 21 or 24
		return 0
	}
	bitSize += increment

	return bitSize
}
