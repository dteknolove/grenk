
		package packagename

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
			//TODO implement me
			panic("implement me")
		}

		func (w *write) Update(ctx context.Context, p Update) error {
			//TODO implement me
			panic("implement me")
		}

		func (w *write) Delete(ctx context.Context, p Delete) error {
			//TODO implement me
			panic("implement me")
		}
		   