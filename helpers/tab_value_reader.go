package helpers

import (
	"bufio"
	"io"
	"strings"
	"github.com/bloomapi/dataloading"
)

type TabReader struct {
	scanner  *bufio.Scanner
	fieldMap map[string]TabField
}

type TabRow struct {
	reader *TabReader
	record string
}

type TabField struct {
	Name string
	StartIndex uint64
	EndIndex uint64
}

func NewTabReader(r io.Reader, fields []TabField) *TabReader {
	fieldMap := map[string]TabField{}

	for _, elm := range fields {
		fieldMap[elm.Name] = elm
	}

	return &TabReader{
		scanner: bufio.NewScanner(r),
		fieldMap: fieldMap,
	}
}

func (r *TabReader) FieldNames() []string {
	fieldNames := make([]string, len(r.fieldMap))
	index := 0
	
	for _, value := range r.fieldMap {
		fieldNames[index] = value.Name
		index += 1
	}

	return fieldNames
}

func (r *TabReader) Read() (dataloading.Valuable, error) {
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return nil, err
		} else {
			return nil, io.EOF
		}
	}

	record := r.scanner.Text()
	
	return &TabRow{
		reader: r,
		record: record,
	}, nil
}

func (r *TabRow) Value(index string) (string, bool) {
	var (
		trueEnd int
	)

	field, ok := r.reader.fieldMap[index]
	if ok == false {
		return "", false
	}

	if field.StartIndex > uint64(len(r.record)) {
		return "", true
	}

	if len(r.record) == 0 {
		return "", true
	} else if field.EndIndex > uint64(len(r.record)) {
		trueEnd = len(r.record)
	} else {
		trueEnd = int(field.EndIndex)
	}

	value := r.record[(field.StartIndex - 1):trueEnd]
	value = strings.TrimSpace(value)
	return value, true
}
