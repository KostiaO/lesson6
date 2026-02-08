package document_store

import "errors"

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

	return nil
}

func (s *Collection) Get(key string) (*Document, bool) {
	document, ok := s.Data[key]

	return document, ok
}

func (s *Collection) Delete(key string) bool {
	_, ok := s.Data[key]

	if !ok {
		return false
	}

	delete(s.Data, key)

	return true
}

func (s *Collection) List() []Document {
	listOfDocuments := make([]Document, 0, len(s.Data))

	for _, doc := range s.Data {
		listOfDocuments = append(listOfDocuments, *doc)
	}

	return listOfDocuments
}
