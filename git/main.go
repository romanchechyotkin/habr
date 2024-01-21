package main

import (
	"fmt"
	"github.com/romanchechyotkin/habr/git/cmd"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wyag",
	Short: "A simple Git-like CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide a valid subcommand. Use 'wyag --help' for usage")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	rootCmd.AddCommand(cmd.InitCmd, cmd.AddCmd)
	Execute()
}

//case "add"          : cmd_add(args)
//case "cat-file"     : cmd_cat_file(args)
//case "check-ignore" : cmd_check_ignore(args)
//case "checkout"     : cmd_checkout(args)
//case "commit"       : cmd_commit(args)
//case "hash-object"  : cmd_hash_object(args)
//case "log"          : cmd_log(args)
//case "ls-files"     : cmd_ls_files(args)
//case "ls-tree"      : cmd_ls_tree(args)
//case "rev-parse"    : cmd_rev_parse(args)
//case "rm"           : cmd_rm(args)
//case "show-ref"     : cmd_show_ref(args)
//case "status"       : cmd_status(args)
//case "tag"          : cmd_tag(args)
//case _              : print("Bad command.")
