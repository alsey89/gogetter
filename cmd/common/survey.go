package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
)

func SelectMultipleOptions(message string, options []string) []string {
	var selectedOptions []string

	// set up multi select prompt
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Help:    "Use the arrow keys to navigate and space to select. Press enter when done.",
	}

	// run the prompt, save result to selectedOptions
	err := survey.AskOne(prompt, &selectedOptions)
	if err != nil {
		log.Fatalf("Failed to select options: %v", err)
	}

	return selectedOptions
}

func SelectOneOption(message string, options []string) string {
	var selectedOption string

	// set up select prompt
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Help:    "Use the arrow keys to navigate and space to select. Press enter when done.",
	}

	// run the prompt, save result to selectedOption
	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		log.Fatalf("Failed to select option: %v", err)
	}

	return selectedOption
}

func SelectYesNo(message string, defaultSelection bool) bool {
	var selectedOption bool

	// set up confirm prompt
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultSelection,
	}

	// run the prompt, save result to selectedOption
	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		log.Fatalf("Failed to select option: %v", err)
	}

	return selectedOption
}

func InputText(message string, isMandatory bool) string {
	var input string
	var opts []survey.AskOpt

	prompt := &survey.Input{
		Message: message,
	}

	if isMandatory {
		// this makes the input a required field
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	// run the prompt, save result to input
	err := survey.AskOne(prompt, &input, opts...)
	if err != nil {
		log.Fatalf("Failed to input text: %v", err)
	}

	return input
}
