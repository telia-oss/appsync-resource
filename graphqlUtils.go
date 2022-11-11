package resource

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
)

func combineSchemas(oldSDL string, newSDL string, resolvers []Resolver) (string, error) {
	oldSchema, err := parser.Parse(parser.ParseParams{
		Source: string(oldSDL),
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})

	if err != nil {
		return "", err
	}

	newSchema, err := parser.Parse(parser.ParseParams{
		Source: string(newSDL),
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})

	if err != nil {
		return "", err
	}

	// Ensure that resolved fields are in the schema
	for _, resolver := range resolvers {
		existingType := getTypeFromSchema(oldSchema, resolver.TypeName)

		if existingType == nil {
			continue
		}

		existingObjectType := existingType.(*ast.ObjectDefinition)

		newType := getTypeFromSchema(newSchema, resolver.TypeName).(*ast.ObjectDefinition)
		newField := getFieldFromObjectDefinition(newType, resolver.FieldName)

		insertTypeAndSubtypesIntoSchema(oldSchema, newSchema, newField.Type)
		for _, arg := range newField.Arguments {
			insertTypeAndSubtypesIntoSchema(oldSchema, newSchema, arg.Type)
		}

		replaceFieldInObjectDefinition(existingObjectType, newField)
	}

	// Add all resolved types and their subtypes to the schema
	// for _, resolver := range resolvers {

	// }

	printed := printer.Print(oldSchema)

	return printed.(string), nil
}

func getTypeFromSchema(doc *ast.Document, typeName string) ast.Node {
	for _, def := range doc.Definitions {
		switch def := def.(type) {
		case *ast.ObjectDefinition:
			if def.Name.Value == typeName {
				return def
			}
		case *ast.InterfaceDefinition:
			if def.Name.Value == typeName {
				return def
			}
		case *ast.UnionDefinition:
			if def.Name.Value == typeName {
				return def
			}
		case *ast.ScalarDefinition:
			if def.Name.Value == typeName {
				return def
			}
		case *ast.EnumDefinition:
			if def.Name.Value == typeName {
				return def
			}
		case *ast.InputObjectDefinition:
			if def.Name.Value == typeName {
				return def
			}
		}
	}
	return nil
}

func replaceOrInsertTypeInSchema(doc *ast.Document, typeName string, typ ast.Node) {
	for i, def := range doc.Definitions {
		switch def := def.(type) {
		case *ast.ObjectDefinition:
			if def.Name.Value == typeName {
				doc.Definitions[i] = typ
				return
			}
		case *ast.InterfaceDefinition:
			if def.Name.Value == typeName {
				doc.Definitions[i] = typ
				return
			}
		case *ast.UnionDefinition:
			if def.Name.Value == typeName {
				doc.Definitions[i] = typ
				return
			}
		case *ast.ScalarDefinition:
			if def.Name.Value == typeName {
				doc.Definitions[i] = typ
				return
			}
		case *ast.EnumDefinition:
			if def.Name.Value == typeName {
				doc.Definitions[i] = typ
				return
			}
		case *ast.InputObjectDefinition:
			if def.Name.Value == typeName {
				doc.Definitions[i] = typ
				return
			}
		}
	}

	doc.Definitions = append(doc.Definitions, typ)
}

func getFieldFromObjectDefinition(obj *ast.ObjectDefinition, fieldName string) *ast.FieldDefinition {
	for _, field := range obj.Fields {
		if field.Name.Value == fieldName {
			return field
		}
	}
	return nil
}

func replaceFieldInObjectDefinition(obj *ast.ObjectDefinition, field *ast.FieldDefinition) {

	for i, existingField := range obj.Fields {
		if existingField.Name.Value == field.Name.Value {
			obj.Fields[i] = field
			return
		}
	}

	obj.Fields = append(obj.Fields, field)
}

func insertTypeAndSubtypesIntoSchema(target *ast.Document, source *ast.Document, typ ast.Type) {
	switch typ := typ.(type) {
	case *ast.Named:
		typeDefinition := getTypeFromSchema(source, typ.Name.Value)

		replaceOrInsertTypeInSchema(target, typ.Name.Value, typeDefinition)
		insertChildTypes(target, source, typeDefinition)

	case *ast.List:
		insertTypeAndSubtypesIntoSchema(target, source, typ.Type)
	case *ast.NonNull:
		insertTypeAndSubtypesIntoSchema(target, source, typ.Type)
	}
}

func insertChildTypes(target *ast.Document, source *ast.Document, typ ast.Node) {
	switch def := typ.(type) {
	case *ast.ObjectDefinition:
		for _, field := range def.Fields {
			insertTypeAndSubtypesIntoSchema(target, source, field.Type)
			for _, arg := range field.Arguments {
				insertTypeAndSubtypesIntoSchema(target, source, arg.Type)
			}
		}
	case *ast.InterfaceDefinition:
		for _, field := range def.Fields {
			insertTypeAndSubtypesIntoSchema(target, source, field.Type)
			for _, arg := range field.Arguments {
				insertTypeAndSubtypesIntoSchema(target, source, arg.Type)
			}
		}

		for _, obj := range source.Definitions {
			if obj, ok := obj.(*ast.ObjectDefinition); ok {
				for _, impl := range obj.Interfaces {
					if impl.Name.Value == def.Name.Value {
						replaceOrInsertTypeInSchema(target, obj.Name.Value, obj)
						insertChildTypes(target, source, obj)
					}
				}
			}
		}

	case *ast.UnionDefinition:
		for _, typ := range def.Types {
			insertTypeAndSubtypesIntoSchema(target, source, typ)
		}
	case *ast.InputObjectDefinition:
		for _, field := range def.Fields {
			insertTypeAndSubtypesIntoSchema(target, source, field.Type)
		}
	}
}
