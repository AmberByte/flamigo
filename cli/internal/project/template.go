package project

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed project/**/*
var tmplFS embed.FS

func InitializeProjectFiles(projectRoot string, data any) error {

	return copyTemplateFiles(projectRoot, "/", data)
}

func copyTemplateFiles(projectRoot string, relativeDir string, data any) error {
	fmt.Printf("Processing directory: %s\n", path.Join("project", relativeDir))
	templateFSPath := path.Join("project", relativeDir)
	templateFiles, err := tmplFS.ReadDir(templateFSPath)
	if err != nil {
		return fmt.Errorf("failed to read directory (%s): %w", templateFSPath, err)
	}

	// Check for conditions file to determine if the directory should be skipped
	if conditionsFile, err := tmplFS.ReadFile(path.Join(templateFSPath, "conditions")); err == nil {
		conditionResult, err := evaluateTemplate(string(conditionsFile), data)
		if err != nil {
			return fmt.Errorf("failed to evaluate conditions file (%s): %w", templateFSPath, err)
		}
		if conditionResult != "true" {
			fmt.Printf("Skipping directory: %s\n", templateFSPath)
			return nil
		}
	}

	for _, file := range templateFiles {
		templateFSFilePath := path.Join(templateFSPath, file.Name())
		fmt.Printf("Processing child: %s, isDir: %t\n", templateFSFilePath, file.IsDir())
		if file.IsDir() {
			err := copyTemplateFiles(projectRoot, path.Join(relativeDir, file.Name()), data)
			if err != nil {
				return err
			}
			continue
		}

		// If file is not a template file skip it
		if !strings.HasSuffix(file.Name(), ".tmpl") {
			continue
		}

		// Check for conditions file to determine if the directory should be skipped
		if conditionsFile, err := tmplFS.ReadFile(path.Join(templateFSPath, file.Name()+".conditions")); err == nil {
			conditionResult, err := evaluateTemplate(string(conditionsFile), data)
			if err != nil {
				return fmt.Errorf("failed to evaluate conditions file (%s): %w", templateFSPath, err)
			}
			if conditionResult != "true" {
				fmt.Printf("Skipping file: %s\n", templateFSPath)
				continue
			}
		}

		fileContent, err := tmplFS.ReadFile(templateFSFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file (%s): %w", templateFSFilePath, err)
		}

		renderedContent, err := evaluateTemplate(string(fileContent), data)
		if err != nil {
			return fmt.Errorf("failed to evaluate template (%s): %w", templateFSFilePath, err)
		}

		outputFilePath := strings.ReplaceAll(path.Join(projectRoot, relativeDir, file.Name()), ".tmpl", "")
		fmt.Printf("Creating file: %s\n", outputFilePath)

		err = createFileIfNotExists(outputFilePath, renderedContent)
		if err != nil {
			return fmt.Errorf("failed to write to file (%s): %w", outputFilePath, err)
		}
	}

	return nil
}

func evaluateTemplate(fileContent string, data any) (string, error) {
	tmpl, err := template.New("main").
		Funcs(template.FuncMap{}).Parse(string(fileContent))
	if err != nil {
		return "", fmt.Errorf("parse template file(%s): %w", string(fileContent), err)
	}
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return result.String(), nil
}

func createFileIfNotExists(filePath string, content string) error {
	dir := path.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory (%s): %w", dir, err)
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file (%s): %w", filePath, err)
		}
		defer file.Close()

		_, err = file.WriteString(content)
		if err != nil {
			return fmt.Errorf("failed to write content to file (%s): %w", filePath, err)
		}
	}
	return nil
}
