package templating

import (
	"reflect"
	"sync"
	"testing"
)

func TestParseCommand(t *testing.T) {

	dataSources := map[string]*DataSource{
		"employees": {
			SourceType: "csv",
			Name:       "employees",
			Data: [][]string{
				{"id", "name", "age", "department"},
				{"1", "John Doe", "30", "Engineering"},
				{"2", "Jane Smith", "40", "Marketing"},
			},
			mu: sync.Mutex{},
		},
	}
	templateDataSource := NewTemplateDataSource()
	templateDataSource.dataSources = dataSources

	tests := []struct {
		query       string
		expected    SQLStatement
		expectError bool
	}{
		{
			query: "SELECT name, age FROM employees WHERE age >= '30'",
			expected: SQLStatement{
				Type:           "SELECT",
				Columns:        []string{"name", "age"},
				Conditions:     []Condition{{Column: "age", Operator: ">=", Value: "30"}},
				DataSourceName: "employees",
			},
			expectError: false,
		},
		{
			query:       "SELECT * FROM employees WHERE department == 'Engineering'",
			expected:    SQLStatement{Type: "SELECT", Columns: []string{"id", "name", "age", "department"}, Conditions: []Condition{{Column: "department", Operator: "==", Value: "Engineering"}}, DataSourceName: "employees"},
			expectError: false,
		},
		{
			query:       "UPDATE employees SET age = '35' WHERE name == 'John Doe'",
			expected:    SQLStatement{Type: "UPDATE", Conditions: []Condition{{Column: "name", Operator: "==", Value: "John Doe"}}, SetClauses: map[string]string{"age": "35"}, DataSourceName: "employees"},
			expectError: false,
		},
		{
			query:       "DELETE FROM employees WHERE age < '30'",
			expected:    SQLStatement{Type: "DELETE", Conditions: []Condition{{Column: "age", Operator: "<", Value: "30"}}, DataSourceName: "employees"},
			expectError: false,
		},
		{
			query:       "SELECT * FROM chillieplants WHERE department == 'Engineering'",
			expectError: true,
		},
		{
			query:       "INVALID QUERY",
			expectError: true,
		},
	}

	for _, test := range tests {
		result, err := parseSqlCommand(test.query, templateDataSource)
		if test.expectError {
			if err == nil {
				t.Errorf("expected error but got none for query: %s", test.query)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for query %s: %v", test.query, err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v for query: %s", test.expected, result, test.query)
			}
		}
	}
}

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`"12345`, `"12345`},
		{`hello`, "hello"},
		{`''`, ""},
	}

	for _, test := range tests {
		result := trimQuotes(test.input)
		if result != test.expected {
			t.Errorf("expected %v, got %v for input %s", test.expected, result, test.input)
		}
	}
}

func TestParseSetClauses_ValidInput(t *testing.T) {
	input := "age = '35', department = 'Engineering'"
	inputHeaders := []string{"age", "department"}
	expected := map[string]string{
		"age":        "35",
		"department": "Engineering",
	}

	result, _ := parseSetClauses(input, inputHeaders)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseSetClauses_InvalidInput(t *testing.T) {
	input := "age = '35', department = 'Engineering'"
	inputHeaders := []string{"fruit", "category"}

	_, err := parseSetClauses(input, inputHeaders)
	if err == nil {
		t.Errorf("expected error but got none.")
	}
}

func TestParseConditions_ValidInput(t *testing.T) {
	wherePart := "id == '1' AND name != 'John' AND age >= '30'"
	expected := []Condition{
		{"id", "==", "1"},
		{"name", "!=", "John"},
		{"age", ">=", "30"},
	}

	conditions, err := parseConditions(wherePart)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(conditions, expected) {
		t.Errorf("expected %v, got %v", expected, conditions)
	}
}

func TestParseConditions_EmptyInput(t *testing.T) {
	wherePart := ""
	expected := []Condition{}

	conditions, err := parseConditions(wherePart)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(conditions, expected) {
		t.Errorf("expected %v, got %v", expected, conditions)
	}
}

func TestParseConditions_SingleCondition(t *testing.T) {
	wherePart := "name == 'Alice'"
	expected := []Condition{
		{"name", "==", "Alice"},
	}

	conditions, err := parseConditions(wherePart)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(conditions, expected) {
		t.Errorf("expected %v, got %v", expected, conditions)
	}
}

func TestParseConditions_MultipleConditions(t *testing.T) {
	wherePart := "id == '1' AND age < '40'"
	expected := []Condition{
		{"id", "==", "1"},
		{"age", "<", "40"},
	}

	conditions, err := parseConditions(wherePart)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(conditions, expected) {
		t.Errorf("expected %v, got %v", expected, conditions)
	}
}

func TestExecuteSqlSelectQuery(t *testing.T) {
	data := [][]string{
		{"id", "name", "age", "department"},
		{"1", "John Doe", "30", "Engineering"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	query := SQLStatement{
		Type:           "SELECT",
		Columns:        []string{"name", "age"},
		Conditions:     []Condition{{Column: "department", Operator: "==", Value: "Engineering"}},
		DataSourceName: "employees",
	}

	expected := []RowMap{
		{"name": "John Doe", "age": "30"},
	}

	result := executeSqlSelectQuery(&data, query)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestExecuteSqlUpdateCommand_DataResult(t *testing.T) {
	data := [][]string{
		{"id", "name", "age", "department"},
		{"1", "John Doe", "30", "Engineering"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	query := SQLStatement{
		Type:           "UPDATE",
		Conditions:     []Condition{{Column: "name", Operator: "==", Value: "John Doe"}},
		SetClauses:     map[string]string{"age": "35"},
		DataSourceName: "employees",
	}

	expected := [][]string{
		{"id", "name", "age", "department"},
		{"1", "John Doe", "35", "Engineering"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	_ = executeSqlUpdateCommand(&data, query)

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}

func TestExecuteSqlUpdateCommand_RowCountResult(t *testing.T) {
	data := [][]string{
		{"id", "name", "age", "department"},
		{"1", "John Doe", "30", "Engineering"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	query := SQLStatement{
		Type:           "UPDATE",
		Conditions:     []Condition{{Column: "name", Operator: "==", Value: "John Doe"}},
		SetClauses:     map[string]string{"age": "35"},
		DataSourceName: "employees",
	}

	expected := 1

	result := executeSqlUpdateCommand(&data, query)

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestExecuteSqlDeleteCommand_DataResult(t *testing.T) {
	data := [][]string{
		{"id", "name", "age", "department"},
		{"1", "John Doe", "30", "Engineering"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	query := SQLStatement{
		Type:           "DELETE",
		Conditions:     []Condition{{Column: "age", Operator: "==", Value: "30"}},
		DataSourceName: "employees",
	}

	expected := [][]string{
		{"id", "name", "age", "department"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	_ = executeSqlDeleteCommand(&data, query)

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}

func TestExecuteSqlDeleteCommand_RowCountResult(t *testing.T) {
	data := [][]string{
		{"id", "name", "age", "department"},
		{"1", "John Doe", "30", "Engineering"},
		{"2", "Jane Smith", "40", "Marketing"},
	}

	query := SQLStatement{
		Type:           "DELETE",
		Conditions:     []Condition{{Column: "age", Operator: "==", Value: "30"}},
		DataSourceName: "employees",
	}

	expected := 1

	result := executeSqlDeleteCommand(&data, query)

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
