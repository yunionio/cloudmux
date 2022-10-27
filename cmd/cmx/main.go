package main

import (
	"fmt"
	"os"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/structarg"

	"yunion.io/x/cloudmux/cmd/cmx/shell"
	_ "yunion.io/x/cloudmux/pkg/multicloud/loader"
)

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	parse, err := structarg.NewArgumentParserWithHelp(&shell.GlobalOptions{},
		"cmuxcli",
		"Command-line interface to cloudmux API.",
		`See "cmuxcli COMMAND --help" for help on a specific command.`)
	if err != nil {
		return nil, err
	}

	subcmd := parse.GetSubcommand()
	if subcmd == nil {
		return nil, errors.Errorf("No subcommand argument.")
	}

	for _, v := range *shell.CommandTable {
		if _, err := subcmd.AddSubParserWithHelp(v.Options, v.Command, v.Desc, v.Callback); err != nil {
			return nil, err
		}
	}

	return parse, nil
}

func showErrorAndExit(e error) {
	fmt.Fprintf(os.Stderr, "%s", e)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func main() {
	parser, e := getSubcommandParser()
	if e != nil {
		showErrorAndExit(e)
	}
	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*shell.GlobalOptions)

	if parser.IsHelpSet() {
		fmt.Print(parser.HelpString())
		return
	}
	subcmd := parser.GetSubcommand()
	subparser := subcmd.GetSubParser()
	if e != nil || subparser == nil {
		if subparser != nil {
			fmt.Print(subparser.Usage())
		} else {
			fmt.Print(parser.Usage())
		}
		showErrorAndExit(e)
		return
	}
	suboptions := subparser.Options()
	if subparser.IsHelpSet() {
		fmt.Print(subparser.HelpString())
		return
	}
	provider, e := shell.NewCloudProvider(options)
	if e != nil {
		showErrorAndExit(e)
	}

	if e := subcmd.Invoke(provider, suboptions); e != nil {
		showErrorAndExit(e)
	}
}
