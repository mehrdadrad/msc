package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	cli "github.com/urfave/cli/v2"
	yml "gopkg.in/yaml.v3"
)

var version string

type config struct {
	Workers  int
	Timeout  int
	Interval int
	Sites    []string
	Method   string
}

type crawlerCLIRequest struct {
	Config string
}

func loadConfig(r *crawlerCLIRequest) (*config, error) {
	f, err := os.Open(r.Config)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	c := &config{
		Workers:  10,
		Method:   "HEAD",
		Timeout:  5,
		Interval: 30,
	}
	err = yml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getCLI(args []string) (*crawlerCLIRequest, error) {
	var r = &crawlerCLIRequest{}

	initCLI()

	app := &cli.App{
		Version: version,
		Flags:   flags,
		Action:  action(r),
	}

	err := app.Run(args)

	return r, err
}

var flags = []cli.Flag{
	&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "", Usage: "path to a file in yaml format to read configuration"},
}

func action(r *crawlerCLIRequest) cli.ActionFunc {
	return func(c *cli.Context) error {
		r.Config = c.String("config")

		return nil
	}
}

func initCLI() {
	cli.AppHelpTemplate = `usage: simple crawler options
	
options:

   {{range .VisibleFlags}}{{.}}
   {{end}}
`

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Simple Crawler version: %s\n", c.App.Version)
		cli.OsExiter(0)
	}

	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		funcMap := template.FuncMap{
			"join": strings.Join,
		}
		t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
		t.Execute(w, data)
		cli.OsExiter(0)
	}
}
