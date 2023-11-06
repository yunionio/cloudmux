// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmx

import (
	"fmt"
	"os"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/structarg"

	"yunion.io/x/cloudmux/pkg/multicloud/cmx/shell"
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

func Main() {
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
