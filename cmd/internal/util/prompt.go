package util

import (
	"github.com/manifoldco/promptui"
)

func PromptSelect(label string, options []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return 0, "", err
	}

	return index, options[index], nil
}

func PromptOptionalText(label string) (string, error) {
	getValue := promptui.Prompt{
		Label: label,
	}

	value, err := getValue.Run()
	if err != nil {
		return "", err
	}

	return value, nil
}
