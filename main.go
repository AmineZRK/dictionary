package main

import (
	"bufio"
	"estiam/dictionary"
	"fmt"
	"os"
)

func main() {
	d := dictionary.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("1. Add Word")
		fmt.Println("2. Define Word")
		fmt.Println("3. Remove Word")
		fmt.Println("4. List Words")
		fmt.Println("5. Exit")

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
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please choose a number between 1 and 5.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word: ")
	word, _ := reader.ReadString('\n')
	word = word[:len(word)-1] // Remove the newline character

	fmt.Print("Enter definition: ")
	definition, _ := reader.ReadString('\n')
	definition = definition[:len(definition)-1] // Remove the newline character

	d.Add(word, definition)
	fmt.Println("Word added successfully!")
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to define: ")
	word, _ := reader.ReadString('\n')
	word = word[:len(word)-1] // Remove the newline character

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
	word = word[:len(word)-1] // Remove the newline character

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
