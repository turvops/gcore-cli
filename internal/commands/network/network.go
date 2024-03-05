package network

import (
	"fmt"
	"github.com/G-core/gcore-cli/internal/errors"
	"regexp"

	"github.com/spf13/cobra"

	cloud "github.com/G-Core/gcore-cloud-sdk-go"
	"github.com/G-core/gcore-cli/internal/core"
	"github.com/G-core/gcore-cli/internal/human"
)

var (
	client *cloud.ClientWithResponses

	projectID     int
	regionID      int
	waitForResult bool
)

var reNetworkName = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z 0-9._\\-]{1,61}[a-zA-Z0-9._]$")

type network struct {
	// Id Network ID
	Id string
	// Name Network name
	Name string
	// Type Network type (vlan, vxlan)
	Type string
	// External True if the network has `router:external` attribute
	External bool
	// Default True if the network has is_default attribute
	Default bool
	// Shared True when the network is shared with your project by external owner
	Shared bool
	// Mtu MTU (maximum transmission unit). Default value is 1450
	Mtu int
	// Subnets List of subnetworks
	Subnets []string
	// Metadata Network metadata
	Metadata []cloud.MetadataItemSchema
	// SegmentationId Id of network segment
	SegmentationId int
	// ProjectId Project ID
	ProjectId int
	// Region Region name
	Region string
	// RegionId Region ID
	RegionId int
	// CreatedAt Datetime when the network was created
	CreatedAt string
	// UpdatedAt Datetime when the network was last updated
	UpdatedAt string
}

func toView(instance cloud.NetworkSchema) network {
	return network{
		Id:             instance.Id,
		Name:           instance.Name,
		Type:           instance.Type,
		External:       instance.External,
		Default:        instance.Default,
		Shared:         instance.Shared,
		Mtu:            instance.Mtu,
		Subnets:        instance.Subnets,
		Metadata:       instance.Metadata,
		SegmentationId: instance.SegmentationId,
		ProjectId:      instance.ProjectId,
		Region:         instance.Region,
		RegionId:       instance.RegionId,
		CreatedAt:      instance.CreatedAt,
		UpdatedAt:      instance.UpdatedAt,
	}
}

func init() {
	human.RegisterMarshalerFunc(cloud.NetworkSchema{}, func(i interface{}, opt *human.MarshalOpt) (body string, err error) {
		instance := i.(cloud.NetworkSchema)
		s := toView(instance)

		return human.Marshal(s, nil)
	})

	human.RegisterMarshalerFunc([]cloud.NetworkSchema{}, func(i interface{}, opt *human.MarshalOpt) (body string, err error) {
		instances := i.([]cloud.NetworkSchema)
		s := make([]network, len(instances))
		for i, instance := range instances {
			s[i] = toView(instance)
		}

		return human.Marshal(s, nil)
	})
}

// top-level cloud network command
func Commands() *cobra.Command {
	// networkCmd represents the network command
	var networkCmd = &cobra.Command{
		Use:     "network",
		Short:   "Cloud network management commands",
		Long:    ``, // TODO:
		GroupID: "cloud",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		Args: cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			var (
				ctx = cmd.Context()
			)

			profile, err := core.GetClientProfile(ctx)
			if err != nil {
				return err
			}

			if profile.ApiKey == nil || *profile.ApiKey == "" {
				return &errors.CliError{
					Err:  fmt.Errorf("subcommand requires APIKEY token"),
					Hint: "See gcore-cli init, gcore-cli config",
				}
			}

			client, err = core.CloudClient(ctx)
			if err != nil {
				return err
			}

			waitForResult = cmd.Flag("wait").Value.String() == "true"

			return nil
		},
	}

	networkCmd.AddCommand(create(), show(), list(), rename(), delete())
	return networkCmd
}
