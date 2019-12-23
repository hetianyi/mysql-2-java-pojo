package command

import (
	"fmt"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/mysql-2-java-pojo/common"
	"github.com/hetianyi/mysql-2-java-pojo/util"
	"github.com/hetianyi/mysql-2-java-pojo/worker"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

// Parse parses command flags using `github.com/urfave/cli`
func Parse(arguments []string) {
	appFlag := cli.NewApp()
	appFlag.Version = "0.0.1"
	appFlag.HideVersion = true
	appFlag.Name = "Mysql Reverse Tool"
	appFlag.Usage = "generator"
	appFlag.HelpName = "generator"
	// config file location
	appFlag.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c,config",
			Value:       "",
			Usage:       `配置文件地址`,
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "o,output",
			Value:       "",
			Usage:       `工作空间(文件夹)`,
			Destination: &configFile,
		},
	}

	cli.AppHelpTemplate = `
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .VisibleCommands}}

Commands:{{range .VisibleCategories}}
{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

Options:

   {{range $index, $option := .VisibleFlags}}{{if $index}}{{end}}{{$option}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = `
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

{{.Usage}}{{if .VisibleFlags}}

Options:

   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	cli.SubcommandHelpTemplate = `
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

{{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}

Commands:
{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{if .VisibleFlags}}

Options:

   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	appFlag.Action = func(c *cli.Context) error {
		return nil
	}

	err := appFlag.Run(arguments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	if workspace == "" {
		workspace, _ = file.GetWorkDir()
	}

	if outputDir == "" {
		outputDir = workspace + "/output"
	}

	if configFile == "" {
		configFile = workspace + "/config.yml"
	}

	bs, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	config := &common.Config{}

	if err = util.ParseYamlFromString(bs, config); err != nil {
		logger.Fatal(err)
	}
	worker.Start(outputDir, config)
}
