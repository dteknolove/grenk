package ptrn

import (
	"fmt"
)

const (
	GRENK_YAML = `database:
  name: "database_name"
  password: "database_password"
  username: "database_user"
  port: 5432 
  host: "database_host"
  schema: "public" # choose your own database schema
  repo_path: "path/to/your/model"
`
	SQL_INSERT     = "`insert into table_name	(item,item2) values($1,$2)`"
	SQL_DELETE     = "`delete from table_name where id = $1`"
	SQL_UPDATE     = "`update table_name set item = $2	where id = $1`"
	SQL_COUNT_ROW  = "`select count(1) from table_name`"
	SQL_FIND_BY_ID = "`select id from table_name where id = $1`"
	SQL_FIND_ALL   = "`select id,name from table_name tn order by updated_at desc limit $1 offset $2`"

	CONNECT_TABLE_GET_COLUMN = `SELECT column_name, data_type FROM information_schema.columns WHERE table_name=$1 AND table_schema=$2`

	FLAG_TABLE   = "table"
	FLAG_PACKAGE = "package"

	PATH_INTERFACE = "/interface.go"
	PATH_ENTITY    = "/entity.go"
	PATH_WRITE     = "/write.go"
	PATH_READ      = "/read.go"
)

func TemplateEntityContent(repoName string) string {
	return fmt.Sprintf(`package %s

type Entity struct {
{{- range . }}
	{{ .TitleCaseColumnName }} {{ .GoType }}{{ if not .IsLastField }} {{ end }}
{{- end }}
}
type Insert struct {}
type Update struct {}
type Delete struct {}
type Search struct {}
`, repoName)
}

func InterfaceContent(flagPackage string) string {
	return fmt.Sprintf(`package %s

import "context"

type Read interface {
	CountRow(ctx context.Context) (int32, error)
	FindById(ctx context.Context, p Entity) (*Entity, error)
	FindAllPagination(ctx context.Context, limit, offset int16,s Search) ([]*Entity, error)
}
type Write interface {
	Create(ctx context.Context, p Insert) error
	Update(ctx context.Context, p Update) error
	Delete(ctx context.Context, p Delete) error
}
`,
		flagPackage)
}

func WriteContent(flagPackage string) string {
	return fmt.Sprintf(`package %s

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
//		return err
//	}
//		if exec.RowsAffected() == 0 {
//			return db.ErrNoAffected
//		}
//	return nil

panic("implement delete")
}
		   `, flagPackage, SQL_INSERT, SQL_UPDATE, SQL_DELETE)
}

func ReadContent(flagPackage string) string {
	return fmt.Sprintf(`package %s

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

//  var count int32
//  q := %s
//  err := r.DB.QueryRow(ctx, q).Scan(&count)
// if err != nil {
//  	return 0, err
//  }
//  return count, nil
panic("implement me")
}

func (r *read) FindById(ctx context.Context, p Entity) (*Entity, error) {

// TODO: find ID
//
// var e Entity
// q := %s
// err := r.DB.QueryRow(ctx, q, p.ID).Scan(&e)
// if err != nil {
//  	return nil, err
// }
// return &e, nil
//
panic("implement find ID")
}

func (r *read) FindAllPagination(ctx context.Context, limit, offset int16, s Search) ([]*Entity, error) {

// TODO: find all
//
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
		`, flagPackage, SQL_COUNT_ROW, SQL_FIND_BY_ID, SQL_FIND_ALL)
}
