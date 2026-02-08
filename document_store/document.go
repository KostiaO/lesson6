package document_store

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNilInput = errors.New("nil input")
var ErrNonStructType = errors.New("non struct type input")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")
var ErrTypicalErrorBuilderMarshaling = errors.New("error in marshaling Document: %v")

var ErrDocumentNilProvided = errors.New("nil document provided")
var ErrNilOutput = errors.New("nil output provided")
var ErrOutputAsValueProvided = errors.New("output as value provided")
var ErrOutputNonStructType = errors.New("non struct type output")
var ErrTypesOfDocAndOutputNotMatch = errors.New("types of document and type doesnot match")
var ErrTypicalErrorBuilderUnmarshaling = errors.New("error in unmarshaling Document: %v")

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

func MarshalDocument(input any) (*Document, error) {

	if input == nil {
		return nil, fmt.Errorf(ErrTypicalErrorBuilderMarshaling.Error(), ErrNilInput)
	}

	var initalTypeInput reflect.Type = reflect.TypeOf(input)

	var initalValueInput reflect.Value = reflect.ValueOf(input)

	if initalTypeInput.Kind() == reflect.Ptr {
		initalTypeInput = initalTypeInput.Elem()
		initalValueInput = initalValueInput.Elem()
	}

	var typeOfInput reflect.Type = initalTypeInput

	var valueOfInput reflect.Value = initalValueInput

	var doc *Document

	if typeOfInput.Kind() != reflect.Struct {
		return nil, fmt.Errorf(ErrTypicalErrorBuilderMarshaling.Error(), ErrNonStructType)
	}

	if fieldsCount := typeOfInput.NumField(); fieldsCount > 0 {
		doc = &Document{
			Fields: map[string]DocumentField{},
		}

		for i := range fieldsCount {
			fieldInput := typeOfInput.Field(i)

			fieldValue := valueOfInput.FieldByName(fieldInput.Name)

			switch fieldInput.Type.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				doc.Fields[fieldInput.Name] = DocumentField{
					Type:  DocumentFieldTypeNumber,
					Value: fieldValue.Int(),
				}
			case reflect.String:
				doc.Fields[fieldInput.Name] = DocumentField{
					Type:  DocumentFieldTypeString,
					Value: fieldValue.String(),
				}
			case reflect.Bool:
				doc.Fields[fieldInput.Name] = DocumentField{
					Type:  DocumentFieldTypeBool,
					Value: fieldValue.Bool(),
				}
			case reflect.Array, reflect.Slice:

				doc.Fields[fieldInput.Name] = DocumentField{
					Type:  DocumentFieldTypeArray,
					Value: fieldValue.Interface(),
				}
			case reflect.Struct:
				doc.Fields[fieldInput.Name] = DocumentField{
					Type:  DocumentFieldTypeObject,
					Value: fieldValue,
				}
			default:
				return nil, fmt.Errorf(ErrTypicalErrorBuilderMarshaling.Error(), ErrUnsupportedDocumentField)
			}
		}
	}

	return doc, nil
}

func UnmarshalDocument(doc *Document, output any) error {
	if doc == nil {
		return fmt.Errorf(ErrTypicalErrorBuilderUnmarshaling.Error(), ErrDocumentNilProvided)
	}

	if output == nil {
		return fmt.Errorf(ErrTypicalErrorBuilderUnmarshaling.Error(), ErrNilOutput)
	}

	outputType := reflect.TypeOf(output)

	if outputType.Kind() != reflect.Ptr {
		return fmt.Errorf(ErrTypicalErrorBuilderUnmarshaling.Error(), ErrOutputAsValueProvided)
	}

	if outputType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf(ErrTypicalErrorBuilderUnmarshaling.Error(), ErrOutputNonStructType)
	}

	outputValue := reflect.ValueOf(output).Elem()

	for i := range outputType.Elem().NumField() {
		field := outputType.Elem().Field(i)

		docField, isPresent := doc.Fields[field.Name]

		if !isPresent {
			return fmt.Errorf(ErrTypicalErrorBuilderUnmarshaling.Error(), ErrUnsupportedDocumentField)
		}

		if outputValue.CanSet() && field.Type.Kind() == reflect.TypeOf(docField.Value).Kind() {
			switch reflect.TypeOf(docField.Value).Kind() {
			case reflect.String:
				outputValue.Field(i).SetString(reflect.ValueOf(docField.Value).String())

			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				outputValue.Field(i).SetInt(reflect.ValueOf(docField.Value).Int())

			case reflect.Bool:
				outputValue.Field(i).SetBool(reflect.ValueOf(docField.Value).Bool())
			case reflect.Array, reflect.Slice:
				outputValue.Field(i).Set(reflect.ValueOf(docField.Value))

			case reflect.Struct:
				outputValue.Field(i).Set(reflect.ValueOf(docField.Value))
			}
		} else {
			return fmt.Errorf(ErrTypicalErrorBuilderUnmarshaling.Error(), ErrTypesOfDocAndOutputNotMatch)
		}
	}

	return nil
}