// dictionary_test.go
package dictionary_test

import (
	"estiam/dictionary"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWord(t *testing.T) {
	// 1. Create a new instance of the Dictionary.
	filename := "test_dictionary.txt"
	d := dictionary.New(filename)
	defer cleanupDictionaryFile(filename)

	// 2. Call the Add function to add a word to the dictionary.
	word := "testWord"
	definition := "testDefinition"
	message, err := d.Add(word, definition)

	// 3. Use assertions to verify that the word was added successfully.
	assert.NoError(t, err, "Unexpected error adding word")
	assert.Equal(t, fmt.Sprintf("Word '%s' Added successfully", word), message, "Unexpected message")

	// Check that the word is actually in the dictionary
	entry, err := d.Get(word)
	assert.NoError(t, err, "Unexpected error getting word")
	assert.Equal(t, definition, entry.Definition, "Unexpected definition for the added word")
}

// Helper function to clean up the test dictionary file
func cleanupDictionaryFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Printf("Error removing test dictionary file: %v\n", err)
	}
}

func TestGetWord(t *testing.T) {
	// 1. Create a new instance of the Dictionary.
	filename := "test_dictionary.txt"
	d := dictionary.New(filename)
	defer cleanupDictionaryFile(filename)

	// 2. Add a word to the dictionary.
	word := "testWord"
	definition := "testDefinition"
	_, err := d.Add(word, definition)
	assert.NoError(t, err, "Unexpected error adding word")

	// 3. Call the Get function to retrieve the added word.
	entry, err := d.Get(word)

	// 4. Use assertions to verify that the word was retrieved successfully.
	assert.NoError(t, err, "Unexpected error getting word")
	assert.Equal(t, definition, entry.Definition, "Unexpected definition for the retrieved word")
}

func TestRemoveWord(t *testing.T) {
	// 1. Create a new instance of the Dictionary.
	filename := "test_dictionary.txt"
	d := dictionary.New(filename)
	defer cleanupDictionaryFile(filename)

	// 2. Add a word to the dictionary.
	word := "testWord"
	definition := "testDefinition"
	_, err := d.Add(word, definition)
	assert.NoError(t, err, "Unexpected error adding word")

	// 3. Call the Remove function to remove the added word.
	message, err := d.Remove(word)

	// 4. Use assertions to verify that the word was removed successfully.
	assert.NoError(t, err, "Unexpected error removing word")
	assert.Equal(t, fmt.Sprintf("Word '%s' removed successfully", word), message, "Unexpected message")

	// Check that the word is no longer in the dictionary
	_, err = d.Get(word)
	assert.Error(t, err, "Expected error getting removed word")
}

func TestListWords(t *testing.T) {
	// 1. Create a new instance of the Dictionary.
	filename := "test_dictionary.txt"
	d := dictionary.New(filename)
	defer cleanupDictionaryFile(filename)

	// 2. Add multiple words to the dictionary.
	wordsToAdd := []struct {
		word       string
		definition string
	}{
		{"word1", "definition1"},
		{"word2", "definition2"},
		{"word3", "definition3"},
	}

	for _, entry := range wordsToAdd {
		_, err := d.Add(entry.word, entry.definition)
		assert.NoError(t, err, "Unexpected error adding word")
	}

	// 3. Call the List function to retrieve the list of words.
	words, _ := d.List()

	// 4. Use assertions to verify the list of words.
	expectedWords := []string{"word1", "word2", "word3"}
	assert.ElementsMatch(t, expectedWords, words, "Unexpected list of words")
}
