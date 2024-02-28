package formats

func ByteToMegaByte(bytes int) float64 {
	// 1 Megabyte = 1024 * 1024 bytes
	return float64(bytes) / (1024 * 1024)
}

func NumberToMegaBytes(num float64) int {
	return int(num * 1024 * 1024)
}
