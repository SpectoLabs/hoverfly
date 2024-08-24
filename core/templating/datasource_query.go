package templating

import (
	"errors"
	"regexp"
	"strings"
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
type SelectQuery struct {
	Type           string // "SELECT", "UPDATE", or "DELETE"
	Columns        []string
	Conditions     []Condition
	SetClauses     map[string]string // For UPDATE queries
	DataSourceName string
}

func parseQuery(query string, datasource map[string]*DataSource) (SelectQuery, error) {
	query = strings.TrimSpace(query)

	var queryType string
	if strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		queryType = "SELECT"
	} else if strings.HasPrefix(strings.ToUpper(query), "UPDATE") {
		queryType = "UPDATE"
	} else if strings.HasPrefix(strings.ToUpper(query), "DELETE") {
		queryType = "DELETE"
	} else {
		return SelectQuery{}, errors.New("invalid query type")
	}

	var selectRegex *regexp.Regexp
	var matches []string
	var columnsPart, dataSourceName, wherePart string

	switch queryType {
	case "SELECT":
		selectRegex = regexp.MustCompile(`(?i)^SELECT\s+(.+)\s+FROM\s+(\w+)\s*(?:WHERE\s+(.+))?$`)
		matches = selectRegex.FindStringSubmatch(query)
		if len(matches) < 3 {
			return SelectQuery{}, errors.New("invalid query format")
		}
		columnsPart = matches[1]
		dataSourceName = matches[2]
		if len(matches) == 4 {
			wherePart = matches[3]
		}
	case "UPDATE":
		//selectRegex = regexp.MustCompile(`(?i)^UPDATE\s+(\w+)\s+SET\s+(.+)\s*(?:WHERE\s+(.+))?$`)
		selectRegex = regexp.MustCompile(`(?i)^UPDATE\s+(\w+)\s+SET\s+(.+?)(?:\s+WHERE\s+(.+))?$`)

		matches = selectRegex.FindStringSubmatch(query)
		if len(matches) < 3 {
			return SelectQuery{}, errors.New("invalid UPDATE query format")
		}
		dataSourceName = matches[1]
		setPart := matches[2]
		if len(matches) == 4 {
			wherePart = matches[3]
		}
		// log.Debug("########## ParseQuery ############")
		// logMessage := fmt.Sprintf("len(matches)=%d", len(matches))
		// log.Debug(logMessage)
		// logMessage = fmt.Sprintf("matches[3]=" + matches[3])
		// log.Debug(logMessage)
		// formattedMatches := strings.Join(matches, " --- ")
		// logMessage = fmt.Sprintf("matches=%v", formattedMatches)
		// log.Debug(logMessage)
		// log.Debug("######################")
		setClauses := parseSetClauses(setPart)
		conditions, err := parseConditions(wherePart)
		if err != nil {
			return SelectQuery{}, err
		}
		return SelectQuery{
			Type:           "UPDATE",
			Conditions:     conditions,
			SetClauses:     setClauses,
			DataSourceName: dataSourceName,
		}, nil

	case "DELETE":
		selectRegex = regexp.MustCompile(`(?i)^DELETE\s+FROM\s+(\w+)\s*(?:WHERE\s+(.+))?$`)
		matches = selectRegex.FindStringSubmatch(query)
		if len(matches) < 2 {
			return SelectQuery{}, errors.New("invalid DELETE query format")
		}
		dataSourceName = matches[1]
		if len(matches) == 3 {
			wherePart = matches[2]
		}
		conditions, err := parseConditions(wherePart)
		if err != nil {
			return SelectQuery{}, err
		}
		return SelectQuery{
			Type:           "DELETE",
			Conditions:     conditions,
			DataSourceName: dataSourceName,
		}, nil
	}

	headers := datasource[dataSourceName].Data[0]
	columns := parseColumns(columnsPart, headers)

	conditions, err := parseConditions(wherePart)
	if err != nil {
		return SelectQuery{}, err
	}

	return SelectQuery{
		Type:           queryType,
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

// TrimQuotes trims matching single or double quotes from the outer edges of a string.
func trimQuotes(s string) string {
	// Check if the string length is less than 2, in which case there's nothing to trim.
	if len(s) < 2 {
		return s
	}

	// Get the first and last characters of the string.
	firstChar := s[0]
	lastChar := s[len(s)-1]

	// Check if both the first and last characters are matching quotes.
	if (firstChar == '\'' || firstChar == '"') && firstChar == lastChar {
		return s[1 : len(s)-1]
	}

	// If no trimming is necessary, return the original string.
	return s
}

// parseSetClauses parses the SET part of an UPDATE query
func parseSetClauses(setPart string) map[string]string {
	setClauses := make(map[string]string)
	parts := strings.Split(setPart, ",")
	for _, part := range parts {
		keyValue := strings.Split(strings.TrimSpace(part), "=")
		if len(keyValue) == 2 {
			setClauses[strings.TrimSpace(keyValue[0])] = trimQuotes(strings.TrimSpace(keyValue[1]))
		}
	}
	return setClauses
}

// parseConditions parses the WHERE part of the query into a slice of Conditions and returns an error if any issues are found.
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

// ExecuteQuery processes different types of SQL-like queries based on the query type
// func executeQuery(data *[][]string, query SelectQuery) ([]RowMap, error) {
// 	switch query.Type {
// 	case "SELECT":
// 		return executeSelectQuery(*data, query), nil
// 	case "UPDATE":
// 		if err := executeUpdateQuery(data, query); err != nil {
// 			return nil, err
// 		}
// 		return nil, nil // UPDATE doesn't return rows
// 	case "DELETE":
// 		if err := executeDeleteQuery(data, query); err != nil {
// 			return nil, err
// 		}
// 		return nil, nil // DELETE doesn't return rows
// 	default:
// 		return nil, errors.New("unsupported query type")
// 	}
// }

// ExecuteSelectQuery executes a SELECT query and returns the results as a slice of RowMaps
func executeSelectQuery(data [][]string, query SelectQuery) []RowMap {
	headers := data[0] // First row as header
	results := []RowMap{}
	for _, row := range data[1:] {
		rowMap := mapRow(headers, row)
		if matchesConditions(rowMap, query.Conditions) {
			results = append(results, projectRow(rowMap, query.Columns))
		}
	}
	return results
}

// ExecuteUpdateQuery executes an UPDATE query and modifies the data in-place
func executeUpdateQuery(data *[][]string, query SelectQuery) error {
	if len(*data) < 2 {
		return errors.New("no data available to update")
	}

	headers := (*data)[0]
	conditions := query.Conditions
	setClauses := query.SetClauses
	for i, row := range (*data)[1:] {
		rowMap := mapRow(headers, row)
		if matchesConditions(rowMap, conditions) {
			for column, newValue := range setClauses {
				colIndex := indexOf(headers, column)
				if colIndex != -1 {
					(*data)[i+1][colIndex] = newValue
				}
			}
		}
	}
	return nil
}

// ExecuteDeleteQuery executes a DELETE query and modifies the data in-place
func executeDeleteQuery(data *[][]string, query SelectQuery) error {
	if len(*data) < 2 {
		return errors.New("no data available to delete")
	}

	headers := (*data)[0]
	conditions := query.Conditions
	filteredData := [][]string{headers}
	for _, row := range (*data)[1:] {
		rowMap := mapRow(headers, row)
		if !matchesConditions(rowMap, conditions) {
			filteredData = append(filteredData, row)
		}
	}
	*data = filteredData
	return nil
}

// Helper function to find the index of a column header
func indexOf(headers []string, column string) int {
	for i, header := range headers {
		if header == column {
			return i
		}
	}
	return -1
}

// mapRow converts a CSV row into a map with column names as keys
func mapRow(headers, row []string) RowMap {
	rowMap := make(RowMap)
	for i, header := range headers {
		rowMap[header] = row[i]
	}
	return rowMap
}

// mapToRow converts a map back to a CSV row
// func mapToRow(headers []string, rowMap RowMap) []string {
// 	row := make([]string, len(headers))
// 	for i, header := range headers {
// 		row[i] = rowMap[header]
// 	}
// 	return row
// }

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
			return false
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
		}
	}
	return true
}
