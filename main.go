package main

import (
	"bufio"
	"estiam/dictionary"
	"fmt"
	"os"
	"strings"
)

const filename = "dictionary.txt"

func main() {
	d := dictionary.New()
	reader := bufio.NewReader(os.Stdin)

	// Load data from the file
	err := d.LoadFromFile(filename)
	if err != nil {
		fmt.Println("Error loading data:", err)
	}

	for {
		fmt.Println("1. Add Word")
		fmt.Println("2. Define Word")
		fmt.Println("3. Remove Word")
		fmt.Println("4. List Words")
		fmt.Println("5. Save and Exit")

		fmt.Print("Choose an action (1-5): ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			actionAdd(d, reader)
		case 2:
			actionDefine(d, reader)
		case 3:
			actionRemove(d, reader)
		case 4:
			actionList(d)
		case 5:
			// Save data to the file
			err := d.SaveToFile(filename)
			if err != nil {
				fmt.Println("Error saving data:", err)
			}
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please choose a number between 1 and 5.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)
	fmt.Println("Word added successfully!")
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to define: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println("Word not found.")
	} else {
		fmt.Println("Definition:", entry.String())
	}
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
	fmt.Println("Word removed successfully!")
}

func actionList(d *dictionary.Dictionary) {
	words, _ := d.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		entry, _ := d.Get(word)
		fmt.Printf("%s: %s\n", word, entry.Definition)
	}
}
