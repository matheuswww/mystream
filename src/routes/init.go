package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	admin_routes "github.com/matheuswww/mystream/src/routes/admin"
	file_routes "github.com/matheuswww/mystream/src/routes/file"
	upload_routes "github.com/matheuswww/mystream/src/routes/upload"
	user_routes "github.com/matheuswww/mystream/src/routes/user"
)

func InitRoutes(r *gin.Engine, sql *sql.DB) {
	upload_routes.InitUploadRoutes(r, sql)
	file_routes.InitFileRoutes(r)
	user_routes.InitUserRoutes(r, sql)
	admin_routes.InitAdminRouter(r, sql)
}
