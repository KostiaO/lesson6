package document_store

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

	return false
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	// Функція повинна створити та проініціалізувати новий `Store`
	// зі всіма колекціями да даними з вхідного дампу.

}

func (s *Store) Dump() ([]byte, error) {
	// Методи повинен віддати дамп нашого стору в який включені дані про колекції та документ

	// TODO: Implement
}

// Значення яке повертає метод `store.Dump()` має без помилок оброблятись функцією `NewStoreFromDump`

func NewStoreFromFile(filename string) (*Store, error) {
	// Робить те ж саме що і функція `NewStoreFromDump`, але сам дамп має діставатись з файлу
	// TODO: Implement
}

func (s *Store) DumpToFile(filename string) error {
	// Робить те ж саме що і метод  `Dump`, але записує у файл замість того щоб повертати сам дамп

	// TODO: Implement
}
