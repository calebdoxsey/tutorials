package main

// START OMIT

// ConcurrentMergeSort sorts a list concurrently
func ConcurrentMergeSort(xs []int) []int {
	switch len(xs) {
	case 0:
		return nil
	case 1, 2:
		return merge(xs[:1], xs[1:])
	default:
		lc, rc := make(chan []int), make(chan []int)
		go func() {
			lc <- ConcurrentMergeSort(xs[:len(xs)/2])
		}()
		go func() {
			rc <- ConcurrentMergeSort(xs[len(xs)/2:])
		}()
		return merge(<-lc, <-rc)
	}
}

// END OMIT
