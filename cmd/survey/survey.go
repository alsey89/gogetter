package survey

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
)

func SelectMultipleOptions(message string, options []string) ([]string, error) {
	var selectedOptions []string

	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Help:    "Use the arrow keys to navigate and space to select. Press enter when done.",
	}

	err := survey.AskOne(prompt, &selectedOptions)
	if err != nil {
		return nil, err
	}

	return selectedOptions, nil
}

func SelectOneOption(message string, options []string) (*string, error) {
	var selectedOption string

	prompt := &survey.Select{
		Message: message,
		Options: options,
		Help:    "Use the arrow keys to navigate and space to select. Press enter when done.",
	}

	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		return nil, err
	}

	return &selectedOption, nil
}

func SelectYesNo(message string, defaultSelection bool) (bool, error) {
	var selectedOption bool

	prompt := &survey.Confirm{
		Message: message,
		Default: defaultSelection,
	}

	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		return false, err
	}

	return selectedOption, nil
}

func InputText(message string, isMandatory bool) (*string, error) {
	var input string
	var opts []survey.AskOpt

	prompt := &survey.Input{
		Message: message,
	}

	if isMandatory {
		// this makes the input a required field
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	err := survey.AskOne(prompt, &input, opts...)
	if err != nil {
		log.Fatalf("error asking for input: %v", err)
	}

	return &input, nil
}
