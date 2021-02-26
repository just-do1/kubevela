package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/oam-dev/kubevela/apis/types"
	cmdutil "github.com/oam-dev/kubevela/pkg/utils/util"
	"github.com/oam-dev/kubevela/references/plugins"
)

// NewWorkloadsCommand creates `workloads` command
func NewWorkloadsCommand(c types.Args, ioStreams cmdutil.IOStreams) *cobra.Command {
	var enforceRefresh bool
	ctx := context.Background()
	cmd := &cobra.Command{
		Use:                   "workloads",
		DisableFlagsInUseLine: true,
		Short:                 "List workloads",
		Long:                  "List workloads",
		Example:               `vela workloads`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return c.SetConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := RefreshDefinitions(ctx, c, ioStreams, true, enforceRefresh); err != nil {
				return err
			}

			workloads, err := plugins.LoadInstalledCapabilityWithType(types.TypeWorkload)
			if err != nil {
				return err
			}
			return printWorkloadList(workloads, ioStreams)
		},
		Annotations: map[string]string{
			types.TagCommandType: types.TypeCap,
		},
	}
	cmd.SetOut(ioStreams.Out)
	cmd.Flags().BoolVarP(&enforceRefresh, "", "r", false, "Enforce refresh from cluster even if cache is not expired")
	return cmd
}

func printWorkloadList(workloadList []types.Capability, ioStreams cmdutil.IOStreams) error {
	table := newUITable()
	table.MaxColWidth = 120
	table.AddRow("NAME", "DESCRIPTION")
	for _, r := range workloadList {
		table.AddRow(r.Name, r.Description)
	}
	ioStreams.Info(table.String())
	return nil
}
