package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/dteknolove/grenk/pkg/db"
	"github.com/dteknolove/grenk/pkg/ptrn"
	"github.com/dteknolove/grenk/pkg/vip"
)

type ColumnInfo struct {
	ColumnName          string
	DataType            string
	GoType              string
	TitleCaseColumnName string
	IsLastField         bool
}

func mapDataType(dataType string) string {
	switch {
	case strings.Contains(dataType, "uuid"):
		return "uuid.UUID"
	case strings.Contains(dataType, "integer"):
		return "int"
	case strings.Contains(dataType, "character varying"):
		return "string"
	case strings.Contains(dataType, "timestamp without time zone"):
		return "time.Time"
	case strings.Contains(dataType, "date"):
		return "time.Time"
	case strings.Contains(dataType, "bigint"):
		return "int64"
	case strings.Contains(dataType, "smallint"):
		return "int"
	case strings.Contains(dataType, "double precision"):
		return "float64"
	case strings.Contains(dataType, "text"):
		return "string"
	default:
		return dataType
	}
}

func addGoTypeAndTitleCase(columnInfoList []ColumnInfo) []ColumnInfo {
	for i, info := range columnInfoList {
		columnInfoList[i].GoType = mapDataType(info.DataType)
		columnInfoList[i].TitleCaseColumnName = toTitleCase(info.ColumnName)
		columnInfoList[i].IsLastField = i == len(columnInfoList)-1 // Set IsLastField
	}
	return columnInfoList
}

func toTitleCase(s string) string {
	caser := cases.Title(language.English)
	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	for i, word := range words {
		words[i] = caser.String(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

func main() {
	app := &cli.App{
		Name:  "grenk",
		Usage: "fight the loneliness!",
		Commands: []*cli.Command{
			{
				Name:        "init",
				Usage:       "generate grenk.yaml in root project",
				Description: "generate grenk.yaml in root project",
				Action: func(_ *cli.Context) error {
					file, err := os.Create("grenk.yaml")
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					tmpl, err := template.New("grenk.yaml").Parse(ptrn.GRENK_YAML)
					if err != nil {
						log.Fatal(err)
					}
					if err := tmpl.Execute(file, ptrn.GRENK_YAML); err != nil {
						log.Fatal(err)
					}

					fmt.Println("create grenk.yaml \n", ptrn.GRENK_YAML)

					return nil
				},
			},
			{
				Name:        "generate",
				Usage:       "generate database column to entity",
				Description: "save a ton of hour",
				Subcommands: []*cli.Command{
					{
						Name:  "handler",
						Usage: "generate handler and service",
						Action: func(_ *cli.Context) error {
							fmt.Println("this will generate handler")
							return nil
						},
					},
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  ptrn.FLAG_TABLE,
						Value: "table_name",
						Usage: "your table database name",
					},
					&cli.StringFlag{
						Name:  ptrn.FLAG_PACKAGE,
						Value: "packageName",
						Usage: "your package name",
					},
				},
				Action: func(cctx *cli.Context) error {
					ctx := context.Background()
					s := db.NewDbService(ctx)
					vipp, errVip := vip.New().App()
					if errVip != nil {
						return errVip
					}

					flagTable := cctx.String(ptrn.FLAG_TABLE)
					flagPackage := cctx.String(ptrn.FLAG_PACKAGE)

					flag.Parse()

					q := ptrn.CONNECT_TABLE_GET_COLUMN
					rows, err := s.DB.Query(ctx, q, flagTable, vipp.DbSchema)
					if err != nil {
						log.Fatal(err)
					}
					defer rows.Close()
					var columnInfoList []ColumnInfo

					for rows.Next() {
						var columnInfo ColumnInfo
						if err := rows.Scan(&columnInfo.ColumnName, &columnInfo.DataType); err != nil {
							log.Fatal(err)
						}
						columnInfoList = append(columnInfoList, columnInfo)
					}

					if err := rows.Err(); err != nil {
						log.Fatal(err)
					}

					columnInfoList = addGoTypeAndTitleCase(columnInfoList)

					for _, info := range columnInfoList {
						fmt.Printf("%s %s\n", info.TitleCaseColumnName, info.GoType)
					}

					entityTempalte := ptrn.TemplateEntityContent(flagPackage)

					folderPath := vipp.RepoPath + "/" + flagPackage
					err = os.MkdirAll(folderPath, 0755)

					file, err := os.Create(folderPath + ptrn.PATH_ENTITY)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					tmpl, err := template.New("").Parse(entityTempalte)
					if err != nil {
						log.Fatal(err)
					}

					if err := tmpl.Execute(file, columnInfoList); err != nil {
						log.Fatal(err)
					}

					fmt.Println("File created successfully.")

					interfaceContent := ptrn.InterfaceContent(flagPackage)
					writeContent := ptrn.WriteContent(flagPackage)
					readContent := ptrn.ReadContent(flagPackage)

					interfaceFilePath := folderPath + ptrn.PATH_INTERFACE
					err = createFile(interfaceFilePath, interfaceContent)
					if err != nil {
						log.Println("Failed to create interface.go file:", err)
						return err
					}

					writeFilePath := folderPath + ptrn.PATH_WRITE
					err = createFile(writeFilePath, writeContent)
					if err != nil {
						log.Println("Failed to create write.go file:", err)
						return err
					}

					readFilePath := folderPath + ptrn.PATH_READ
					err = createFile(readFilePath, readContent)
					if err != nil {
						log.Println("Failed to create write.go file:", err)
						return err
					}
					return nil
				},
			},
		},
		Action: func(*cli.Context) error {
			fmt.Println("Example: \n" +
				"grenk init \n" +
				"grenk generate --table table_name --package packageName \n" +
				" \n" +
				"----Copyright Teknolove 2024----")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
