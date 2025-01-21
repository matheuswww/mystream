package routes

import (
	"database/sql"

	"github.com/matheuswww/mystream/src/router"
	file_routes "github.com/matheuswww/mystream/src/routes/file"
	upload_routes "github.com/matheuswww/mystream/src/routes/upload"
)

func InitRoutes(r *router.Router, db *sql.DB) {
	upload_routes.InitUploadRoutes(r, db)
	file_routes.InitFileRoutes(r)
}
