package main

import (
	"fmt"
	"lesson6/document_store"
	"log/slog"
	"os"
)

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(slogHandler)

	slog.SetDefault(logger)

	s := document_store.NewStore()

	isBool, collection := s.CreateCollection("col", &document_store.CollectionConfig{PrimaryKey: "prime"})

	if !isBool {
		return
	}

	documentOne := &document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"prime": {Type: document_store.DocumentFieldTypeString, Value: "DOC"},
			"info":  {Type: document_store.DocumentFieldTypeString, Value: "Info"},
		},
	}

	err := collection.Put(*documentOne)

	if err != nil {
		fmt.Println(err)
		return
	}
}
