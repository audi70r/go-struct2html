// Package struct2html provides a convenient way to convert Go structs or slices of structs
// into HTML table representations. These tables are designed to be responsive and
// email-friendly, making them suitable for use in automated email reporting systems,
// client communications, and any application that requires structured data to be presented
// in a clear, tabular format within an email.
package struct2html

import (
	"fmt"
	"html"
	"reflect"
	"strings"
)

// StructToHTMLTable takes a struct or a slice of structs and converts it into an email-friendly HTML table representation.
func StructToHTMLTable(data interface{}) (string, error) {
	val := reflect.Indirect(reflect.ValueOf(data))
	if val.Kind() != reflect.Struct && !(val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Struct) {
		return "", fmt.Errorf("expected a struct or a slice of structs but got a %s", val.Kind())
	}

	var htmlTable strings.Builder
	// Set width to 100% to make the table responsive
	htmlTable.WriteString(`<table border="0" cellpadding="5" cellspacing="0" style="border-collapse: collapse; width: 100%;">` + "\n")

	if val.Kind() == reflect.Struct {
		htmlTable.WriteString(generateTableRows(val))
	} else {
		for i := 0; i < val.Len(); i++ {
			item := reflect.Indirect(val.Index(i)) // Ensure we're working with the concrete struct, not a pointer
			htmlTable.WriteString(generateTableRows(item))
		}
	}

	htmlTable.WriteString("</table>")
	return htmlTable.String(), nil
}

// generateTableRows generates the table rows for a single struct, handling nested structs and pointers.
func generateTableRows(val reflect.Value) string {
	var rows strings.Builder

	// Create headers
	rows.WriteString(`<tr style="background-color: #f8f8f8;">` + "\n")
	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := valType.Field(i)
		if !field.IsExported() { // Skip unexported fields
			continue
		}
		rows.WriteString(`<th style="border: 1px solid #ddd; padding: 5px;">`)
		rows.WriteString(html.EscapeString(field.Name))
		rows.WriteString("</th>")
	}
	rows.WriteString("</tr>\n")

	// Create row data
	rows.WriteString(`<tr style="background-color: #ffffff;">` + "\n")
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		if !valType.Field(i).IsExported() {
			continue
		}

		rows.WriteString(`<td style="border: 1px solid #ddd; padding: 5px;">`)
		if fieldValue.Kind() == reflect.Struct {
			// Recursive call for nested structs
			rows.WriteString(`<table border="0" cellpadding="5" cellspacing="0" style="border-collapse: collapse; width: 100%;">`)
			rows.WriteString(generateTableRows(fieldValue))
			rows.WriteString("</table>")
		} else if fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct {
			// Recursive call for pointers to structs, if not nil
			if !fieldValue.IsNil() {
				rows.WriteString(`<table border="0" cellpadding="5" cellspacing="0" style="border-collapse: collapse; width: 100%;">`)
				rows.WriteString(generateTableRows(fieldValue.Elem()))
				rows.WriteString("</table>")
			}
		} else {
			// Print field value for non-struct fields
			rows.WriteString(html.EscapeString(fmt.Sprint(fieldValue.Interface())))
		}
		rows.WriteString("</td>")
	}
	rows.WriteString("</tr>\n")

	return rows.String()
}
