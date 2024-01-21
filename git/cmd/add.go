package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to the index",
	Run: func(cmd *cobra.Command, args []string) {
		add()
	},
}

func add() {
	fmt.Println("add")
}
