package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	entries  map[string]Entry
	mu       sync.Mutex
	addCh    chan EntryOperation
	removeCh chan string
	filename string
}

type EntryOperation struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

func New(filename string) *Dictionary {
	d := &Dictionary{
		entries:  make(map[string]Entry),
		addCh:    make(chan EntryOperation),
		removeCh: make(chan string),
		filename: filename,
	}

	go d.processOperations()

	return d
}

func (d *Dictionary) processOperations() {
	for {
		select {
		case op := <-d.addCh:
			d.mu.Lock()
			d.entries[op.Word] = Entry{Definition: op.Definition}
			d.mu.Unlock()
		case word := <-d.removeCh:
			d.mu.Lock()
			delete(d.entries, word)
			d.mu.Unlock()
		}
	}
}

func (d *Dictionary) Add(word string, definition string) (string, error) {
	d.addCh <- EntryOperation{Word: word, Definition: definition}

	if err := d.SaveToFile(d.filename); err != nil {
		return "", fmt.Errorf("error saving data: %v", err)
	}

	return fmt.Sprintf("Word '%s' Added successfully", word), nil
}

func (d *Dictionary) Get(word string) (Entry, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	entry, found := d.entries[word]
	if !found {
		return Entry{}, fmt.Errorf("word not found: %s", word)
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, found := d.entries[word]; !found {
		return fmt.Sprintf("Word '%s' does not exist in the dictionary", word), nil
	}

	delete(d.entries, word)

	d.removeCh <- word

	if err := d.SaveToFile(d.filename); err != nil {
		return "", fmt.Errorf("error saving data: %v", err)
	}

	return fmt.Sprintf("Word '%s' removed successfully", word), nil

}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	words := make([]string, 0, len(d.entries))
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries
}

// SaveToFile saves the dictionary data to a file.
func (d *Dictionary) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for word, entry := range d.entries {
		_, err := fmt.Fprintln(writer, word+":", entry.Definition)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// LoadFromFile loads the dictionary data from a file.
func (d *Dictionary) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line format: %s", line)
		}

		word := strings.TrimSpace(parts[0])
		definition := strings.TrimSpace(parts[1])

		d.Add(word, definition)
	}

	return scanner.Err()
}
