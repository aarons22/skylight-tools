package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

// PrintTable renders a JSON:API response as a formatted table.
// It handles:
//   - {"data": [...]}  — array of JSON:API resources
//   - {"data": {...}}  — single JSON:API resource
//   - [...]            — plain array of objects
//
// Falls back to raw JSON output if parsing fails.
func PrintTable(data []byte, noColor bool) error {
	// Try JSON:API list response
	var listResp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &listResp); err == nil && listResp.Data != nil {
		// Try as array
		var items []json.RawMessage
		if err := json.Unmarshal(listResp.Data, &items); err == nil {
			return printResources(items, noColor)
		}
		// Try as single object
		var item map[string]json.RawMessage
		if err := json.Unmarshal(listResp.Data, &item); err == nil {
			return printResources([]json.RawMessage{listResp.Data}, noColor)
		}
	}

	// Try plain array
	var rawItems []json.RawMessage
	if err := json.Unmarshal(data, &rawItems); err == nil {
		return printPlainArray(rawItems, noColor)
	}

	// Fallback: print raw JSON
	fmt.Println(string(data))
	return nil
}

// printResources prints JSON:API resource objects as a table by flattening
// id + attributes fields into columns.
func printResources(items []json.RawMessage, noColor bool) error {
	if len(items) == 0 {
		fmt.Println("(no results)")
		return nil
	}

	type resource struct {
		ID         string                     `json:"id"`
		Type       string                     `json:"type"`
		Attributes map[string]json.RawMessage `json:"attributes"`
	}

	var rows []map[string]string
	colSet := make(map[string]struct{})

	for _, raw := range items {
		var r resource
		if err := json.Unmarshal(raw, &r); err != nil {
			// Not a JSON:API resource — treat as plain object
			return printPlainArray(items, noColor)
		}

		row := make(map[string]string)
		row["id"] = r.ID
		colSet["id"] = struct{}{}

		for k, v := range r.Attributes {
			row[k] = jsonValueToString(v)
			colSet[k] = struct{}{}
		}
		rows = append(rows, row)
	}

	return renderTable(rows, colSet, noColor)
}

// printPlainArray prints a plain JSON array of objects as a table.
func printPlainArray(items []json.RawMessage, noColor bool) error {
	if len(items) == 0 {
		fmt.Println("(no results)")
		return nil
	}

	var rows []map[string]string
	colSet := make(map[string]struct{})

	for _, raw := range items {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			fmt.Println(string(raw))
			continue
		}
		row := make(map[string]string)
		for k, v := range obj {
			row[k] = jsonValueToString(v)
			colSet[k] = struct{}{}
		}
		rows = append(rows, row)
	}

	return renderTable(rows, colSet, noColor)
}

func renderTable(rows []map[string]string, colSet map[string]struct{}, noColor bool) error {
	// Sort columns: id first, then alphabetical
	cols := make([]string, 0, len(colSet))
	for k := range colSet {
		if k != "id" {
			cols = append(cols, k)
		}
	}
	sort.Strings(cols)
	if _, hasID := colSet["id"]; hasID {
		cols = append([]string{"id"}, cols...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cols)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("  ")
	table.SetRowSeparator("")
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)

	if noColor {
		table.SetHeaderColor()
	}

	for _, row := range rows {
		vals := make([]string, len(cols))
		for i, col := range cols {
			vals[i] = row[col]
		}
		table.Append(vals)
	}

	table.Render()
	return nil
}

// jsonValueToString converts a raw JSON value to a human-readable string.
func jsonValueToString(v json.RawMessage) string {
	if len(v) == 0 {
		return ""
	}
	// Null
	if string(v) == "null" {
		return ""
	}
	// String — unquote
	var s string
	if err := json.Unmarshal(v, &s); err == nil {
		return s
	}
	// Bool, number, array, object — use raw JSON
	return string(v)
}
