package configuration

import (
	"bufio"
	"os"
	"strings"
)

type EnvConfiguration struct {
	Listen string `env:"listen"`
}

func (c *EnvConfiguration) Configure(path string) error {
	values := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.setValue(values, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// TODO: make it reflection later and validate data
	c.Listen = values["listen"]

	return file.Close()
}

func (c *EnvConfiguration) setValue(values map[string]string, text string) {
	trimmed := strings.Split(text, "#")[0]
	line := strings.Split(trimmed, "=")
	if len(line) == 2 {
		values[line[0]] = line[1]
	}
}
