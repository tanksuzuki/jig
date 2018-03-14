package config

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Imports []Import `toml:"import"`
	Scripts []Script `toml:"script"`
}

type Import struct {
	Path   string `toml:"path"`
	Prefix string `toml:"prefix"`
}

func (c Config) validate() error {
	hasWhitespace := regexp.MustCompile(`\s`)
	duplicateCheck := map[string]struct{}{}

	for _, script := range c.Scripts {
		if hasWhitespace.MatchString(script.Name) {
			return fmt.Errorf("script name cannot contain whitespace: \"%s\"", script.Name)
		}
		if _, ok := duplicateCheck[script.Name]; ok {
			return fmt.Errorf("script name must be unique: \"%s\"", script.Name)
		}
		duplicateCheck[script.Name] = struct{}{}
	}

	return nil
}

func (c Config) get(name string) (Script, error) {
	for _, s := range c.Scripts {
		if s.Name == name {
			return s, nil
		}
	}
	return Script{}, fmt.Errorf("script not found")
}

func (c Config) List() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, ' ', 0)
	for _, script := range c.Scripts {
		description := strings.Split(script.Description, "\n")[0]
		if description == "" {
			description = "n/a"
		}
		fmt.Fprintf(w, "%s %s\t%s\n", script.Name, color.YellowString(script.Args), description)
	}
	w.Flush()
}

func (c Config) Help(name string) error {
	script, err := c.get(name)
	if err != nil {
		return err
	}
	fmt.Println(strings.TrimSpace(script.Help()))
	return nil
}

func (c Config) Run(name string, args []string) (int, error) {
	script, err := c.get(name)
	if err != nil {
		return -1, err
	}
	return script.Run(args), nil
}

func Read(path, prefix string) (Config, error) {
	c := Config{}

	homedirExpandedPath, err := homedir.Expand(path)
	if err != nil {
		return c, err
	}

	if _, err := toml.DecodeFile(homedirExpandedPath, &c); err != nil {
		return c, err
	}

	for _, i := range c.Imports {
		imported, err := Read(i.Path, i.Prefix)
		if err != nil {
			return c, err
		}
		c.Scripts = append(c.Scripts, imported.Scripts...)
	}

	scriptCount := len(c.Scripts)
	for i := 0; i < scriptCount; i++ {
		c.Scripts[i].Name = prefix + c.Scripts[i].Name
	}

	if err := c.validate(); err != nil {
		return c, err
	}

	sort.Slice(c.Scripts, func(i int, j int) bool {
		return c.Scripts[i].Name < c.Scripts[j].Name
	})

	return c, nil
}

type Script struct {
	Name        string `toml:"name"`
	Exec        string `toml:"exec"`
	Args        string `toml:"args"`
	Description string `toml:"description"`
}

func (s Script) Help() string {
	var help string

	lines := strings.Split(s.Description, "\n")

	if len(lines) > 0 {
		help += lines[0] + "\n\n"
	}

	help += fmt.Sprintf("Usage:\n  jig %s %s\n", s.Name, s.Args)

	if len(lines) > 1 {
		help += "\n"
		for i := 1; i < len(lines); i++ {
			help += lines[i] + "\n"
		}
	}

	return help
}

func (s Script) Run(args []string) int {
	args = append([]string{"-c", s.Exec, s.Name}, args...)

	cmd := exec.Command("sh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if sw, ok := ee.Sys().(syscall.WaitStatus); ok {
				return sw.ExitStatus()
			}
		}
		return 1
	}
	return 0
}
