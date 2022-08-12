package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/handler"
)

// Main function
func main() {
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

	// healthz
	router.GET("/healthz", healthGET())

	router.POST("/client", handler.AddClient(h))
	router.GET("/clients", handler.GetAllClients(h))
	router.GET("/clients/:client_id", handler.GetClient(h))
	router.PUT("/clients/:client_id", handler.UpdateClient(h))
	router.DELETE("/clients/:client_id", handler.DeleteClient(h))

	router.POST("/project", handler.AddProject(h))
	router.GET("/projects", handler.GetAllProjects(h))
	router.GET("/projects/:project_id", handler.GetProject(h))
	router.PUT("/projects/:project_id", handler.UpdateProject(h))
	router.DELETE("/projects/:project_id", handler.DeleteProject(h))

	router.POST("/tag", handler.AddTag(h))
	router.GET("/tags", handler.GetAllTags(h))
	router.GET("/tags/:tag_id", handler.GetTag(h))
	router.PUT("/tags/:tag_id", handler.UpdateTag(h))
	router.DELETE("/tags/:tag_id", handler.DeleteTag(h))

	router.POST("/task", handler.AddTask(h))
	router.GET("/tasks", handler.GetAllTasks(h))
	router.GET("/tasks/:task_id", handler.GetTask(h))
	router.PUT("/tasks/:task_id", handler.UpdateTask(h))
	router.DELETE("/tasks/:task_id", handler.DeleteTask(h))

	router.POST("/team_group", handler.AddTeamGroup(h))
	router.GET("/team_groups", handler.GetAllTeamGroups(h))
	router.GET("/team_groups/:team_group_id", handler.GetTeamGroup(h))
	router.PUT("/team_groups/:team_group_id", handler.UpdateTeamGroup(h))
	router.DELETE("/team_groups/:team_group_id", handler.DeleteTeamGroup(h))

	router.POST("/team_member", handler.AddTeamMember(h))
	router.GET("/team_members", handler.GetAllTeamMembers(h))
	router.GET("/team_members/:team_member_id", handler.GetTeamMember(h))
	router.PUT("/team_members/:team_member_id", handler.UpdateTeamMember(h))
	router.DELETE("/team_members/:team_member_id", handler.DeleteTeamMember(h))

	router.POST("/team_role", handler.AddTeamRole(h))
	router.GET("/team_roles", handler.GetAllTeamRoles(h))
	router.GET("/team_roles/:team_role_id", handler.GetTeamRole(h))
	router.PUT("/team_roles/:team_role_id", handler.UpdateTeamRole(h))
	router.DELETE("/team_roles/:team_role_id", handler.DeleteTeamRole(h))

	router.POST("/user", handler.AddUser(h))
	router.GET("/users", handler.GetAllUsers(h))
	router.GET("/users/:user_id", handler.GetUser(h))
	router.PUT("/users/:user_id", handler.UpdateUser(h))
	router.DELETE("/users/:user_id", handler.DeleteUser(h))

	router.POST("/workspace", handler.AddWorkspace(h))
	router.GET("/workspaces", handler.GetAllWorkspaces(h))
	router.GET("/workspaces/:workspace_id", handler.GetWorkspace(h))
	router.PUT("/workspaces/:workspace_id", handler.UpdateWorkspace(h))
	router.DELETE("/workspaces/:workspace_id", handler.DeleteWorkspace(h))

	return router
}

func healthGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "clockify-api")
	}
}
