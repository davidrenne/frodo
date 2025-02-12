package cli

import (
	"log"

	"github.com/davidrenne/frodo/generate"
	"github.com/davidrenne/frodo/parser"
	"github.com/spf13/cobra"
)

// GenerateGatewayRequest contains all of the CLI options used in the "frodo client" command.
type GenerateGatewayRequest struct {
	templateOption
	// InputFileName is the service definition to parse/process (the "--service" option)
	InputFileName string
}

// GenerateGateway handles the registration and execution of the 'frodo gateway' CLI subcommand.
type GenerateGateway struct{}

// Command creates the Cobra struct describing this CLI command and its options.
func (c GenerateGateway) Command() *cobra.Command {
	request := &GenerateGatewayRequest{}
	cmd := &cobra.Command{
		Use:   "gateway [flags] FILENAME",
		Short: "Process a Go source file with your service interface to generate an RPC/API gateway.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			request.InputFileName = args[0]
			crapPants(c.Exec(request))
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().StringVar(&request.Template, "template", "", "Path to a custom Go template file used to generate this artifact.")
	return cmd
}

// Exec actually executes the parsing/generating logic creating the gateway for the given declaration.
func (c GenerateGateway) Exec(request *GenerateGatewayRequest) error {
	log.Printf("Parsing service definitions: %s", request.InputFileName)
	ctx, err := parser.ParseFile(request.InputFileName)
	if err != nil {
		return err
	}

	artifact := request.ToFileTemplate("gateway.go")
	log.Printf("Generating artifact '%s'", artifact.Name)
	return generate.File(ctx, artifact)
}
