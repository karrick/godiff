package godiff

import (
	"testing"
)

func TestReverse(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var got []string
		reverse(got)
		ensureStringSlicesMatch(t, got, nil)
	})
	t.Run("single", func(t *testing.T) {
		got := []string{"hydrogen"}
		reverse(got)
		ensureStringSlicesMatch(t, got, []string{"hydrogen"})
	})
	t.Run("double", func(t *testing.T) {
		got := []string{"hydrogen", "helium"}
		reverse(got)
		ensureStringSlicesMatch(t, got, []string{"helium", "hydrogen"})
	})
}

func TestPartition(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Run("left", func(t *testing.T) {
			right := []string{"hydrogen", "helium", "lithium"}
			prefix, suffix, uniqueLeft, uniqueRight := partition(nil, right)

			ensureStringSlicesMatch(t, prefix, nil)
			ensureStringSlicesMatch(t, suffix, nil)
			ensureStringSlicesMatch(t, uniqueLeft, nil)
			ensureStringSlicesMatch(t, uniqueRight, []string{"hydrogen", "helium", "lithium"})
		})
		t.Run("right", func(t *testing.T) {
			left := []string{"hydrogen", "helium", "lithium"}
			prefix, suffix, uniqueLeft, uniqueRight := partition(left, nil)

			ensureStringSlicesMatch(t, prefix, nil)
			ensureStringSlicesMatch(t, suffix, nil)
			ensureStringSlicesMatch(t, uniqueLeft, []string{"hydrogen", "helium", "lithium"})
			ensureStringSlicesMatch(t, uniqueRight, nil)
		})
	})
	t.Run("short", func(t *testing.T) {
		t.Run("short left", func(t *testing.T) {
			left := []string{"hydrogen", "helium", "lithium"}
			right := []string{"hydrogen", "beryllium", "boron", "helium", "lithium"}

			prefix, suffix, uniqueLeft, uniqueRight := partition(left, right)

			ensureStringSlicesMatch(t, prefix, []string{"hydrogen"})
			ensureStringSlicesMatch(t, suffix, []string{"helium", "lithium"})
			ensureStringSlicesMatch(t, uniqueLeft, nil)
			ensureStringSlicesMatch(t, uniqueRight, []string{"beryllium", "boron"})
		})
		t.Run("right", func(t *testing.T) {
			left := []string{"hydrogen", "beryllium", "boron", "helium", "lithium"}
			right := []string{"hydrogen", "helium", "lithium"}

			prefix, suffix, uniqueLeft, uniqueRight := partition(left, right)

			ensureStringSlicesMatch(t, prefix, []string{"hydrogen"})
			ensureStringSlicesMatch(t, suffix, []string{"helium", "lithium"})
			ensureStringSlicesMatch(t, uniqueLeft, []string{"beryllium", "boron"})
			ensureStringSlicesMatch(t, uniqueRight, nil)
		})
	})
	t.Run("from strings", func(t *testing.T) {
		left := []string{"hydrogen", "helium", "hydrogen", "lithium", "carbon"}
		right := []string{"hydrogen", "boron", "helium", "hydrogen", "carbon", "nitrogen"}

		prefix, suffix, uniqueLeft, uniqueRight := partition(left, right)

		ensureStringSlicesMatch(t, prefix, []string{"hydrogen"})
		ensureStringSlicesMatch(t, suffix, nil)
		ensureStringSlicesMatch(t, uniqueLeft, []string{"helium", "hydrogen", "lithium", "carbon"})
		ensureStringSlicesMatch(t, uniqueRight, []string{"boron", "helium", "hydrogen", "carbon", "nitrogen"})
	})
}

func TestFindNextMatch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Run("both", func(t *testing.T) {
			var left, right []string

			lm, rm := findNextMatch(left, right)

			if got, want := lm, -1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := rm, -1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("left", func(t *testing.T) {
			var left []string
			right := []string{"oxygen", "lithium", "flourine", "helium"}

			lm, rm := findNextMatch(left, right)

			if got, want := lm, -1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := rm, -1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("right", func(t *testing.T) {
			left := []string{"oxygen", "lithium", "flourine", "helium"}
			var right []string

			lm, rm := findNextMatch(left, right)

			if got, want := lm, -1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := rm, -1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})
	t.Run("all distinct", func(t *testing.T) {
		left := []string{"hydrogen", "helium", "beryllium", "boron"}
		right := []string{"carbon", "nitrogen", "oxygen", "flourine"}

		lm, rm := findNextMatch(left, right)

		if got, want := lm, -1; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := rm, -1; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})
	t.Run("prefers alpha when both matches have same index", func(t *testing.T) {
		left := []string{"hydrogen", "beryllium", "helium", "lithium"}
		right := []string{"beryllium", "hydrogen", "helium", "lithium"}

		lm, rm := findNextMatch(left, right)

		if got, want := lm, 0; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := rm, 1; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})
	t.Run("prefers bravo when its index is smaller than alpha index", func(t *testing.T) {
		left := []string{"distinct1", "distinct2", "same1", "distinct3", "same2"}
		right := []string{"distinct4", "same1", "distinct5", "same2"}

		lm, rm := findNextMatch(left, right)

		if got, want := lm, 2; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := rm, 1; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})
}

func TestStrings(t *testing.T) {
	t.Run("given", func(t *testing.T) {
		left := []string{"hydrogen", "helium", "hydrogen", "lithium", "carbon"}
		right := []string{"hydrogen", "boron", "helium", "hydrogen", "carbon", "nitrogen"}

		//  hydrogen
		// +boron
		//  helium
		//  hydrogen
		// -lithium
		//  carbon
		// +nitrogen
		want := []string{" hydrogen", "+boron", " helium", " hydrogen", "-lithium", " carbon", "+nitrogen"}

		got := Strings(left, right)
		ensureStringSlicesMatch(t, got, want)
	})
}

func BenchmarkStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		left := []string{"hydrogen", "helium", "hydrogen", "lithium", "carbon"}
		right := []string{"hydrogen", "boron", "helium", "hydrogen", "carbon", "nitrogen"}

		//  hydrogen
		// +boron
		//  helium
		//  hydrogen
		// -lithium
		//  carbon
		// +nitrogen
		want := []string{" hydrogen", "+boron", " helium", " hydrogen", "-lithium", " carbon", "+nitrogen"}

		got := Strings(left, right)
		ensureStringSlicesMatch(b, got, want)
	}
}
