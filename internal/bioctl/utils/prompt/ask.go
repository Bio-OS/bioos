package prompt

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
)

type PromptOptions struct {
	message string
}

type PromptOption func(o *PromptOptions)

func WithInputMessage(message string) PromptOption {
	return func(o *PromptOptions) {
		o.message = message
	}
}

func PromptRequiredString(name string, opts ...PromptOption) (string, error) {
	return PromptStringWithValidator(name, survey.Required, opts...)
}

func PromptOptionalString(name string, opts ...PromptOption) (string, error) {
	return PromptStringWithValidator(name, OptionalInput, opts...)
}

func PromptStringWithValidator(name string, validator survey.Validator, opts ...PromptOption) (string, error) {
	res := ""
	msg := name
	options := &PromptOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.message != "" {
		msg = fmt.Sprintf("%s (%s)", name, options.message)
	}
	prompt := &survey.Input{
		Message: msg,
	}
	if err := survey.AskOne(prompt, &res, survey.WithValidator(validator), survey.WithIcons(func(icons *survey.IconSet) {
		// you can set any icons
		icons.Question.Text = SELECTICON
		// for more information on formatting the icons, see here: https://github.com/mgutz/ansi#style-format
		//icons.Question.Format = "yellow+hb"
		icons.Question.Format = "cyan"
	})); err != nil {
		fmt.Println(err)
		return "", utils.HandlePromptError(err)
	}
	return res, nil
}

// PromptStringSlice prompt for string slice
func PromptStringSlice(label string) ([]string, error) {
	res := make([]string, 0)

	var result string
	var err error
	for result != "quit" {
		result, err = PromptRequiredString(fmt.Sprintf("%s ('quit' to exit)", label))
		if err != nil {
			return nil, err
		}

		if strings.ToLower(result) != "quit" {
			res = append(res, result)
		}
	}

	return res, nil
}

func PromptRequiredInt32(name string, opts ...PromptOption) (int32, error) {
	var res int32 = 0
	msg := name
	options := &PromptOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.message != "" {
		msg = options.message
	}
	prompt := &survey.Input{
		Message: msg,
	}
	if err := survey.AskOne(prompt, &res, survey.WithValidator(ValidateIntegerNumberInput), survey.WithIcons(func(icons *survey.IconSet) {
		// you can set any icons
		icons.Question.Text = SELECTICON
		// for more information on formatting the icons, see here: https://github.com/mgutz/ansi#style-format
		//icons.Question.Format = "yellow+hb"
		icons.Question.Format = "cyan"
	})); err != nil {
		return 0, utils.HandlePromptError(err)
	}
	return res, nil
}
