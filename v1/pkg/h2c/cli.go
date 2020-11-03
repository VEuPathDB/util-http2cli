package h2c

import cli "github.com/Foxcapades/Argonaut/v0"

const (
	descTools = "List of tools this service is allowed to call.\n\n" +
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
)

var (
	defTools  = []string{"blastn", "blastp", "blastx", "tblastn", "tblastx"}
	defDbDir  = "/db"
	defOutDir = "/out"
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
				Default(defTools))).
		Flag(cli.NewFlag().
			Short('d').
			Long("db-dir").
			Description(descDbDir).
			Arg(cli.NewArg().
				Bind(&config.DbDir).
				Require().
				Default(defDbDir))).
		Flag(cli.NewFlag().
			Short('o').
			Long("out-dir").
			Description(descOutDir).
			Arg(cli.NewArg().
				Bind(&config.OutDir).
				Require().
				Default(defOutDir))).
		MustBuild()
}
