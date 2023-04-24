package bubble_bath

func GetMaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetMinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Clamp(value, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return GetMinInt(high, GetMaxInt(low, value))
}
