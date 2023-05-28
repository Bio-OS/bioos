package notebook

import (
	"os"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	applog "github.com/Bio-OS/bioos/pkg/log"
)

func TestMain(m *testing.M) {
	applog.RegisterLogger(&applog.Options{
		Level: "fatal",
	})
	os.Exit(m.Run())
}

func TestValidateIPythonNotebook(t *testing.T) {
	cases := []struct {
		content     string
		expectError bool
	}{
		{
			content: `{
  "cells": [
	{
	  "cell_type": "code",
	  "execution_count": null,
	  "id": "ee2a969b",
	  "metadata": {},
	  "outputs": [],
	  "source": []
	}
  ],
  "metadata": {
	"kernelspec": {
	  "display_name": "Python 3 (ipykernel)",
	  "language": "python",
	  "name": "python3"
	},
	"language_info": {
	  "codemirror_mode": {
		"name": "ipython",
		"version": 3
	  },
	  "file_extension": ".py",
	  "mimetype": "text/x-python",
	  "name": "python",
	  "nbconvert_exporter": "python",
	  "pygments_lexer": "ipython3",
	  "version": "3.9.7"
	}
  },
  "nbformat": 4,
  "nbformat_minor": 5
}`,
			expectError: false,
		},
	}
	g := gomega.NewWithT(t)
	for _, c := range cases {
		var match types.GomegaMatcher
		if c.expectError {
			match = gomega.HaveOccurred()
		} else {
			match = gomega.BeNil()
		}
		g.Expect(validateIPythonNotebook([]byte(c.content))).To(match)
	}
}
