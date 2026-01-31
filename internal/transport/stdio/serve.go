package stdio

import (
	"fmt"

	"github.com/hkionline/prompter/internal/configuration"
)

// Serve for stdio transport layer (old implementation - kept for compatibility)
func Serve(configuration configuration.Configuration) {
	fmt.Println("Warning: Using old Serve function - MCP SDK now handles stdio transport directly")
	// This will be removed after full migration
}
