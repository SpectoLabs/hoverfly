package zeus

func trim(s, w string) string {
	ss := len(s)
	ws := len(w)

	if ss >= ws && s[:ws] == w {
		s = s[ws:]
		ss -= ws
	}

	if ss >= ws && s[ss-ws:] == w {
		s = s[:ss-ws]
	}

	return s
}

func split(s, d string) []string {
	var i, c, n int

	l := len(s)

	for i < l {
		if s[i:i+1] == d {
			c++
		}
		i++
	}

	a := make([]string, c+1)
	c, i = 0, 0

	for i < l {
		if i == l-1 {
			a[c] = s[n : i+1]
			break
		}
		if s[i:i+1] == d {
			a[c] = s[n:i]
			n = i + 1
			c++
		}
		i++
	}

	return a
}
