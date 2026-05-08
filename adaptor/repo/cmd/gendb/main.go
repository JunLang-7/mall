package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const defaultQueryPath = "./query"

type dbConfig struct {
	DSN               string   `yaml:"dsn"`
	DB                string   `yaml:"db"`
	Tables            []string `yaml:"tables"`
	OnlyModel         bool     `yaml:"onlyModel"`
	OutPath           string   `yaml:"outPath"`
	OutFile           string   `yaml:"outFile"`
	WithUnitTest      bool     `yaml:"withUnitTest"`
	ModelPkgName      string   `yaml:"modelPkgName"`
	FieldNullable     bool     `yaml:"fieldNullable"`
	FieldCoverable    bool     `yaml:"fieldCoverable"`
	FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"`
	FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`
	FieldSignable     bool     `yaml:"fieldSignable"`
}

type yamlConfig struct {
	Version  string    `yaml:"version"`
	Database *dbConfig `yaml:"database"`
}

func main() {
	configPath := flag.String("c", "gen.yaml", "path to gorm/gen yaml config")
	flag.Parse()

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatalln("parse config fail:", err)
	}
	config.revise()

	db, err := connectDB(config)
	if err != nil {
		log.Fatalln("connect db server fail:", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           config.OutPath,
		OutFile:           config.OutFile,
		ModelPkgPath:      config.ModelPkgName,
		WithUnitTest:      config.WithUnitTest,
		FieldNullable:     config.FieldNullable,
		FieldCoverable:    config.FieldCoverable,
		FieldWithIndexTag: config.FieldWithIndexTag,
		FieldWithTypeTag:  config.FieldWithTypeTag,
		FieldSignable:     config.FieldSignable,
	})
	g.UseDB(db)

	models, err := genModels(g, db, config.Tables)
	if err != nil {
		log.Fatalln("get tables info fail:", err)
	}

	if !config.OnlyModel {
		g.ApplyBasic(models...)
	}
	g.Execute()
}

func parseConfig(path string) (*dbConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg yamlConfig
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.Database == nil {
		return nil, fmt.Errorf("missing database config")
	}
	return cfg.Database, nil
}

func (c *dbConfig) revise() {
	if c.DB == "" {
		c.DB = "mysql"
	}
	if c.OutPath == "" {
		c.OutPath = defaultQueryPath
	}

	tables := make([]string, 0, len(c.Tables))
	for _, table := range c.Tables {
		if table = strings.TrimSpace(table); table != "" {
			tables = append(tables, table)
		}
	}
	c.Tables = tables
}

func connectDB(c *dbConfig) (*gorm.DB, error) {
	if c.DSN == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}
	if c.DB != "mysql" {
		return nil, fmt.Errorf("unsupported db %q", c.DB)
	}
	return gorm.Open(mysql.Open(c.DSN))
}

func genModels(g *gen.Generator, db *gorm.DB, tables []string) ([]interface{}, error) {
	if len(tables) == 0 {
		var err error
		tables, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}

	models := make([]interface{}, len(tables))
	for i, tableName := range tables {
		models[i] = g.GenerateModel(tableName, gen.FieldModify(escapeGORMCommentTag))
	}
	return models, nil
}

func escapeGORMCommentTag(f gen.Field) gen.Field {
	values, ok := f.GORMTag[field.TagKeyGormComment]
	if !ok {
		return f
	}

	for i, value := range values {
		// GORM comments come from DB metadata and may contain JSON examples.
		values[i] = strings.ReplaceAll(value, `"`, `\"`)
	}
	f.GORMTag.Set(field.TagKeyGormComment, values...)
	return f
}
