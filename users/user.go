package users

import (
	"errors"
	"fmt"
	"lesson6/document_store"
)

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrUserCollectionNotFound       = errors.New("user collection is not present in provided store")
	ErrUserAlreadyExists            = errors.New("user already exists by this id: %v")
	ErrNewCollectionInStoreCreating = errors.New("error in creating new collection")
	ErrGetListOfUsers               = errors.New("error in getting list of users: %v")
	ErrInGetingUser                 = errors.New("error in getting user: %v")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll *document_store.Collection
}

func UserService(store *document_store.Store, table string) (*Service, error) {
	collection, tableExist := store.GetCollection(table)

	if !tableExist {
		IsCreated, newCollection := store.CreateCollection(table, &document_store.CollectionConfig{
			PrimaryKey: "ID",
		})

		if !IsCreated {
			return nil, ErrNewCollectionInStoreCreating
		}

		return &Service{
			coll: newCollection,
		}, nil
	}

	return &Service{
		coll: collection,
	}, nil
}

func (s *Service) CreateUser(id, name string) (*User, error) {
	if _, exist := s.coll.Get(id); exist {
		return nil, fmt.Errorf(ErrUserAlreadyExists.Error(), id)
	}

	user := &User{
		ID:   id,
		Name: name,
	}

	userDoc, err := document_store.MarshalDocument(user)

	if err != nil {
		return nil, fmt.Errorf("error in creating user on marshaling user %v", err)
	}

	errCollPut := s.coll.Put(*userDoc)

	if errCollPut != nil {
		return nil, errCollPut
	}

	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	userList := make([]User, 0, len(s.coll.Data))

	user := &User{}

	for _, doc := range s.coll.Data {
		err := document_store.UnmarshalDocument(doc, user)

		if err != nil {
			return nil, fmt.Errorf(ErrGetListOfUsers.Error(), err)
		}

		userList = append(userList, *user)
	}

	return userList, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	userDoc, ok := s.coll.Data[userID]

	if !ok {
		return nil, fmt.Errorf(ErrInGetingUser.Error(), ErrUserNotFound)
	}

	user := &User{}

	err := document_store.UnmarshalDocument(userDoc, user)

	if err != nil {
		return nil, fmt.Errorf(ErrInGetingUser.Error(), err)
	}

	return user, nil
}

func (s *Service) DeleteUser(userID string) error {
	isDeleted := s.coll.Delete(userID)

	if !isDeleted {
		return fmt.Errorf("error in deleting user: %v", ErrUserNotFound)
	}

	return nil
}
