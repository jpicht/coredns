package tree

// less returns <0 when a is less than b, 0 when they are equal and
// >0 when a is larger than b.
// The function orders names in DNSSEC canonical order: RFC 4034s section-6.1
//
// See https://bert-hubert.blogspot.co.uk/2015/10/how-to-do-fast-canonical-ordering-of.html
// for a blog article on this implementation, although here we still go label by label.
//
// The values of a and b are *not* lowercased before the comparison!
func less(ap, bp prepared) int {
	i := 0
	al := len(ap)
	bl := len(bp)
	for i < al && i < bl {
		if ap[i] == '.' && bp[i] != '.' {
			return -1
		} else if ap[i] != '.' && bp[i] == '.' {
			return 1
		}
		res := int(ap[i]) - int(bp[i])
		if res != 0 {
			return res
		}

		i++
	}
	return al - bl
}

func dddToByte(s []byte) byte { return (s[1]-'0')*100 + (s[2]-'0')*10 + (s[3] - '0') }
