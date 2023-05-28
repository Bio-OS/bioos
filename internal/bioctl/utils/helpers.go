package utils

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	DefaultErrorExitCode = 1
)

var fatalErrHandler = fatal
var PromptManualExitSignal = fmt.Errorf("prompt")

func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

var ErrExit = fmt.Errorf("exit")

func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

func checkErr(err error, handleErr func(string, int)) {

	if err == nil {
		return
	}

	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	case err == PromptManualExitSignal:
		handleErr("Prompt exit", DefaultErrorExitCode)
	default:
		msg, ok := StandardErrorMessage(err)
		if !ok {
			msg = err.Error()
			if !strings.HasPrefix(msg, "error: ") {
				msg = fmt.Sprintf("error: %s", msg)
			}
		}
		handleErr(msg, DefaultErrorExitCode)
	}
}

func StandardErrorMessage(err error) (string, bool) {
	switch t := err.(type) {
	case *url.Error:
		switch {
		case strings.Contains(t.Err.Error(), "connection refused"):
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			return fmt.Sprintf("The connection to the server %s was refused - did you specify the right host or port?", host), true
		}
		return fmt.Sprintf("Unable to connect to the server: %v", t.Err), true
	}
	return "", false
}

func HandlePromptError(err error) error {
	if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrAbort) {
		return PromptManualExitSignal
	}
	return fmt.Errorf("Prompt failed %v\n", err)
}
