package util

import (
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

func PromptSelect(label string, options []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}
	index, _, err := prompt.Run()
	return index, options[index], errors.Wrap(err, "error running select promt")
}

func PromptOptionalText(label string) (string, error) {
	getValue := promptui.Prompt{
		Label: label,
	}
	value, err := getValue.Run()
	return value, errors.Wrap(err, "error running optional text promt")
}
