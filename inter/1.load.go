package inter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Load external definitions
func (i *Interpreter) Load(ior io.Reader) {

	old := i.scanner
	i.SetReader(ior)
	i.Run()
	i.scanner = old

}

// LoadFile loads the specified file name.
func (i *Interpreter) LoadFile(fileName string) {

	if i.Err != nil {
		return
	}

	f, err := os.Open(fileName)
	if err != nil {
		i.Err = err
		return
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	i.Load(bf)
	if i.Err != nil {
		fmt.Println("Error trying to load ", fileName, " : ", i.Err)
		return
	}

}

// initForth will load and compile the forth.forth file.
func (i *Interpreter) initForth() {

	var f *os.File
	fname := "forth.forth"
	names := []string{
		fname,
		filepath.Join(".", fname),
		filepath.Join("..", fname),
		filepath.Join("inter", fname),
		filepath.Join(".", "inter", fname),
		filepath.Join("..", "inter", fname),
	}
	for _, name := range names { // loop until we found the file
		// fmt.Println("DEBUG : Trying ", name)
		if f, i.Err = os.Open(name); i.Err == nil {
			// fmt.Println("DEBUG : Loading ", name)
			break
		}
	}

	if i.Err != nil {
		panic(i.Err)
	}

	defer f.Close()
	bf := bufio.NewReader(f)
	i.Load(bf)

	if i.Err != nil {
		panic(i.Err)
	}

}
