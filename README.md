# godiff

godiff is a small library for creating unified style diffs.

[![GoDoc](https://godoc.org/github.com/karrick/godiff?status.svg)](https://godoc.org/github.com/karrick/godiff)

## Description

This library generates a line by line comparison of its two input
slices of strings.

## Example

```Go
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
	
	
	for _, s := range godiff.Strings(left, right) {
        fmt.Println(s)
	}
```

MIT License

Copyright (c) 2019 Karrick McDermott

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
