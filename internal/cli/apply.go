package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewApplyCmd returns the apply command.
func NewApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply patch (unified diff) with preview/rollback",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = getApp(cmd)
			fmt.Println("Apply not implemented yet (Phase D).")
			return nil
		},
	}
}
