package helm

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/deislabs/porter/pkg/printer"
	yaml "gopkg.in/yaml.v2"
)

type StatusAction struct {
	Steps []StatusStep `yaml:"status"`
}

// StatusStep represents the structure of an Status action
type StatusStep struct {
	StatusArguments `yaml:"helm"`
}

// StatusArguments are the arguments available for the Status action
type StatusArguments struct {
	Step `yaml:",inline"`

	Releases []string `yaml:"releases"`
}

// Status reports the status for a provided set of Helm releases
func (m *Mixin) Status(opts printer.PrintOptions) error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var action StatusAction
	err = yaml.Unmarshal(payload, &action)
	if err != nil {
		return err
	}
	if len(action.Steps) != 1 {
		return errors.Errorf("expected a single step, but got %d", len(action.Steps))
	}
	step := action.Steps[0]

	format := ""
	switch opts.Format {
	case printer.FormatPlaintext:
		// do nothing, as default output is plaintext
	case printer.FormatYaml:
		format = `-o yaml`
	case printer.FormatJson:
		format = `-o json`
	default:
		return fmt.Errorf("invalid format: %s", opts.Format)
	}

	for _, release := range step.Releases {
		cmd := m.NewCommand("helm", "status", strings.TrimSpace(fmt.Sprintf(`%s %s`, release, format)))

		cmd.Stdout = m.Out
		cmd.Stderr = m.Err

		prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
		fmt.Fprintln(m.Out, prettyCmd)

		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("could not execute command, %s: %s", prettyCmd, err)
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
