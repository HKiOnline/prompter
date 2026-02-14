package templa

import (
	"strings"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	// Test that Date returns a string in the correct format
	result := Date()

	// Parse the result to verify it's a valid date
	_, err := time.Parse("2006-01-02", result)
	if err != nil {
		t.Errorf("Date() returned invalid format: %v, got: %s", err, result)
	}

	// Test that it returns today's date
	expected := time.Now().Format("2006-01-02")
	if result != expected {
		t.Errorf("Date() returned %s, expected %s", result, expected)
	}
}

func TestCreateTemplate(t *testing.T) {
	// Test that createTemplate registers the date function
	tmpl := createTemplate("test")

	// Test that the template can use the date function
	testTemplate := `Today is {{date}}`
	_, err := tmpl.Parse(testTemplate)
	if err != nil {
		t.Fatalf("Failed to parse template with date function: %v", err)
	}

	// Execute the template
	var result strings.Builder
	err = tmpl.Execute(&result, nil)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	// Verify the result contains a date
	if result.Len() == 0 {
		t.Error("Template execution produced empty result")
	}
}
