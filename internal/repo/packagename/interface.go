
		package packagename

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

		