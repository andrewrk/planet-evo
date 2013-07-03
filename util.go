package main

func sign(x float64) float64 {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}
