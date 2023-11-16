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
	addCh    chan entryOperation
	removeCh chan string
}

type entryOperation struct {
	word       string
	definition string
}

func New() *Dictionary {
	d := &Dictionary{
		entries:  make(map[string]Entry),
		addCh:    make(chan entryOperation),
		removeCh: make(chan string),
	}

	go d.processOperations()

	return d
}

func (d *Dictionary) processOperations() {
	for {
		select {
		case op := <-d.addCh:
			d.mu.Lock()
			d.entries[op.word] = Entry{Definition: op.definition}
			d.mu.Unlock()
		case word := <-d.removeCh:
			d.mu.Lock()
			delete(d.entries, word)
			d.mu.Unlock()
		}
	}
}

func (d *Dictionary) Add(word string, definition string) {
	d.addCh <- entryOperation{word: word, definition: definition}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	entry, found := d.entries[word]
	if !found {
		return Entry{}, nil
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	d.removeCh <- word
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
		_, err := fmt.Fprintf(writer, "%s: %s\n", word, entry.Definition)
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
