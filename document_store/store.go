package document_store

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
)

type Store struct {
	Storage map[string]*Collection
}

func NewStore() *Store {

	return &Store{
		Storage: make(map[string]*Collection),
	}

}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	if _, ok := s.Storage[name]; ok {
		return false, nil
	}
	newCollection := &Collection{
		Data: make(map[string]*Document),
	}
	newCollection.PrimaryKey = cfg.PrimaryKey

	s.Storage[name] = newCollection

	slog.Default().Info("created new collection:", name, "primaryKey:", cfg.PrimaryKey)

	return true, newCollection
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	if collection, ok := s.Storage[name]; ok {
		return collection, true
	}

	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.Storage[name]; ok {
		delete(s.Storage, name)

		return true
	}

	slog.Default().Info("deleted collection:", name)

	return false
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	var collections *map[string]*Collection

	err := json.Unmarshal(dump, collections)

	if err != nil {
		return nil, err
	}

	if collections == nil {
		return nil, err
	}

	store := NewStore()

	store.Storage = *collections

	return store, nil
}

func (s *Store) Dump() ([]byte, error) {
	dump, err := json.Marshal(s.Storage)

	if err != nil {
		return nil, err
	}

	return dump, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	var errs error = nil

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)

	if err != nil {
		return nil, errors.Join(errs, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}()

	reader := bufio.NewReader(file)

	fileBytes, errReading := reader.ReadBytes('\n')

	if errReading != nil && !errors.Is(errReading, io.EOF) {
		return nil, errors.Join(errs, errReading)
	}

	store, errStoringDump := NewStoreFromDump(fileBytes)

	if errStoringDump != nil {
		return nil, errors.Join(errs, errStoringDump)
	}

	return store, errs
}

func (s *Store) DumpToFile(filename string) error {
	var errs error = nil

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return errors.Join(errs, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}()

	dumpToWrite, errDumping := s.Dump()

	if errDumping != nil {
		return errors.Join(errs, errDumping)
	}

	_, errWriting := file.Write(dumpToWrite)

	if errWriting != nil {
		return errors.Join(errs, errWriting)
	}

	return errs
}
