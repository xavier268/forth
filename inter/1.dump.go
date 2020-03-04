package inter

import (
	"fmt"
	"sort"
)

// dump
func (i *Interpreter) dump() {

	fmt.Printf("\nWords dumps, (size : %d)\n--NFA-----CFA--Flags-----Word---------\n", len(i.words))

	var keys sort.IntSlice
	for k := range i.words {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	for _, k := range keys {
		w := i.words[k]
		flags := []byte("---")
		if w.immediate {
			flags[0] = 'i'
		}
		if i.isPrimitive(w) {
			flags[2] = 'P'
		}
		fmt.Printf("%4d\t%4d\t%s\t%s\n", w.nfa, w.cfa, flags, w.name)
	}

	fmt.Printf("Memory dump, (size : %d) ", len(i.mem))
	for k, v := range i.mem {
		if k%5 == 0 {
			fmt.Printf("\n%5d --: ", k)
		}
		if k == i.ip {
			fmt.Printf("%s%5d%s ", ColorGreen, v, ColorOff)
		} else {
			fmt.Printf("%5d ", v)
		}
	}
	fmt.Println()

	fmt.Println("IP           : ", i.ip)
	fmt.Println("LastNFA      : ", i.lastNfa)
	fmt.Println("LastNFA prim : ", i.lastPrimitiveNfa)
	fmt.Println("Compile      : ", i.compileMode)
	fmt.Println("DS           : ", i.ds.data)
	fmt.Println("RS           : ", i.rs.data)
	fmt.Println()

}
