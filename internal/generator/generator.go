package generator

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hajna1/bron/internal/sql/models"
	"github.com/iancoleman/strcase"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"os"
	"text/template"
)

//go:embed templates/model/*
var modelContent embed.FS

//go:embed templates/repository/*
var repoContent embed.FS

func New(packageName, projectDirectory, modelDirectory string) (*Generator, error) {
	return &Generator{
		PackageName:      packageName,
		ProjectDirectory: projectDirectory,
		ModelDirectory:   modelDirectory,
	}, nil
}

type Generator struct {
	PackageName      string
	ProjectDirectory string
	ModelDirectory   string
}

func (g *Generator) Generate(m models.Create) error {
	t := Repository{
		PackageName:    g.PackageName,
		ConstName:      strcase.ToLowerCamel(fmt.Sprintf("%sTableName", m.TableName)),
		OriginName:     m.TableName,
		ModelType:      fmt.Sprintf("%s.%s", g.ModelDirectory, strcase.ToCamel(m.TableName)),
		ModelDirectory: g.ModelDirectory,
		GoName:         strcase.ToCamel(m.TableName),
		ID:             Column{},
		CreatedAt:      nil,
		UpdatedAt:      nil,
		Columns:        make([]Column, 0, len(m.CreateColumn)),
	}
	for _, c := range m.CreateColumn {
		goType, err := g.ToGoType(c.ColumnType)
		if err != nil {
			return err
		}
		c.ColumnType = goType

		column := Column{
			TableName:  strcase.ToCamel(m.TableName),
			OriginName: c.ColumnName,
			GoName:     strcase.ToCamel(c.ColumnName),
			ConstName:  strcase.ToLowerCamel(fmt.Sprintf("%sColumn%s", m.TableName, strcase.ToCamel(c.ColumnName))),
			GoType:     goType,
			IsNullable: !c.IsNotNull,
		}
		switch c.ColumnName {
		case "id":
			t.ID = column
		case "created_at":
			t.CreatedAt = &column
		case "updated_at":
			t.UpdatedAt = &column
		default:
			t.Columns = append(t.Columns, column)
		}
	}
	if err := g.GenerateModels(t); err != nil {
		return err
	}
	if err := g.GenerateRepository(t); err != nil {
		return err
	}
	return nil
}

func (g *Generator) GenerateModels(r Repository) error {
	dir := fmt.Sprintf("%s/internal/%s/%s", g.ProjectDirectory, r.OriginName, g.ModelDirectory)
	return g.generate(r, modelContent, dir, "model.tmpl", "templates/model/*.tmpl")
}

func (g *Generator) GenerateRepository(r Repository) error {
	dir := fmt.Sprintf("%s/internal/%s/repository", g.ProjectDirectory, r.OriginName)
	return g.generate(r, repoContent, dir, "repository.tmpl", "templates/repository/*.tmpl")
}

func (g *Generator) generate(r Repository, fs embed.FS, dir, repoTempl, pattern string) error {

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%s.go", dir, r.OriginName))
	if err != nil {
		return fmt.Errorf("cannot create file : %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	templ := template.Must(template.New(repoTempl).ParseFS(fs, pattern))
	var buf bytes.Buffer

	err = templ.Execute(&buf, r)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, "", buf.Bytes(), parser.ParseComments)

	if err != nil {
		return err
	}

	err = format.Node(f, fset, node)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) CopyColumnWithoutID(columns []Column) []Column {
	result := make([]Column, 0, len(columns))
	for _, c := range columns {
		if c.OriginName != "id" {
			result = append(result, c)
		}
	}
	return result
}

//LoadGo deprecated
func (g *Generator) LoadGo(patterns ...string) error {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedName,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return err
	}
	pkg := pkgs[0]
	fmt.Println(pkg.Name)
	return nil
}

func (g *Generator) ToGoType(sqlType string) (string, error) {
	m := map[string]string{
		"bigserial": "int64",
		"varchar":   "string",
		"timestamp": "time.Time",
		"real":      "float64",
		"boolean":   "bool",
	}
	s, has := m[sqlType]
	if !has {
		return "", fmt.Errorf("sql type not supported: %s", sqlType)
	}
	return s, nil
}
