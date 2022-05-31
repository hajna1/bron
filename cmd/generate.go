/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/hajna1/bron/internal/generator"
	"github.com/hajna1/bron/internal/sql/lexer"
	"github.com/hajna1/bron/internal/sql/merger"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// generateCmd represents the create command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates models and repository from your migration files",
	Long:  `generates models and repository from your migration files`,

	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.WarnLevel)
		p := lexer.New()
		migrationDir, err := cmd.Flags().GetString("migrations")
		if err != nil {
			return
		}
		outputDir, err := cmd.Flags().GetString("output")
		if err != nil {
			return
		}
		packageName, err := cmd.Flags().GetString("package")
		if err != nil {
			return
		}
		modelName, err := cmd.Flags().GetString("model")
		if err != nil {
			return
		}
		tableName, err := cmd.Flags().GetString("table")
		if err != nil {
			return
		}
		qs, err := p.ParseAllDirectory(migrationDir)
		if err != nil {
			fmt.Printf("err: %s", err)
			return
		}
		m := merger.New()
		tables := m.ParseAll(qs...)
		g, err := generator.New(packageName, outputDir, modelName)
		if err != nil {
			logrus.Fatal(err)
		}
		switch tableName {
		case "all":
			for _, table := range tables {
				if err := g.Generate(table); err != nil {
					logrus.Fatal(err)
				}
			}
		default:
			target, has := tables[tableName]
			if !has {
				logrus.WithField("table", tableName).Warnf("table not found")
				return
			}
			if err := g.Generate(target); err != nil {
				logrus.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("migrations", "m", "", "migration directory (required)")
	_ = generateCmd.MarkFlagRequired("migrations")
	generateCmd.Flags().StringP("output", "o", "", "project directory (required)")
	_ = generateCmd.MarkFlagRequired("output")
	generateCmd.Flags().StringP("package", "p", "", "package name (required)")
	_ = generateCmd.MarkFlagRequired("package")
	generateCmd.Flags().String("model", "model", "model package name")
	generateCmd.Flags().StringP("table", "t", "all", "table name")
}
