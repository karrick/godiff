package godiff // import "github.com/karrick/godiff"

func reverse(list []string) {
	// Walk i forward from first index and l backwards from final index

	l := len(list)
	if l < 2 {
		return // done when 0 or 1 item in list
	}
	l--

	var i int

	for i < l {
		list[i], list[l] = list[l], list[i] // swap elements
		i++
		l--
	}
}

// partition splits the two lists into logical chunks of a shared prefix, a
// shared suffix, and any middle bits for the first and second side.
func partition(alpha, bravo []string) (prefix []string, suffix []string, uniqueAlpha []string, uniqueBravo []string) {
	al := len(alpha)
	if al == 0 {
		uniqueBravo = bravo
		return
	}

	bl := len(bravo)
	if bl == 0 {
		uniqueAlpha = alpha
		return
	}

	var ai, bi int   // index = 0
	ae, be := al, bl // end   = length

	// build common prefix
	for alpha[ai] == bravo[bi] {
		prefix = append(prefix, alpha[ai])
		ai++
		bi++
		// When either slice is at its end, there is no common suffix.
		if ai == ae {
			for i := bi; i < be; i++ {
				uniqueBravo = append(uniqueBravo, bravo[i])
			}
			return
		}
		if bi == be {
			for i := ai; i < ae; i++ {
				uniqueAlpha = append(uniqueAlpha, alpha[i])
			}
			return
		}
	}
	// POST: alpha and bravo strings do not match; but we still have more in both lists

	// build common suffix
	defer func() { reverse(suffix) }() // in-place reversal of the suffix string slice

	as, bs := ai, bi // remember right side of slices where each started from
	ae, be = ai, bi  //
	ai = al - 1      // start alpha index at final element in alpha slice
	bi = bl - 1      // start bravo index at final element in bravo slice

	for alpha[ai] == bravo[bi] {
		suffix = append(suffix, alpha[ai])
		ai--
		bi--
		// When either slice is at its end, we are done.
		if ai < ae {
			for i := be; i < bs+2; i++ {
				uniqueBravo = append(uniqueBravo, bravo[i])
			}
			return
		}
		if bi < be {
			for i := ae; i < as+2; i++ {
				uniqueAlpha = append(uniqueAlpha, alpha[i])
			}
			return
		}
	}
	// POST: everything remaining is unique to each respective slice.

	for i := ae; i < ai+1; i++ {
		uniqueAlpha = append(uniqueAlpha, alpha[i])
	}
	for i := be; i < bi+1; i++ {
		uniqueBravo = append(uniqueBravo, bravo[i])
	}

	return
}

// findNextMatch returns the index in the alpha and bravo lists where both list
// share a matching string that occurs only once in both lists.
func findNextMatch(alpha, bravo []string) (sa int, sb int) {
	// used in this function to signal never found a distinct item
	sa = len(alpha)
	if sa == 0 {
		return -1, -1
	}
	sb = len(bravo)
	if sb == 0 {
		return -1, -1
	}

	// Walk bravo first, because lower alpha index is preferred, and we want to
	// use alpha as final map to walk, and we want to walk smaller map, thus
	// will build alpha with bravo data in mind.
	bc := make(map[string]int, len(bravo)) // value -> count
	bi := make(map[string]int, len(bravo)) // value -> index
	for i, s := range bravo {
		bc[s]++
		bi[s] = i
	}

	// Build alpha count histogram, but only include items where b count
	// histogram has recorded a single item for that same string value.
	ac := make(map[string]int, len(alpha)) // value -> count
	ai := make(map[string]int, len(alpha)) // value -> index
	for i, s := range alpha {
		if c, ok := bc[s]; ok && c == 1 {
			ac[s]++
			ai[s] = i
		}
	}

	for s, c := range ac {
		if c == 1 {
			// This item only shows up once in both alpha and bravo, but might
			// not be the smallest index
			if ia := ai[s]; sa > ia {
				// This item is the smallest index.
				sa = ia
				sb = bi[s]
			}
		}
	}

	if sa < len(alpha) && sb < len(bravo) {
		return
	}

	return -1, -1
}

// Strings returns the diff in unified format of two string slices.
func Strings(alpha, bravo []string) (diff []string) {
	prefix, suffix, uniqueAlpha, uniqueBravo := partition(alpha, bravo)

	for _, s := range prefix {
		diff = append(diff, " "+s)
	}

	// uniqueAlpha :=        dog, cat, elephant, hippo
	// uniqueBravo := tiger, dog, cat, hippo, monkey

	am, bm := findNextMatch(uniqueAlpha, uniqueBravo)
	// POST: Both values are -1, or both values are greater or equal to 0.

	if am >= 0 {
		rdiff := Strings(uniqueAlpha[:am], uniqueBravo[:bm])
		for _, s := range rdiff {
			diff = append(diff, s)
		}

		diff = append(diff, " "+uniqueAlpha[am])

		rdiff = Strings(uniqueAlpha[am+1:], uniqueBravo[bm+1:])
		for _, s := range rdiff {
			diff = append(diff, s)
		}
	} else {
		for _, s := range uniqueAlpha {
			diff = append(diff, "-"+s)
		}
		for _, s := range uniqueBravo {
			diff = append(diff, "+"+s)
		}
	}

	for _, s := range suffix {
		diff = append(diff, " "+s)
	}

	return
}
