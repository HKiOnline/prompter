package templa

import (
	"strings"
	"text/template"
)

// TemplateData holds the data for template execution
type TemplateData struct {
	// Add fields as needed for template data
}

// processTemplate processes template content with arguments
func Process(content string, args map[string]string) string {
	// Create template with built-in functions
	tmpl := createTemplate("prompt")

	// Parse the template
	parsedTemplate, err := tmpl.Parse(content)
	if err != nil {
		// If template parsing fails, return original content
		return content
	}

	// Execute template with arguments
	var result strings.Builder
	err = parsedTemplate.Execute(&result, convertArgsToInterface(args))
	if err != nil {
		// If template execution fails, return original content
		return content
	}

	return result.String()
}

// convertArgsToInterface converts string arguments to interface{}
func convertArgsToInterface(args map[string]string) map[string]interface{} {
	if args == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range args {
		result[k] = v
	}
	return result
}

// CreateTemplate creates a new template with built-in functions
func createTemplate(name string) *template.Template {
	t := template.New(name)
	// Register built-in functions
	t = t.Funcs(template.FuncMap{
		"date": Date,
	})
	return t
}
