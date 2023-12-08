
		package packagename

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
			//TODO implement me
			panic("implement me")
		}

		func (r *read) FindById(ctx context.Context, p Entity) (*Entity, error) {
			var e Entity
			//TODO implement me
			panic("implement me")
		}

		func (r *read) FindAllPagination(ctx context.Context, limit, offset int16) ([]*Entity, error) {
			var es []*Entity
			//TODO implement me
			panic("implement me")
		}
		