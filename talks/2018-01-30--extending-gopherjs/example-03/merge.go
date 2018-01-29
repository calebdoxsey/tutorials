package main

// START OMIT

func merge(l, r []int) []int {
	m := make([]int, 0, len(l)+len(r))
	for len(l) > 0 || len(r) > 0 {
		switch {
		case len(l) == 0:
			m = append(m, r[0])
			r = r[1:]
		case len(r) == 0:
			m = append(m, l[0])
			l = l[1:]
		case l[0] <= r[0]:
			m = append(m, l[0])
			l = l[1:]
		case l[0] > r[0]:
			m = append(m, r[0])
			r = r[1:]
		}
	}
	return m
}

// END OMIT
