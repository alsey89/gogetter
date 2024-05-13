package common

import "fmt"

// GetConfigPath returns scope.key in string format
func GetConfigPath(scope string, key string) string {
	return fmt.Sprintf("%s.%s", scope, key)
}
