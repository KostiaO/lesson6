package document_store

import (
	"errors"
	"log/slog"
)

var (
	ErrPrimaryKeysDoesNotMatch = errors.New("primary key of doc does not correspond to primary key of collection")
	ErrPrimaryKeyIsNotString   = errors.New("primary key of document is not string")
)

type CollectionConfig struct {
	PrimaryKey string
}

type Collection struct {
	CollectionConfig
	Data map[string]*Document
}

func (s *Collection) Put(doc Document) error {
	keyField, ok := doc.Fields[s.PrimaryKey]

	if !ok {
		return ErrPrimaryKeysDoesNotMatch
	}

	if keyField.Type != DocumentFieldTypeString {
		return ErrPrimaryKeyIsNotString
	}

	if newDocKey, isString := keyField.Value.(string); isString {
		s.Data[newDocKey] = &doc
	}

	slog.Default().Info("added document to collection:", slog.Any(keyField.Value.(string), doc))

	return nil
}

func (s *Collection) Get(key string) (*Document, bool) {
	document, ok := s.Data[key]

	return document, ok
}

func (s *Collection) Delete(key string) bool {
	doc, ok := s.Data[key]

	if !ok {
		return false
	}

	delete(s.Data, key)

	slog.Default().Info("deleted document in collection:", slog.Any(key, *doc))

	return true
}

func (s *Collection) List() []Document {
	listOfDocuments := make([]Document, 0, len(s.Data))

	for _, doc := range s.Data {
		listOfDocuments = append(listOfDocuments, *doc)
	}

	return listOfDocuments
}
