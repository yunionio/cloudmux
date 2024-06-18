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

package ocli

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/http/httpproxy"

	"yunion.io/x/pkg/util/shellutils"
	"yunion.io/x/structarg"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/oracle"
	_ "yunion.io/x/cloudmux/pkg/multicloud/oracle/shell"
)

type BaseOptions struct {
	Debug         bool   `help:"debug mode"`
	TenancyOCID   string `help:"Tenancy OCID" default:"$ORACLE_TENANCY_OCID" metavar:"ORACLE_TENANCY_OCID"`
	UserOCID      string `help:"User OCID" default:"$ORACLE_USER_OCID" metavar:"ORACLE_USER_OCID"`
	CompartmentId string `help:"Compartment Id" default:"$ORACLE_COMPARTMENT_ID" metavar:"ORACLE_COMPARTMENT_ID"`
	PrivateKey    string `help:"Private Key" default:"$ORACLE_PRIVATE_KEY" metavar:"ORACLE_PRIVATE_KEY"`
	RegionId      string `help:"RegionId" default:"$ORACLE_REGION_ID|ap-singapore-1" metavar:"ORACLE_REGION_ID"`
	SUBCOMMAND    string `help:"oraclecli subcommand" subcommand:"true"`
}

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	parse, e := structarg.NewArgumentParserWithHelp(&BaseOptions{},
		"ocli",
		"Command-line interface to oracle API.",
		`See "ocli COMMAND --help" for help on a specific command.`)

	if e != nil {
		return nil, e
	}

	subcmd := parse.GetSubcommand()
	if subcmd == nil {
		return nil, fmt.Errorf("No subcommand argument.")
	}
	for _, v := range shellutils.CommandTable {
		_, e := subcmd.AddSubParserWithHelp(v.Options, v.Command, v.Desc, v.Callback)
		if e != nil {
			return nil, e
		}
	}
	return parse, nil
}

func showErrorAndExit(e error) {
	fmt.Fprintf(os.Stderr, "%s", e)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func newClient(options *BaseOptions) (*oracle.SRegion, error) {
	if len(options.TenancyOCID) == 0 {
		return nil, fmt.Errorf("Missing tenancy ocid")
	}

	if len(options.UserOCID) == 0 {
		return nil, fmt.Errorf("Missing user ocid")
	}

	if len(options.PrivateKey) == 0 {
		return nil, fmt.Errorf("Missing private key")
	}

	cfg := &httpproxy.Config{
		HTTPProxy:  os.Getenv("HTTP_PROXY"),
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		NoProxy:    os.Getenv("NO_PROXY"),
	}
	cfgProxyFunc := cfg.ProxyFunc()
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		return cfgProxyFunc(req.URL)
	}

	opts, err := oracle.NewOracleClientConfig(
		options.TenancyOCID,
		options.UserOCID,
		options.CompartmentId,
		options.PrivateKey,
	)
	if err != nil {
		return nil, err
	}
	cli, err := oracle.NewOracleClient(
		opts.Debug(options.Debug).
			CloudproviderConfig(
				cloudprovider.ProviderConfig{
					ProxyFunc: proxyFunc,
					RegionId:  options.RegionId,
				},
			),
	)
	if err != nil {
		return nil, err
	}

	return cli.GetRegion(options.RegionId)
}

func Main() {
	parser, e := getSubcommandParser()
	if e != nil {
		showErrorAndExit(e)
	}

	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*BaseOptions)

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
	var region *oracle.SRegion
	region, e = newClient(options)
	if e != nil {
		showErrorAndExit(e)
	}
	e = subcmd.Invoke(region, suboptions)
	if e != nil {
		showErrorAndExit(e)
	}
}
