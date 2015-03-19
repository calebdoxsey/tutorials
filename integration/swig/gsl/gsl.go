package gsl

func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	return Gsl_stats_mean(&xs[0], 1, int64(len(xs)))
}
