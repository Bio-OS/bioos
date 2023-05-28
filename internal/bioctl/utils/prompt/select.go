package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
)

type GetItemsFunc func(...interface{}) ([]interface{}, error)

// PromptStringSelect prompt select from slice
func PromptStringSelect(label string, size int, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("%s  {{ . | cyan }}", SELECTICON),
			Inactive: "  {{ . | cyan }}",
			Selected: fmt.Sprintf("%s  {{ . | cyan }}", SELECTICON),
		},
		Size: size,
	}
	i, _, err := prompt.Run()

	if err != nil {
		return "", utils.HandlePromptError(err)
	}

	return items[i], nil
}

func PromptBoolSelect(label string) (bool, error) {
	items := []bool{true, false}
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("%s  {{ . | cyan }}", SELECTICON),
			Inactive: "  {{ . | cyan }}",
			Selected: fmt.Sprintf("%s  {{ . | cyan }}", SELECTICON),
		},
		Size: 2,
	}
	i, _, err := prompt.Run()

	if err != nil {
		return false, utils.HandlePromptError(err)
	}

	return items[i], nil
}

// PromptSelectWithFunc prompt select from slice with func
func PromptSelectWithFunc(label string, size int, getItemsFunc GetItemsFunc, funcArgs ...interface{}) (interface{}, error) {
	items, err := getItemsFunc(funcArgs...)
	if err != nil {
		return nil, utils.HandlePromptError(err)
	}
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("%s  {{ . | cyan }}", SELECTICON),
			Inactive: "  {{ . | cyan }}",
			Selected: fmt.Sprintf("%s  {{ . | cyan }}", SELECTICON),
		},
		Size: size,
	}
	i, _, err := prompt.Run()

	if err != nil {
		return nil, utils.HandlePromptError(err)
	}

	return items[i], nil
}

// PromptStringMultiSelect prompt select from slice
func PromptStringMultiSelect(label string, size int, items []string) ([]string, error) {
	res := make([]string, 0)
	prompt := &survey.MultiSelect{
		Message:  label,
		Options:  items,
		PageSize: size,
	}
	err := survey.AskOne(prompt, &res, survey.WithIcons(func(icons *survey.IconSet) {
		// you can set any icons
		icons.Question.Text = SELECTICON
		// for more information on formatting the icons, see here: https://github.com/mgutz/ansi#style-format
		//icons.Question.Format = "yellow+hb"
		icons.Question.Format = "cyan"
	}))
	if err != nil {
		return nil, utils.HandlePromptError(err)
	}

	return res, nil
}
