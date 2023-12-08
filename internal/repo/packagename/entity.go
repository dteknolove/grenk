package packagename

import 	"github.com/gofrs/uuid/v5"

type Entity struct {
	Id uuid.UUID 
	Code string 
	Title string 
	CreatedAt time.Time 
	UpdatedAt time.Time
}
type Insert struct {}
type Update struct {}
type Delete struct {}
