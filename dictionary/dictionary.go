package dictionary

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Entry represents a dictionary entry containing a definition.
type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

// Dictionary represents a MongoDB-backed dictionary.
type Dictionary struct {
	collection *mongo.Collection
}

// EntryOperation represents a dictionary operation for adding or updating an entry.
type EntryOperation struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

// NewDictionary creates a new instance of the Dictionary.
func NewDictionary(databaseURI, databaseName, collectionName string) (*Dictionary, error) {
	clientOptions := options.Client().ApplyURI(databaseURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &Dictionary{
		collection: client.Database(databaseName).Collection(collectionName),
	}, nil
}

// Add adds a word with its definition to the dictionary.
func (d *Dictionary) Add(word string, definition string) (string, error) {
	entry := Entry{Definition: definition}

	_, err := d.collection.InsertOne(context.TODO(), map[string]interface{}{
		"word":       word,
		"definition": entry.Definition,
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Word '%s' Added successfully", word), nil
}

// Get retrieves the definition of a word from the dictionary.
func (d *Dictionary) Get(word string) (Entry, error) {
	var result map[string]interface{}
	err := d.collection.FindOne(context.TODO(), map[string]interface{}{
		"word": word,
	}).Decode(&result)

	if err != nil {
		return Entry{}, fmt.Errorf("word not found: %s", word)
	}

	definition, ok := result["definition"].(string)
	if !ok {
		return Entry{}, fmt.Errorf("invalid data structure for definition: %v", result["definition"])
	}

	return Entry{Definition: definition}, nil
}

// Remove removes a word and its definition from the dictionary.
func (d *Dictionary) Remove(word string) (string, error) {
	filter := map[string]interface{}{
		"word": word,
	}

	result, err := d.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error removing word: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Sprintf("Word '%s' not found", word), nil
	}

	return fmt.Sprintf("Word '%s' removed successfully", word), nil
}

// List retrieves a list of all words in the dictionary.
func (d *Dictionary) List() ([]string, error) {
	cursor, err := d.collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var words []string
	for cursor.Next(context.TODO()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		// Exclude _id field and extract word
		words = append(words, result["word"].(string))
	}

	return words, nil
}
