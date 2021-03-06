package main

import (
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

func removeDotfiles(de []os.DirEntry) (ret []os.DirEntry) {
	for _, e := range de {
		if e.Name()[0] != 46 {
			ret = append(ret, e)
		}
	}
	return
}

// contains checks if a string is present in
// a string slice
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func printTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Input", "Output", "Status"})
	table.SetAutoWrapText(false)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func filenameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func walk(paths map[string][]os.DirEntry, includeHidden bool) (map[string][]os.DirEntry, error) {
	iterated := []string{}
	var n = make(map[string][]os.DirEntry)

loop:
	for k, v := range paths {
		if contains(iterated, k) {
			continue
		}

		if !includeHidden {
			v = removeDotfiles(v)
		}

		iterated = append(iterated, k)
		for _, de := range v {
			if de.IsDir() {
				fp := filepath.Join(k, de.Name())
				dirEntry, err := os.ReadDir(fp)
				if err != nil {
					return nil, err
				}

				n[fp] = dirEntry
			}
		}
	}

	if len(n) > 0 {
		for k, v := range n {
			paths[k] = v
			delete(n, k)
		}

		goto loop
	}

	return paths, nil
}
