package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <hash>\n")
			os.Exit(1)
		}

		hash := os.Args[3]
		if len(hash) != 40 {
			fmt.Fprintf(os.Stderr, "Invalid hash: %s\n", hash)
			os.Exit(1)
		}

		content, _ := readBlob(hash)

		fmt.Print(content)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

func readBlob(sha string) (string, error) {
	path := ".git/objects/" + sha[0:2] + "/" + sha[2:]

	commpressedBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	raw, err := zlib.NewReader(bytes.NewReader(commpressedBytes))
	if err != nil {
		return "", err
	}
	defer raw.Close()

	decompressed, err := ioutil.ReadAll(raw)
	if err != nil {
		return "", err
	}

	x := bytes.Index(decompressed, []byte{' '})
	y := bytes.Index(decompressed[x:], []byte{'\x00'}) + x

	return string(decompressed)[y+1:], nil
}
