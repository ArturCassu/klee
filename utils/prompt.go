package utils

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Confirm prompts user for yes/no confirmation
func Confirm(message string) bool {
	result := false
	prompt := &survey.Confirm{
		Message: message,
	}
	survey.AskOne(prompt, &result)
	return result
}

// PromptInput prompts user for text input
func PromptInput(message string, defaultValue string) string {
	result := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}
	survey.AskOne(prompt, &result)
	return result
}

// PromptSelect prompts user to select from options with a default
func PromptSelect(message string, options []string, defaultIndex int) string {
	if len(options) == 0 {
		return ""
	}

	if defaultIndex < 0 || defaultIndex >= len(options) {
		defaultIndex = 0
	}

	result := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Default: options[defaultIndex],
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		// Return default on error
		return options[defaultIndex]
	}

	return result
}

// PromptMultiSelect prompts user to select multiple options using arrow keys and space
func PromptMultiSelect(message string, options []string) ([]string, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("no options provided")
	}

	result := []string{}
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
