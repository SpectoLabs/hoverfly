package templating

import (
	"errors"
	"regexp"
	"strings"
)

// A mock structure representing your CSV data
// var data = [][]string{
// 	{"id", "name", "age", "city"},
// 	{"1", "John", "30", "New York"},
// 	{"2", "Jane", "25", "Los Angeles"},
// 	{"3", "Doe", "40", "New York"},
// }

// RowMap represents a single row in the result set
type RowMap map[string]string

// Condition represents a single condition in the WHERE clause
type Condition struct {
	Column   string
	Operator string
	Value    string
}

// SelectQuery represents a simple SQL-like SELECT query
type SelectQuery struct {
	Columns    []string
	Conditions []Condition
}

// ParseQuery parses a query string into a SelectQuery struct
func ParseQuery(query string) (SelectQuery, error) {
	query = strings.TrimSpace(query)

	selectRegex := regexp.MustCompile(`(?i)^SELECT\s+(.+)\s+WHERE\s+(.+)$`)
	matches := selectRegex.FindStringSubmatch(query)
	if len(matches) != 3 {
		return SelectQuery{}, errors.New("invalid query format")
	}

	columnsPart := matches[1]
	wherePart := matches[2]

	columns := strings.Split(columnsPart, ",")
	for i, column := range columns {
		columns[i] = strings.TrimSpace(column)
	}

	conditions, err := parseConditions(wherePart)
	if err != nil {
		return SelectQuery{}, err
	}

	return SelectQuery{
		Columns:    columns,
		Conditions: conditions,
	}, nil
}

// parseConditions parses the WHERE part of the query into a slice of Conditions
func parseConditions(wherePart string) ([]Condition, error) {
	conditions := []Condition{}

	conditionRegex := regexp.MustCompile(`(\w+)\s*(==|!=)\s*'([^']*)'`)
	conditionMatches := conditionRegex.FindAllStringSubmatch(wherePart, -1)

	if len(conditionMatches) == 0 {
		return nil, errors.New("no valid conditions found")
	}

	for _, match := range conditionMatches {
		if len(match) != 4 {
			return nil, errors.New("invalid condition format")
		}
		conditions = append(conditions, Condition{
			Column:   match[1],
			Operator: match[2],
			Value:    match[3],
		})
	}

	return conditions, nil
}

// ExecuteQuery runs the query against the in-memory data
func ExecuteQuery(data [][]string, query SelectQuery) []RowMap {
	headers := data[0] // first row as header
	results := []RowMap{}

	for _, row := range data[1:] {
		rowMap := mapRow(headers, row)
		if matchesConditions(rowMap, query.Conditions) {
			results = append(results, projectRow(rowMap, query.Columns))
		}
	}
	return results
}

// mapRow converts a CSV row into a map with column names as keys
func mapRow(headers, row []string) RowMap {
	rowMap := make(RowMap)
	for i, header := range headers {
		rowMap[header] = row[i]
	}
	return rowMap
}

// projectRow filters the row based on the selected columns
func projectRow(row RowMap, columns []string) RowMap {
	if len(columns) == 0 {
		return row
	}
	projected := make(RowMap)
	for _, col := range columns {
		if val, ok := row[col]; ok {
			projected[col] = val
		}
	}
	return projected
}

// matchesConditions checks if a row matches all given conditions
func matchesConditions(row RowMap, conditions []Condition) bool {
	for _, condition := range conditions {
		val, ok := row[condition.Column]
		if !ok {
			return false // Column doesn't exist
		}

		switch condition.Operator {
		case "==":
			if val != condition.Value {
				return false
			}
		case "!=":
			if val == condition.Value {
				return false
			}
		default:
			return false // Unsupported operator
		}
	}
	return true // All conditions matched
}

// Example usage
// func main() {
// 	queryStr := "SELECT age, city WHERE age!='30' AND city=='New York'"

// 	query, err := ParseQuery(queryStr)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	results := ExecuteQuery(data, query)
// 	for _, row := range results {
// 		fmt.Println(row)
// 	}
// }
