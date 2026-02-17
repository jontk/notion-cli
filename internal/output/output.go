package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func JSON(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func Table(headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return nil
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	for i, h := range headers {
		fmt.Printf("%-*s", widths[i]+2, h)
	}
	fmt.Println()

	// Print separator
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w+2))
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf("%-*s", widths[i]+2, cell)
			}
		}
		fmt.Println()
	}

	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func Error(err error) error {
	resp := ErrorResponse{
		Error: err.Error(),
		Code:  "ERROR",
	}
	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	if encErr := encoder.Encode(resp); encErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}
