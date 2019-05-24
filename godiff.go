package godiff // import "github.com/karrick/godiff"
import (
	"fmt"
	"os"
)

func reverse(list []string) {
	// Walk i forward from first index and l backwards from final index

	l := len(list)
	if l < 2 {
		return // done when 0 or 1 item in list
	}
	l--

	var i int

	for i < l {
		t := list[i]
		list[i] = list[l]
		list[l] = t
		i++
		l--
	}
}

// partition splits the two lists into logical chunks of a shared prefix, a
// shared suffix, and any first over middle bits for the first and second side.
func partition(alpha, bravo []string) (prefix []string, suffix []string, uniqueAlpha []string, uniqueBravo []string) {
	debug("PARTITION: alpha: %v\n", alpha)
	debug("PARTITION: bravo: %v\n", bravo)

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
	for {
		debug("PARTITION: prefix ai: %v (%s); bi: %v (%s)\n", ai, alpha[ai], bi, bravo[bi])
		if alpha[ai] != bravo[bi] {
			debug("PARTITION: no prefix match: %v; %v\n", alpha[ai], bravo[bi])
			break
		}
		prefix = append(prefix, alpha[ai])
		ai++
		bi++
		if ai == ae || bi == be {
			debug("PARTITION: prefix reached end\n")
			// either one or both of these loops will be no-ops
			for i := bi; i < be; i++ {
				debug("PARTITION: adding unique bravo: %s\n", bravo[i])
				uniqueBravo = append(uniqueBravo, bravo[i])
			}
			for i := ai; i < ae; i++ {
				debug("PARTITION: adding unique alpha: %s\n", alpha[i])
				uniqueAlpha = append(uniqueAlpha, alpha[i])
			}
			return
		}
	}

	// POST: alpha and bravo strings do not match; but we still have more in both lists
	debug("PARTITION: ai: %d; alpha: %v\n", ai, alpha[ai:])
	debug("PARTITION: bi: %d; bravo: %v\n", bi, bravo[bi:])

	// build common suffix
	defer func() { reverse(suffix) }() // reverse string slice in place

	as, bs := ai, bi // remember right side of slices where each started from
	ae, be = ai, bi  //
	debug("PARTITION: alpha end: %v\n", ae)
	debug("PARTITION: bravo end: %v\n", be)
	ai = al - 1 // start alpha index at final element in alpha slice
	bi = bl - 1 // start bravo index at final element in bravo slice
	debug("PARTITION: moving alpha reverse from [%d to %d]: %v\n", ai, ae, alpha[ae:ai+1])
	debug("PARTITION: moving bravo reverse from [%d to %d]: %v\n", bi, be, bravo[be:bi+1])

	for {
		debug("PARTITION: suffix ai: %v (%s); bi: %v (%s)\n", ai, alpha[ai], bi, bravo[bi])
		if alpha[ai] != bravo[bi] {
			debug("PARTITION: no suffix match: %v; %v\n", alpha[ai], bravo[bi])
			break
		}
		suffix = append(suffix, alpha[ai])
		debug("PARTITION: added item to suffix: %v\n", suffix)
		ai--
		bi--
		if ai < ae || bi < be {
			debug("PARTITION: suffix reached end\n")
			if ai < 0 {
				debug("PARTITION: as: %d; ai: %d: HERE\n", as, ai)
			} else {
				debug("PARTITION: as: %d; ai: %d; alpha: %v\n", as, ai, alpha[ai:])
			}
			if bi < 0 {
				debug("PARTITION: bs: %d; bi: %d: HERE\n", bs, bi)
			} else {
				debug("PARTITION: bs: %d; bi: %d; bravo: %v\n", bs, bi, bravo[bi:])
			}
			if ai < ae {
				debug("PARTITION: reached end of alpha: be: %d; bs: %d\n", be, bs)
				for i := be; i < bs+2; i++ {
					debug("PARTITION: adding unique bravo: %s\n", bravo[i])
					uniqueBravo = append(uniqueBravo, bravo[i])
				}
			}
			if bi < be {
				debug("PARTITION: reached end of bravo: ae: %d; as: %d\n", ae, as)
				for i := ae; i < as+2; i++ {
					debug("PARTITION: adding unique alpha: %s\n", alpha[i])
					uniqueAlpha = append(uniqueAlpha, alpha[i])
				}
			}
			return
		}
	}

	debug("PARTITION: ai: %d; alpha: %v\n", ai, alpha[ai:])
	debug("PARTITION: bi: %d; bravo: %v\n", bi, bravo[bi:])

	for i := ae; i < ai+1; i++ {
		debug("PARTITION: adding unique alpha: %s\n", alpha[i])
		uniqueAlpha = append(uniqueAlpha, alpha[i])
	}
	for i := be; i < bi+1; i++ {
		debug("PARTITION: adding unique bravo: %s\n", bravo[i])
		uniqueBravo = append(uniqueBravo, bravo[i])
	}

	return
}

// findNextMatch returns the index in the alpha and bravo lists where both list
// share a matching string that occurs only once in both lists.
func findNextMatch(alpha, bravo []string) (sa int, sb int) {
	debug("FIND NEXT MATCH: %v\n", alpha)
	debug("FIND NEXT MATCH: %v\n", bravo)

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
	// prefix -- print with spaces
	// L1 | R1 -- recurse
	// pivot -- print with space
	// L2 | R2 -- recurse
	// suffix -- print with spaces

	var i int
	debug("STRINGS %d: alpha: %v\n", i, alpha)
	debug("STRINGS %d: bravo: %v\n", i, bravo)
	prefix, suffix, uniqueAlpha, uniqueBravo := partition(alpha, bravo)
	debug("STRINGS %d: prefix: %v\n", i, prefix)
	debug("STRINGS %d: suffix: %v\n", i, suffix)
	debug("STRINGS %d: ualpha: %v\n", i, uniqueAlpha)
	debug("STRINGS %d: ubravo: %v\n", i, uniqueBravo)

	for _, s := range prefix {
		debug("STRINGS %d: appending prefix %s\n", i, s)
		diff = append(diff, " "+s)
	}

	// uniqueAlpha :=        dog, cat, elephant, hippo
	// uniqueBravo := tiger, dog, cat, hippo, monkey

	am, bm := findNextMatch(uniqueAlpha, uniqueBravo)
	// NOTE: Both values are -1, or both values are greater or equal to 0.

	if am >= 0 {
		debug("STRINGS %d: match found; am: %d; bm: %d; %v\n", i, am, bm, uniqueAlpha[am])

		rdiff := Strings(uniqueAlpha[:am], uniqueBravo[:bm])
		for _, s := range rdiff {
			debug("STRINGS %d: appending pre pivot rdiff: %s\n", i, s)
			diff = append(diff, s)
		}

		debug("STRINGS %d: appending pivot: %s\n", i, uniqueAlpha[am])
		diff = append(diff, " "+uniqueAlpha[am])

		rdiff = Strings(uniqueAlpha[am+1:], uniqueBravo[bm+1:])
		for _, s := range rdiff {
			debug("STRINGS %d: appending post pivot rdiff: %s\n", i, s)
			diff = append(diff, s)
		}
	} else {
		debug("STRINGS %d: no match found\n", i)
		for _, s := range uniqueAlpha {
			diff = append(diff, "-"+s)
		}
		for _, s := range uniqueBravo {
			diff = append(diff, "+"+s)
		}
	}

	for _, s := range suffix {
		debug("STRINGS %d: appending suffix %s\n", i, s)
		diff = append(diff, " "+s)
	}

	return
}

func debug(f string, v ...interface{}) {
	if false {
		fmt.Fprintf(os.Stderr, f, v...)
	}
}
