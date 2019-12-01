package tree

type prepared []byte

func reversePart(a []byte) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func prepareName(s string) prepared {
	out := make([]byte, 0, len(s))
	l := len(s)
	lp := 0
	if s[l-1] == '.' {
		l--
	}
	for i := l - 1; i >= 0; i-- {
		if s[i] == '\\' && i+3 < l && s[i+1] >= '0' && s[i+2] >= '0' && s[i+3] >= '0' && s[i+1] <= '9' && s[i+2] <= '9' && s[i+3] <= '9' {
			out = append(out[0:len(out)-3], dddToByte([]byte(s[i:])))
			continue
		} else if s[i] >= 'A' && s[i] <= 'Z' {
			out = append(out, s[i]-32)
			continue
		} else if s[i] == '.' {
			reversePart(out[lp:])
			lp = len(out) + 1
		}
		out = append(out, s[i])
	}
	reversePart(out[lp:])
	return out
}
