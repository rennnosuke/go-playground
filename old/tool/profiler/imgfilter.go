package profiler

func FlipFilter(img [][][3]uint8) [][][3]uint8 {
	newImg := make([][][3]uint8, len(img))
	for r := range img {
		newImg[r] = make([][3]uint8, len(img[r]))
		for c := range img[r] {
			newImg[r][c][0] = 255 - img[r][c][0]
			newImg[r][c][1] = 255 - img[r][c][1]
			newImg[r][c][2] = 255 - img[r][c][2]
		}
	}
	return newImg
}
