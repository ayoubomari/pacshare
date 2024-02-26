package formats

func ByteToMegaByte(bytes int) float64 {
	// 1 Megabyte = 1024 * 1024 bytes
	return float64(bytes) / (1024 * 1024)
}
