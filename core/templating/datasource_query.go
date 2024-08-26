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

// SQLStatement represents a simple SQL-like query
type SQLStatement struct {
	Type           string // "SELECT", "UPDATE", or "DELETE"
	Columns        []string
	Conditions     []Condition
	SetClauses     map[string]string // For UPDATE queries
	DataSourceName string
}

func parseSqlCommand(query string, datasource *TemplateDataSource) (SQLStatement, error) {
	query = strings.TrimSpace(query)
	var commandType string
	if strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		commandType = "SELECT"
	} else if strings.HasPrefix(strings.ToUpper(query), "UPDATE") {
		commandType = "UPDATE"
	} else if strings.HasPrefix(strings.ToUpper(query), "DELETE") {
		commandType = "DELETE"
	} else {
		return SQLStatement{}, errors.New("invalid query type")
	}

	var selectRegex *regexp.Regexp
	var matches []string
	var columnsPart, dataSourceName, wherePart string

	switch commandType {
	case "SELECT":
		selectRegex = regexp.MustCompile(`(?i)^SELECT\s+(.+)\s+FROM\s+([\w-]+)\s*(?:WHERE\s+(.+))?$`)
		matches = selectRegex.FindStringSubmatch(query)
		if len(matches) < 3 {

			return SQLStatement{}, errors.New("invalid query format")

		}
		columnsPart = matches[1]
		dataSourceName = matches[2]
		if !dataSourceExists(datasource, dataSourceName) {
			return SQLStatement{}, errors.New("data source does not exist")
		}

		if len(matches) == 4 {
			wherePart = matches[3]
		}
		headers := datasource.DataSources[dataSourceName].Data[0]
		columns := parseColumns(columnsPart, headers)

		conditions, err := parseConditions(wherePart)
		if err != nil {
			return SQLStatement{}, err
		}

		return SQLStatement{
			Type:           commandType,
			Columns:        columns,
			Conditions:     conditions,
			DataSourceName: dataSourceName,
		}, nil
	case "UPDATE":
		selectRegex = regexp.MustCompile(`(?i)^UPDATE\s+([\w-]+)\s+SET\s+(.+?)(?:\s+WHERE\s+(.+))?$`)
		matches = selectRegex.FindStringSubmatch(query)
		if len(matches) < 3 {
			return SQLStatement{}, errors.New("invalid UPDATE query format")
		}
		dataSourceName = matches[1]
		if !dataSourceExists(datasource, dataSourceName) {
			return SQLStatement{}, errors.New("data source does not exist")
		}
		setPart := matches[2]
		if len(matches) == 4 {
			wherePart = matches[3]
		}

		setClauses := parseSetClauses(setPart)
		conditions, err := parseConditions(wherePart)
		if err != nil {
			return SQLStatement{}, err
		}
		return SQLStatement{
			Type:           "UPDATE",
			Conditions:     conditions,
			SetClauses:     setClauses,
			DataSourceName: dataSourceName,
		}, nil

	case "DELETE":
		selectRegex = regexp.MustCompile(`(?i)^DELETE\s+FROM\s+([\w-]+)\s*(?:WHERE\s+(.+))?$`)
		matches = selectRegex.FindStringSubmatch(query)
		if len(matches) < 2 {
			return SQLStatement{}, errors.New("invalid DELETE query format")
		}
		dataSourceName = matches[1]
		if !dataSourceExists(datasource, dataSourceName) {
			return SQLStatement{}, errors.New("data source does not exist")
		}
		if len(matches) == 3 {
			wherePart = matches[2]
		}
		conditions, err := parseConditions(wherePart)
		if err != nil {
			return SQLStatement{}, err
		}
		return SQLStatement{
			Type:           "DELETE",
			Conditions:     conditions,
			DataSourceName: dataSourceName,
		}, nil
	}
	return SQLStatement{}, errors.New("invalid query format")
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
		return conditions, nil
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

// ExecuteSelectQuery executes a SELECT query and returns the results as a slice of RowMaps
func executeSqlSelectQuery(data *[][]string, query SQLStatement) []RowMap {
	headers := (*data)[0] // First row as header
	results := []RowMap{}
	for _, row := range (*data)[1:] {
		rowMap := mapRow(headers, row)
		if matchesConditions(rowMap, query.Conditions) {
			results = append(results, projectRow(rowMap, query.Columns))
		}
	}
	return results
}

// ExecuteUpdateQuery executes an UPDATE query and modifies the data in-place
func executeSqlUpdateCommand(data *[][]string, query SQLStatement) error {
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
func executeSqlDeleteCommand(data *[][]string, query SQLStatement) error {
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
