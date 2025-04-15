/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/amberbyte/flamigo/tools/flamigo/internal/project"
	"github.com/spf13/cobra"
)

type DomainTemplateData struct {
	DomainName string
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a new domain or interface",
	Long:    `Adds a new domain or interface to the project`,
	Example: `flamigo add domain <domain-name>`,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "domain":
			workingDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			rootPath, err := project.DetermineProjectRoot(workingDir)
			if err != nil {
				panic(err)
			}
			domainName := args[1]
			data := DomainTemplateData{
				DomainName: domainName,
			}
			err = project.InitializeDomainFiles(path.Join(rootPath, "internal", "domain", domainName), data)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created domain %s in %s\n", domainName, path.Join(rootPath, "internal", "domain", domainName))
			fmt.Println("You must add your domains app, and infrastructure to your cmd/main.go")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Args = cobra.ExactArgs(2)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
