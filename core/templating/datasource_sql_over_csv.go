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
		source, exits := datasource.GetDataSource(dataSourceName)
		if !exits {
			return SQLStatement{}, errors.New("data source does not exist")
		}

		if len(matches) == 4 {
			wherePart = matches[3]
		}

		headers := source.Data[0]
		columns, err := parseColumns(columnsPart, headers)
		if err != nil {
			return SQLStatement{}, err
		}

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
		source, exits := datasource.GetDataSource(dataSourceName)
		if !exits {
			return SQLStatement{}, errors.New("data source does not exist")
		}
		setPart := matches[2]
		if len(matches) == 4 {
			wherePart = matches[3]
		}
		headers := source.Data[0]
		setClauses, err := parseSetClauses(setPart, headers)
		if err != nil {
			return SQLStatement{}, err
		}
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
		if _, exits := datasource.GetDataSource(dataSourceName); !exits {
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
func parseColumns(columnsPart string, headers []string) ([]string, error) {
	columnsPart = strings.TrimSpace(columnsPart)
	if columnsPart == "*" {
		return headers, nil
	}
	columns := strings.Split(columnsPart, ",")
	for i, column := range columns {
		if !stringExists(headers, strings.TrimSpace(column)) {
			return nil, errors.New("invalid column provided: " + strings.TrimSpace(column))
		}
		columns[i] = strings.TrimSpace(column)
	}
	return columns, nil
}

func stringExists(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
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
func parseSetClauses(setPart string, headers []string) (map[string]string, error) {
	setClauses := make(map[string]string)
	parts, err := splitOnCommasOutsideQuotes(setPart)
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		keyValue := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(keyValue) != 2 {
			return nil, errors.New("invalid SET clause: " + part)
		}
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		if !stringExists(headers, key) {
			return nil, errors.New("invalid column provided: " + key)
		}
		setClauses[key] = trimQuotes(value)
	}
	return setClauses, nil
}

// splitOnCommasOutsideQuotes splits a string by commas only if the commas are outside of quotes
func splitOnCommasOutsideQuotes(s string) ([]string, error) {
	var parts []string
	var sb strings.Builder
	var inQuotes bool
	var quoteChar rune

	for _, r := range s {
		switch {
		case r == '\'' || r == '"':
			if inQuotes {
				if r == quoteChar {
					inQuotes = false
				}
			} else {
				inQuotes = true
				quoteChar = r
			}
			sb.WriteRune(r)
		case r == ',' && !inQuotes:
			parts = append(parts, sb.String())
			sb.Reset()
		default:
			sb.WriteRune(r)
		}
	}

	// Add the last part
	parts = append(parts, sb.String())

	return parts, nil
}

// parseConditions parses the WHERE part of the query into a slice of Conditions and returns an error if any issues are found.
func parseConditions(wherePart string) ([]Condition, error) {
	conditions := []Condition{}

	conditionRegex := regexp.MustCompile(`(\w+)\s*(==|=|!=|<=|>=|<|>)\s*'([^']*)'`)
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
func executeSqlUpdateCommand(data *[][]string, query SQLStatement) int {
	if len(*data) < 2 {
		log.Debug("no data available to update")
		return 0
	}

	headers := (*data)[0]
	conditions := query.Conditions
	setClauses := query.SetClauses
	rowsAffected := 0
	for i, row := range (*data)[1:] {
		rowMap := mapRow(headers, row)
		if matchesConditions(rowMap, conditions) {
			rowsAffected += 1
			for column, newValue := range setClauses {
				colIndex := indexOf(headers, column)
				if colIndex != -1 {
					(*data)[i+1][colIndex] = newValue
				}
			}
		}
	}
	return rowsAffected
}

// executeSqlDeleteCommand executes a DELETE query and modifies the data in-place
func executeSqlDeleteCommand(data *[][]string, query SQLStatement) int {
	if len(*data) < 2 {
		log.Debug("no data available to delete")
		return 0
	}

	headers := (*data)[0]
	conditions := query.Conditions
	rowsAffected := 0
	// Iterate in reverse to avoid index shifting issues
	for i := len(*data) - 1; i > 0; i-- {
		row := (*data)[i]
		rowMap := mapRow(headers, row)
		if matchesConditions(rowMap, conditions) {
			removeRow(data, i)
			rowsAffected++
		}
	}
	return rowsAffected
}

// removeRow removes the row at index rowIndex from the data.
func removeRow(data *[][]string, rowIndex int) {
	if rowIndex < 0 || rowIndex >= len(*data) {
		// Return early if rowIndex is out of bounds
		return
	}
	// Overwrite the row at rowIndex with the next rows
	copy((*data)[rowIndex:], (*data)[rowIndex+1:])
	// Truncate the slice to remove the last row which is now duplicate
	*data = (*data)[:len(*data)-1]
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
		case "==", "=":
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
