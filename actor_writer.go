package actor

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

func (a *Actor) Write(w io.Writer) error {

	writer := newWriter(w)

	if err := writer.writeTags(a.Tags); err != nil {
		return fmt.Errorf("Write initial tags: %s", err)
	}

	if err := writer.writeKeyword("Actor", a.Name); err != nil {
		return fmt.Errorf("Write actor keyword: %s", err)
	}

	writer.indent()

	for _, blurb := range a.Blurb {
		if err := writer.writeBlurb(blurb); err != nil {
			return fmt.Errorf("Write blurbs: %s", err)
		}
	}

	if len(a.Blurb) > 0 {
		if err := writer.newLine(); err != nil {
			return fmt.Errorf("New line: %s", err)
		}
	}

	goalsWithTags := 0

	// Goals with tags
	for _, goal := range a.Goals {

		if len(goal.Tags) > 0 {

			if err := writer.writeTags(goal.Tags); err != nil {
				return fmt.Errorf("Write goal tags: %s", err)
			}

			if err := writer.writeKeyword("Goal", goal.Name); err != nil {
				return fmt.Errorf("Write goal name: %s", err)
			}

			goalsWithTags++
		}
	}

	if goalsWithTags > 0 && (len(a.Goals)-goalsWithTags) > 0 {

		if err := writer.newLine(); err != nil {
			return fmt.Errorf("New line: %s", err)
		}

		if err := writer.writeKeyword("Goals", ""); err != nil {
			return fmt.Errorf("Write goals tag: %s", err)
		}

		writer.indent()

		// Goals without tags
		for _, goal := range a.Goals {

			if len(goal.Tags) == 0 {

				if err := writer.writeBlurb(goal.Name); err != nil {
					return fmt.Errorf("Write goal name: %s", err)
				}

				goalsWithTags++
			}
		}
	}

	return nil
}

func (a *Actor) WriteToFile(name string) error {

	buf := &bytes.Buffer{}

	if err := a.Write(buf); err != nil {
		return err
	}

	return ioutil.WriteFile(name, buf.Bytes(), 0644)
}
