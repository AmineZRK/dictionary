package dictionary

import "fmt"

type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	entries map[string]Entry
}

func New() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Entry),
	}
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Definition: definition}
	d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, found := d.entries[word]
	if !found {
		return Entry{}, nil
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	words := make([]string, 0, len(d.entries))
	for word, entry := range d.entries {
		words = append(words, word)
		fmt.Println("Word:", word, "Definition:", entry.String()) // For debugging
	}
	return words, d.entries
}
