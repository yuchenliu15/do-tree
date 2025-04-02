package main

import (
	"os"
	"log"
	"strings"
)

type dir struct {
	name string
	indent int
	parent string
}

func tree(root string) ([]string, error) {
	stack := []dir{{root, 0, ""}} 
	result := []string{}
	for len(stack) > 0 {
		last := stack[len(stack)-1]
		result = append(result, strings.Repeat("\t", last.indent)+last.name)
		stack = stack[:len(stack)-1]
		full_path := last.parent + last.name
		entries, err := os.ReadDir(full_path)
		if err != nil {
			log.Printf("Error reading '%s': %s\n", last.name, err)
			log.Printf("Result: %s\n", result)
			return nil,  err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				stack = append(stack, dir{entry.Name(), last.indent+1, full_path+"/"})
			} else {
				result = append(result, strings.Repeat("\t", last.indent)+entry.Name())
			}
		}
	}
	return result, nil
}
