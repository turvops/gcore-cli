package fastedge

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cobra"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

var client *sdk.ClientWithResponses

// top-level FastEdge command
func Commands(baseUrl string, authFunc func(ctx context.Context, req *http.Request) error) (*cobra.Command, error) {
	var local bool
	var cmdFastedge = &cobra.Command{
		Use:   "fastedge <subcommand>",
		Short: "Gcore Edge compute solution",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			url := baseUrl
			if !local {
				url += "/fastedge"
			}
			client, err = sdk.NewClientWithResponses(
				url,
				sdk.WithRequestEditorFn(authFunc),
			)
			if err != nil {
				return fmt.Errorf("cannot init SDK: %w", err)
			}

			carbon.SetDefault(carbon.Default{
				Timezone: carbon.UTC,
				Locale:   "en",
			})

			return nil
		},
	}
	cmdFastedge.PersistentFlags().BoolVar(&local, "local", false, "local testing")
	cmdFastedge.PersistentFlags().MarkHidden("local")

	cmdFastedge.AddCommand(app(), binary(), plan(), stat())
	return cmdFastedge, nil
}

func newPointer[T any](val T) *T {
	return &val
}
