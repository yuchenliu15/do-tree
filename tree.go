package dotree

import (
	"fmt"
	"os"
)

func tree(current_dir string) ([]string, bool) {

	entries, err := os.ReadDir(current_dir)
	if err != nil {
		fmt.Println("Error reading dir: %s", err)
		return nil, false
	}

}
