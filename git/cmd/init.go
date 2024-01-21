package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// InitCmd wyag init [path]
var InitCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize new repository",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		initRepo(args)
	},
}

func initRepo(args []string) {
	// getwd, err := os.Getwd()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	cmd.NewGitRepository(getwd)

	fmt.Println(args)
	fmt.Println("init")
}

type GitRepository struct {
	// WorkTree path to project where we want to init repository
	WorkTree string

	// GitDir
	GitDir string

	// Conf
	Conf string
}

func NewGitRepository(repositoryPath string) *GitRepository {
	repo := &GitRepository{
		WorkTree: repositoryPath,
		GitDir:   path.Join(repositoryPath, ".git"),
	}

	if !isDir(repo.GitDir) {
		fmt.Println(repo.GitDir, "not dir")
	}

	return repo
}

func isDir(gitDir string) bool {
	fileInfo, err := os.Stat(gitDir)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	return fileInfo.IsDir()
}
