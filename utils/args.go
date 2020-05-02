package utils

import "strings"

// ParseArgs are defined for parsing parameters
func ParseArgs(osArgs []string) ([]string, []string) {
	content := make([]string, 0)
	args := []string{osArgs[0]}
	lastArg := ""
	for _, arg := range osArgs[1:] {
		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
			lastArg = arg
			continue
		}
		if strings.HasPrefix(lastArg, "-") && len(content) != 0 {
			continue
		}
		content = append(content, arg)
		lastArg = arg
	}
	return content, args[1:]
}