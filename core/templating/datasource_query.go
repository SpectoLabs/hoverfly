package templating

import (
	"errors"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RowMap represents a single row in the result set
type RowMap map[string]string

// Condition represents a single condition in the WHERE clause
type Condition struct {
	Column   string
	Operator string
	Value    string
}

// SelectQuery represents a simple SQL-like SELECT query
//
//	type SelectQuery struct {
//		Columns    []string
//		Conditions []Condition
//	}
type SelectQuery struct {
	Columns        []string
	Conditions     []Condition
	DataSourceName string
}

// ParseQuery parses a query string into a SelectQuery struct
// func ParseQuery(query string, headers []string) (SelectQuery, error) {
// 	query = strings.TrimSpace(query)

// 	selectRegex := regexp.MustCompile(`(?i)^SELECT\s+(.+)\s+WHERE\s+(.+)$`)
// 	matches := selectRegex.FindStringSubmatch(query)
// 	if len(matches) != 3 {
// 		return SelectQuery{}, errors.New("invalid query format")
// 	}

// 	columnsPart := matches[1]
// 	wherePart := matches[2]

// 	columns := parseColumns(columnsPart, headers)
// 	conditions, err := parseConditions(wherePart)
// 	if err != nil {
// 		return SelectQuery{}, err
// 	}

//		return SelectQuery{
//			Columns:    columns,
//			Conditions: conditions,
//		}, nil
//	}
func ParseQuery(query string, datasource map[string]*DataSource) (SelectQuery, error) {
	query = strings.TrimSpace(query)

	// Regular expression to match SELECT, FROM, and WHERE clauses
	selectRegex := regexp.MustCompile(`(?i)^SELECT\s+(.+)\s+FROM\s+(\w+)\s*(?:WHERE\s+(.+))?$`)
	matches := selectRegex.FindStringSubmatch(query)
	if len(matches) < 3 {
		return SelectQuery{}, errors.New("invalid query format")
	}

	columnsPart := matches[1]
	dataSourceName := matches[2]
	wherePart := ""
	if len(matches) == 4 {
		wherePart = matches[3]
	}

	// Assuming headers can be retrieved from the data source at a later stage
	//headers := []string{} // This will be populated once the data source is retrieved
	headers := datasource[dataSourceName].Data[0]
	log.Debug("###################")
	log.Debug(headers)
	log.Debug("###################")
	// Determine the columns to select
	columns := parseColumns(columnsPart, headers)

	// Parse the WHERE clause if present
	var conditions []Condition
	var err error
	if wherePart != "" {
		conditions, err = parseConditions(wherePart)
		if err != nil {
			return SelectQuery{}, err
		}
	}

	return SelectQuery{
		Columns:        columns,
		Conditions:     conditions,
		DataSourceName: dataSourceName,
	}, nil
}

// parseColumns determines the columns to select based on the query part and headers
func parseColumns(columnsPart string, headers []string) []string {
	columnsPart = strings.TrimSpace(columnsPart)
	if columnsPart == "*" {
		return headers
	}
	columns := strings.Split(columnsPart, ",")
	for i, column := range columns {
		columns[i] = strings.TrimSpace(column)
	}
	return columns
}

// parseConditions parses the WHERE part of the query into a slice of Conditions
func parseConditions(wherePart string) ([]Condition, error) {
	conditions := []Condition{}

	conditionRegex := regexp.MustCompile(`(\w+)\s*(==|!=|<=|>=|<|>)\s*'([^']*)'`)
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
		case "<":
			if isGreaterThanOrEqual(val, condition.Value) {
				return false
			}
		case "<=":
			if isGreaterThan(val, condition.Value) {
				return false
			}
		case ">":
			if isLessThanOrEqual(val, condition.Value) {
				return false
			}
		case ">=":
			if isLessThan(val, condition.Value) {
				return false
			}
		default:
			return false // Unsupported operator
		}
	}
	return true // All conditions matched
}
