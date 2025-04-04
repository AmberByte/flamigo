/*
Copyright © 2025 AmberByte
*/
package cmd

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/amberbyte/flamigo/tools/flamigo/internal/project"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var features = []huh.Option[string]{
	huh.NewOption("Authentication (Domain)", "auth"),
	huh.NewOption("Realtime - Websocket (Interface)", "realtime"),
	// TODO: must complete implementation
	// huh.NewOption("Configuration (Core)", "config"),
}

type FormData struct {
	modulePath string
	features   []string
	ok         bool
}

type TemplateData struct {
	ProjectModulePath string
	ProjectName       string
	Features          []string
}

func (t TemplateData) HasFeature(feature string) bool {
	return slices.Contains(t.Features, feature)
}

func newForm() (*huh.Form, *FormData) {
	data := &FormData{}

	var form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Whats the name of your go module").
				Placeholder("github.com/foo/bar").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("go module path cannot be empty")
					}
					matched, err := regexp.MatchString(`^[a-z0-9-_.]+(?:/[a-z0-9-_.]+)*$`, s)
					if err != nil {
						return fmt.Errorf("error validating go module path: %v", err)
					}
					if !matched {
						return fmt.Errorf("go module path can only contain lowercase letters, numbers, and slashes, and must not start or end with a slash")
					}
					return nil
				}).
				Value(&data.modulePath),
			huh.NewMultiSelect[string]().
				Title("Which optional features do you want to use?").
				Options(features...).
				Value(&data.features),
		),
		huh.NewGroup(
			huh.NewNote().
				DescriptionFunc(func() string { return fmt.Sprintf("\n\nModule path: %s\nFeatures: %s", data.modulePath, data.features) }, data),
			huh.NewConfirm().
				Title("Do you want to create a new project with these options?").
				Value(&data.ok),
		),
	)

	return form, data
}

func extractDirPath(modulePath string) string {
	if strings.Contains(modulePath, "/") {
		parts := strings.Split(modulePath, "/")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
		panic("Invalid module path")
	}
	return modulePath
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new project",
	Long:  `This initializes a new flamigo project with the given name`,
	Run: func(cmd *cobra.Command, args []string) {

		form, data := newForm()
		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if !data.ok {
			fmt.Println("Aborted by user")
			return
		}

		folderPath := extractDirPath(data.modulePath)
		project.InitializeDirectories(folderPath)
		project.InitializeGoMod(folderPath, data.modulePath)
		err = project.InitializeProjectFiles(folderPath, TemplateData{
			ProjectModulePath: data.modulePath,
			ProjectName:       folderPath,
			Features:          data.features,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Propagated templates to %s\n", folderPath)
		err = project.TidyGoMod(folderPath)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
