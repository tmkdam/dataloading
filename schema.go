package dataloading

import (
	"regexp"
	"io"
	"errors"
)

type fieldType struct {
	Name string
	Expression *regexp.Regexp
}

var types = []fieldType{
	fieldType{
		"datetime",
		regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\d)$`),
	},
	fieldType{
		"bigint",
		regexp.MustCompile(`^\-?[1-9]\d{9,17}$`),
	},
	fieldType{
		"int",
		regexp.MustCompile(`^\-?[1-9]\d{0,8}$`),
	},
	fieldType{
		"decimal",
		regexp.MustCompile(`^\-?\d*\.\d+$`),
	},
	fieldType{
		"boolean",
		regexp.MustCompile(`^(true|false|TRUE|FALSE|True|False)$`),
	},
}

type FieldInfo struct {
	FieldName string
	FieldType string
	MaxLength int
}

type SourceSchema struct {
	SourceName string
	Fields []FieldInfo
}

func Schema (desc Description) ([]SourceSchema, error) {
	sources, err := desc.Available()
	if err != nil {
		return nil, err
	}

	sourcesByName := make(map[string][]Source)

	for _, source := range sources {
		if sourcesByName[source.Name] == nil {
			sourcesByName[source.Name] = []Source{}
		}

		sourcesByName[source.Name] = append(sourcesByName[source.Name], source)
	}

	sourceSchemas := make([]SourceSchema, len(sourcesByName))

	sourceIndex := 0
	for sourceName, sources := range sourcesByName {
		fieldNames, err := desc.FieldNames(sourceName)
		if err != nil {
			return nil, err
		}

		discoveredTypeIndexes := make([]int, len(fieldNames))
		discoveredMaxLengths := make([]int, len(fieldNames))

		for _, source := range sources {
			reader, err := desc.Reader(source)
			if err != nil {
				return nil, err
			}


			for {
	  		row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}

				for fieldIndex, fieldName := range fieldNames {
					value, ok := row.Value(fieldName)
					if !ok {
						return nil, errors.New("Field '" + fieldName + "' not found in source")
					}
					valueLength := len(value)

					if discoveredMaxLengths[fieldIndex] < valueLength {
						discoveredMaxLengths[fieldIndex] = valueLength
					}

					if discoveredTypeIndexes[fieldIndex] == len(types) {
						continue
					}

					for {
						if discoveredTypeIndexes[fieldIndex] == len(types) {
							break
						}
						match := types[discoveredTypeIndexes[fieldIndex]].Expression.MatchString(value)
						if match == true || value == "" {
							break
						}
						discoveredTypeIndexes[fieldIndex] += 1
					}
				}
			}
		}

		sourceSchemas[sourceIndex] = SourceSchema{
			sourceName,
			make([]FieldInfo, len(fieldNames)),
		}

		for fieldIndex, fieldName := range fieldNames {
			var fieldType string
			if discoveredTypeIndexes[fieldIndex] == len(types) {
				fieldType = "string"
			} else {
				fieldType = types[discoveredTypeIndexes[fieldIndex]].Name
			}

			sourceSchemas[sourceIndex].Fields[fieldIndex] = FieldInfo{
				fieldName,
				fieldType,
				discoveredMaxLengths[fieldIndex],
			}
		}

		sourceIndex += 1
	}

	return sourceSchemas, nil
}