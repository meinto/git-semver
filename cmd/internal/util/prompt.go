package util

import "github.com/manifoldco/promptui"

func PromptSelect(label string, options []string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return index, nil
}
