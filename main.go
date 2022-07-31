package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/auth"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/handler"
)

// Main function
func main() {
	conf.InitConfigs()

	apiHandler, err := handler.NewHandler()
	if err != nil {
		panic(err)
	}

	router := setupRouter(apiHandler)

	err = router.Run(conf.GetServerAddress())
	if err != nil {
		panic(err)
	}
}

func setupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// healthz
	router.GET("/healthz", healthGET())

	router.POST("/signup", handler.SignUpUser(h))
	router.POST("/login", auth.LoginUser(h))
	router.POST("/refresh_user", auth.RefreshUserTokenPOST(h))

	router.POST("/client", auth.IsUserAuthorized(handler.AddClient, h))
	router.GET("/clients", auth.IsUserAuthorized(handler.GetAllClients, h))
	router.GET("/clients/:client_id", auth.IsUserAuthorized(handler.GetClient, h))
	router.PUT("/clients/:client_id", auth.IsUserAuthorized(handler.UpdateClient, h))
	router.DELETE("/clients/:client_id", auth.IsUserAuthorized(handler.DeleteClient, h))

	router.POST("/project", auth.IsUserAuthorized(handler.AddProject, h))
	router.GET("/projects", auth.IsUserAuthorized(handler.GetAllProjects, h))
	router.GET("/projects/:project_id", auth.IsUserAuthorized(handler.GetProject, h))
	router.PUT("/projects/:project_id", auth.IsUserAuthorized(handler.UpdateProject, h))
	router.DELETE("/projects/:project_id", auth.IsUserAuthorized(handler.DeleteProject, h))

	router.POST("/tag", auth.IsUserAuthorized(handler.AddTag, h))
	router.GET("/tags", auth.IsUserAuthorized(handler.GetAllTags, h))
	router.GET("/tags/:tag_id", auth.IsUserAuthorized(handler.GetTag, h))
	router.PUT("/tags/:tag_id", auth.IsUserAuthorized(handler.UpdateTag, h))
	router.DELETE("/tags/:tag_id", auth.IsUserAuthorized(handler.DeleteTag, h))

	router.POST("/task", auth.IsUserAuthorized(handler.AddTask, h))
	router.GET("/tasks", auth.IsUserAuthorized(handler.GetAllTasks, h))
	router.GET("/tasks/:task_id", auth.IsUserAuthorized(handler.GetTask, h))
	router.PUT("/tasks/:task_id", auth.IsUserAuthorized(handler.UpdateTask, h))
	router.DELETE("/tasks/:task_id", auth.IsUserAuthorized(handler.DeleteTask, h))

	router.POST("/team_group", auth.IsUserAuthorized(handler.AddTeamGroup, h))
	router.GET("/team_groups", auth.IsUserAuthorized(handler.GetAllTeamGroups, h))
	router.GET("/team_groups/:team_group_id", auth.IsUserAuthorized(handler.GetTeamGroup, h))
	router.PUT("/team_groups/:team_group_id", auth.IsUserAuthorized(handler.UpdateTeamGroup, h))
	router.DELETE("/team_groups/:team_group_id", auth.IsUserAuthorized(handler.DeleteTeamGroup, h))

	router.POST("/team_member", auth.IsUserAuthorized(handler.AddTeamMember, h))
	router.GET("/team_members", auth.IsUserAuthorized(handler.GetAllTeamMembers, h))
	router.GET("/team_members/:team_member_id", auth.IsUserAuthorized(handler.GetTeamMember, h))
	router.PUT("/team_members/:team_member_id", auth.IsUserAuthorized(handler.UpdateTeamMember, h))
	router.DELETE("/team_members/:team_member_id", auth.IsUserAuthorized(handler.DeleteTeamMember, h))

	router.POST("/team_role", auth.IsUserAuthorized(handler.AddTeamRole, h))
	router.GET("/team_roles", auth.IsUserAuthorized(handler.GetAllTeamRoles, h))
	router.GET("/team_roles/:team_role_id", auth.IsUserAuthorized(handler.GetTeamRole, h))
	router.PUT("/team_roles/:team_role_id", auth.IsUserAuthorized(handler.UpdateTeamRole, h))
	router.DELETE("/team_roles/:team_role_id", auth.IsUserAuthorized(handler.DeleteTeamRole, h))

	router.GET("/users", auth.IsUserAuthorized(handler.GetAllUsers, h))
	router.GET("/users/:user_id", auth.IsUserAuthorized(handler.GetUser, h))
	router.PUT("/users/:user_id", auth.IsUserAuthorized(handler.UpdateUser, h))
	router.DELETE("/users/:user_id", auth.IsUserAuthorized(handler.DeleteUser, h))

	router.POST("/workspace", auth.IsUserAuthorized(handler.AddWorkspace, h))
	router.GET("/workspaces", auth.IsUserAuthorized(handler.GetAllWorkspaces, h))
	router.GET("/workspaces/:workspace_id", auth.IsUserAuthorized(handler.GetWorkspace, h))
	router.PUT("/workspaces/:workspace_id", auth.IsUserAuthorized(handler.UpdateWorkspace, h))
	router.DELETE("/workspaces/:workspace_id", auth.IsUserAuthorized(handler.DeleteWorkspace, h))

	return router
}

func healthGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "clockify-api")
	}
}
