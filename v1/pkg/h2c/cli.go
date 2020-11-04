package h2c

import (
	"fmt"
	"os"

	"github.com/Foxcapades/Argonaut/v0"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

const (
	descTools = "List of tools this service is allowed to call.\n\n" +
		"This list may be specified by using this flag more than once.\n\n" +
		"Defaults to [blastn, blastp, blastx, tblastn, tblastx]."
	descDbDir = "Root directory this service will use as the base path when" +
		" configuring a blast tool run.\n\n" +
		"This should match the volume mount point set when starting the docker" +
		" container.\n\n" +
		"Defaults to /db"
	descOutDir = "Output directory this service will use as the base path when" +
		" configuring a blast tool run.\n\n" +
		"This should match the volume mount point set when starting the docker" +
		" container.\n\n" +
		"Defaults to /out"
	descPort = "Port the HTTP server should bind to."
	descVers = "Prints the server binary version."
)

var (
	DefaultTools  = []string{"blastn", "blastp", "blastx", "tblastn", "tblastx"}
	DefaultDbDir  = "/db"
	DefaultOutDir = "/out"
	DefaultPort   = uint16(80)
)

func InitCLI(config *Config) {
	cli.NewCommand().
		Flag(cli.NewFlag().
			Short('t').
			Long("tools").
			Description(descTools).
			Arg(cli.NewArg().
				Bind(&config.Tools).
				Require().
				Name("cmd").
				Default(DefaultTools))).
		Flag(cli.NewFlag().
			Short('d').
			Long("db-dir").
			Description(descDbDir).
			Arg(cli.NewArg().
				Bind(&config.DbDir).
				Require().
				Name("path").
				Default(DefaultDbDir))).
		Flag(cli.NewFlag().
			Short('o').
			Long("out-dir").
			Description(descOutDir).
			Arg(cli.NewArg().
				Bind(&config.OutDir).
				Require().
				Name("path").
				Default(DefaultOutDir))).
		Flag(cli.NewFlag().
			Short('p').
			Long("port").
			Description(descPort).
			Arg(cli.NewArg().
				Bind(&config.ServerPort).
				Require().
				Name("port").
				Default(DefaultPort))).
		Flag(cli.NewFlag().
			Short('V').
			Long("version").
			Description(descVers).
			OnHit(func(argo.Flag) {
				fmt.Println(config.Version)
				os.Exit(0)
			})).
		MustParse()
}
