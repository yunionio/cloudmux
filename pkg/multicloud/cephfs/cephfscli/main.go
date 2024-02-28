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

package cephfscli

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/http/httpproxy"

	"yunion.io/x/pkg/util/shellutils"
	"yunion.io/x/structarg"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/cephfs"
	_ "yunion.io/x/cloudmux/pkg/multicloud/cephfs/shell"
)

type BaseOptions struct {
	Debug      bool   `help:"debug mode"`
	Username   string `help:"Username" default:"$CEPHFS_USERNAME" metavar:"CEPHFS_USERNAME"`
	Password   string `help:"Password" default:"$CEPHFS_PASSWORD" metavar:"CEPHFS_PASSWORD"`
	Port       int    `help:"Service port" default:"$CEPHFS_PORT|8443" metavar:"CEPHFS_PORT"`
	Host       string `help:"RegionId" default:"$CEPHFS_HOST" metavar:"CEPHFS_HOST"`
	SUBCOMMAND string `help:"cephfscli subcommand" subcommand:"true"`
}

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	parse, e := structarg.NewArgumentParserWithHelp(&BaseOptions{},
		"cephfscli",
		"Command-line interface to cephfs API.",
		`See "cephfscli COMMAND --help" for help on a specific command.`)

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

func newClient(options *BaseOptions) (*cephfs.SCephFSClient, error) {
	if len(options.Username) == 0 {
		return nil, fmt.Errorf("Missing username")
	}

	if len(options.Password) == 0 {
		return nil, fmt.Errorf("Missing password")
	}

	if len(options.Host) == 0 {
		return nil, fmt.Errorf("Missing host")
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

	cli, err := cephfs.NewCephFSClient(
		cephfs.NewCephFSClientConfig(
			options.Host,
			options.Port,
			options.Username,
			options.Password,
			"",
		).Debug(options.Debug).
			CloudproviderConfig(
				cloudprovider.ProviderConfig{
					ProxyFunc: proxyFunc,
				},
			),
	)
	if err != nil {
		return nil, err
	}
	return cli, nil
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
	var cli *cephfs.SCephFSClient
	cli, e = newClient(options)
	if e != nil {
		showErrorAndExit(e)
	}
	e = subcmd.Invoke(cli, suboptions)
	if e != nil {
		showErrorAndExit(e)
	}
}
