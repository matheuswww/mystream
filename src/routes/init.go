package routes

import (
	"database/sql"

	"github.com/matheuswww/mystream/src/router"
	upload_routes "github.com/matheuswww/mystream/src/routes/upload"
)

func InitRoutes(r *router.Router, db *sql.DB) {
	upload_routes.InitUploadRoutes(r, db)
}