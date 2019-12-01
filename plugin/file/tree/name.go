package tree

type prepared []byte

func reversePart(a []byte) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

/*
	In Master Files [STD13] and other human-readable and -writable ASCII
	contexts, an escape is needed for the byte value for period (0x2E,
	".") and all octet values outside of the inclusive range from 0x21
	("!") to 0x7E ("~").  That is to say, 0x2E and all octet values in
	the two inclusive ranges from 0x00 to 0x20 and from 0x7F to 0xFF.
*/

var charmapping = [256]byte{}

func init() {
	for i := 1; i < 256; i++ {
		if byte(i) == '.' {
			charmapping[i] = 0
		} else if byte(i) >= 'A' && byte(i) <= 'Z' {
			charmapping[i] = byte(i) + 32
		} else {
			charmapping[i] = byte(i)
		}
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
		if s[i] == '\\' {
			if i+3 < l && s[i+1] >= '0' && s[i+2] >= '0' && s[i+3] >= '0' && s[i+1] <= '9' && s[i+2] <= '9' && s[i+3] <= '9' {
				out = append(out[0:len(out)-3], dddToByte([]byte(s[i:])))
			}
			continue
		} else if s[i] == '.' {
			reversePart(out[lp:])
			lp = len(out) + 1
		}

		out = append(out, charmapping[s[i]])
	}
	reversePart(out[lp:])
	return out
}
