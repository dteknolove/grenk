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

	"github.com/dteknolove/grenk/pkg/db"
	"github.com/dteknolove/grenk/pkg/vip"
)

type ColumnInfo struct {
	ColumnName          string
	DataType            string
	GoType              string
	TitleCaseColumnName string
	IsLastField         bool
}

const (
	insertTemplate   = "`insert into table_name	(item,item2) values($1,$2)`"
	deleteTemplate   = "`delete from table_name where id = $1`"
	updateTemplate   = "`update from table_name	set item = $2	where id = $1`"
	countRowTemplate = "`select count(1) from table_name`"
	findIDTemplate   = "`select id from table_name where id = $1`"
	findAllTemplate  = "`select id,name from table_name tn order by updated_at desc limit $1 offset $2`"
)

func templateContent(repoName string) string {
	return fmt.Sprintf(`package %s

import 	"github.com/gofrs/uuid/v5"

type Entity struct {
{{- range . }}
	{{ .TitleCaseColumnName }} {{ .GoType }}{{ if not .IsLastField }} {{ end }}
{{- end }}
}
type Insert struct {}
type Update struct {}
type Delete struct {}
`, repoName)
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
	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

func main() {
	app := &cli.App{
		Name:  "grenk",
		Usage: "fight the loneliness!",
		Commands: []*cli.Command{
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
						Name:  "table",
						Value: "table_name",
						Usage: "your table database name",
					},
					&cli.StringFlag{
						Name:  "reponame",
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

					tableName := cctx.String("table")
					repoName := cctx.String("reponame")

					flag.Parse()

					q := `SELECT column_name, data_type FROM information_schema.columns WHERE table_name=$1 AND table_schema=$2`
					rows, err := s.DB.Query(ctx, q, tableName, vipp.DbSchema)
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

					entityTempalte := templateContent(repoName)

					folderPath := vipp.RepoPath + "/" + repoName
					err = os.MkdirAll(folderPath, 0755)

					file, err := os.Create(folderPath + "/entity.go")
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

					interfaceContent := fmt.Sprintf(`
		package %s

		import "context"

		type Read interface {
			CountRow(ctx context.Context) (int32, error)
			FindById(ctx context.Context, p Entity) (*Entity, error)
			FindAllPagination(ctx context.Context, limit, offset int16) ([]*Entity, error)
		}
		type Write interface {
			Create(ctx context.Context, p Insert) error
			Update(ctx context.Context, p Update) error
			Delete(ctx context.Context, p Delete) error
		}

		`,
						repoName)
					writeContent := fmt.Sprintf(`
		package %s

		import (
			"context"
			"github.com/jackc/pgx/v5"
		)

		type write struct {
			TX pgx.Tx
		}

		func NewWrite(tx pgx.Tx) Write {
			return &write{TX: tx}
		}

		func (w *write) Create(ctx context.Context, p Insert) error {
		// TODO: create
				
		// q := %s
		// _, err := w.TX.Exec(ctx, q, p.Name)
		// if err != nil {
			// return err
		// }
		// return nil
	
		panic("implement create")
		}

		func (w *write) Update(ctx context.Context, p Update) error {
		// TODO: update

		// q := %s
		// _, err := w.TX.Exec(ctx, q, p.Name)
		// 	if err != nil {
		//		return err
		// 	}
		// return nil

		panic("implement update")
		}

		func (w *write) Delete(ctx context.Context, p Delete) error {
		// TODO: delete

		//	q := %s
		//	exec, err := w.TX.Exec(ctx, q, p.Item)
		//	if err != nil {
		//		if exec.RowsAffected() == 0 {
		//			return db.ErrNoAffected
		//		}
		//		return err
		//	}
		//	return nil

		panic("implement delete")
		}
		   `, repoName, insertTemplate, updateTemplate, deleteTemplate)
					readContent := fmt.Sprintf(`
		package %s

		import (
			"context"
			"github.com/jackc/pgx/v5/pgxpool"
		)

		type read struct {
			DB *pgxpool.Pool
		}

		func NewRead(db *pgxpool.Pool) Read {
			return &read{DB: db}
		}

		func (r *read) CountRow(ctx context.Context) (int32, error) {
		
	// TODO: countrow

	// var count int32
	// q := %s
	// err := r.DB.QueryRow(ctx, q).Scan(&count)
	//if err != nil {
	// 	return 0, err
	// }
	// return count, nil
	panic("implement me")
		}

		func (r *read) FindById(ctx context.Context, p Entity) (*Entity, error) {
		
	// TODO: find ID

	// var e Entity
	// q := %s
	// err := r.DB.QueryRow(ctx, q, p.ID).Scan(&e)
	// if err != nil {
	// 	return nil, err
	// }
	// return &e, nil

	panic("implement me")
		}

		func (r *read) FindAllPagination(ctx context.Context, limit, offset int16) ([]*Entity, error) {
		
// TODO: find all

	// var es []*Entity
	// q := %s
	// rows, errs := r.DB.Query(ctx, q, limit, offset)
	// if errs != nil {
	// 	return nil, errs
	// }
	//
	// for rows.Next() {
	// 	var e Entity
	// 	err := rows.Scan(&e)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	es = append(es, &e)
	// }
	// return es, nil
			panic("implemented find all")
		}
		`, repoName, countRowTemplate, findIDTemplate, findAllTemplate)

					interfaceFilePath := folderPath + "/interface.go"
					err = createFile(interfaceFilePath, interfaceContent)
					if err != nil {
						log.Println("Failed to create interface.go file:", err)
						return err
					}

					writeFilePath := folderPath + "/write.go"
					err = createFile(writeFilePath, writeContent)
					if err != nil {
						log.Println("Failed to create write.go file:", err)
						return err
					}

					readFilePath := folderPath + "/read.go"
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
				"grenk generate --table table_name --reponame packageName \n" +
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
