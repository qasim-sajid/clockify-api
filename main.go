package main

import (
	"github.com/qasim-sajid/clockify-api/dbhandler"
)

// Main function
func main() {
	dbC, err := dbhandler.NewDBClient("DB")
	if err != nil {
		panic(err)
	}

	// task := models.Task{}
	// task.ID = "ID8"
	// task.Billable = true
	// task.Description = "Description8"
	// task.StartTime = time.Now()
	// task.EndTime = time.Now()
	// task.Date = time.Now()
	// task.Project = &models.Project{ID: "Project2"}
	// _, status, err := dbC.AddTask(&task)
	// if err != nil {
	// 	panic(err)
	// }

	// tasks, err := dbC.GetAllTasks()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(tasks); i++ {
	// 	fmt.Println(tasks[i])
	// }

	// task := models.Task{ID: "ID5"}
	// updates := make(map[string]interface{})
	// updates["project_id"] = "p_99544cf2-bba4-4274-82a3-7f0707a5404c"
	// _, err = dbC.UpdateTask(task.ID, updates)
	// if err != nil {
	// 	panic(err)
	// }

	// client := models.Client{}
	// client.Name = "Client2"
	// client.Address = "Client 2 Address"
	// client.Note = "Client 2 Note"
	// client.IsArchived = true
	// _, status, err := dbC.AddClient(&client)
	// if err != nil {
	// 	panic(err)
	// }

	// clients, err := dbC.GetAllClients()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(clients); i++ {
	// 	fmt.Println(clients[i])
	// }

	// project := models.Project{}
	// project.Name = "Project1"
	// project.ColorTag = "Blue"
	// project.IsPublic = true
	// project.TrackedHours = 0
	// project.ProgressPercentage = 0
	// project.Workspace = &models.Workspace{ID: "w_3ef2b6ed-6b36-43f9-825f-b6396cb4bd56"}
	// _, status, err := dbC.AddProject(&project)
	// if err != nil {
	// 	panic(err)
	// }

	// projects, err := dbC.GetAllProjects()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(projects); i++ {
	// 	fmt.Println(projects[i])
	// }

	// workspace := models.Workspace{}
	// workspace.Name = "Workspace1"
	// _, status, err := dbC.AddWorkspace(&workspace)
	// if err != nil {
	// 	panic(err)
	// }

	// workspace, err := dbC.GetAllWorkspaces()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(workspace); i++ {
	// 	fmt.Println(workspace[i])
	// }

	// tag := models.Tag{}
	// tag.Name = "Tag3"
	// _, status, err := dbC.AddTag(&tag)
	// if err != nil {
	// 	panic(err)
	// }

	// tag, err := dbC.GetAllTags()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(tag); i++ {
	// 	fmt.Println(tag[i])
	// }

	// teamGroup := models.TeamGroup{}
	// teamGroup.Name = "TeamGroup2"
	// teamGroup.Workspace = &models.Workspace{ID: "w_3ef2b6ed-6b36-43f9-825f-b6396cb4bd56"}
	// _, status, err := dbC.AddTeamGroup(&teamGroup)
	// if err != nil {
	// 	panic(err)
	// }

	// teamGroup, err := dbC.GetAllTeamGroups()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(teamGroup); i++ {
	// 	fmt.Println(teamGroup[i])
	// }

	// teamMember := models.TeamMember{}
	// teamMember.BillableRate = 25
	// teamMember.User = &models.User{Email: "user3@gmail.com"}
	// teamMember.Workspace = &models.Workspace{ID: "w_3ef2b6ed-6b36-43f9-825f-b6396cb4bd56"}
	// teamMember.TeamRole = &models.TeamRole{ID: "tr_538112a5-dae6-4f84-b1b3-9b9d1dd1666b"}
	// _, status, err := dbC.AddTeamMember(&teamMember)
	// if err != nil {
	// 	panic(err)
	// }

	// teamMember, err := dbC.GetAllTeamMembers()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(teamMember); i++ {
	// 	fmt.Println(teamMember[i])
	// }

	// teamRole := models.TeamRole{}
	// teamRole.Role = "TeamRole3"
	// _, status, err := dbC.AddTeamRole(&teamRole)
	// if err != nil {
	// 	panic(err)
	// }

	// teamRoles, err := dbC.GetAllTeamRoles()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(teamRoles); i++ {
	// 	fmt.Println(teamRoles[i])
	// }

	// user := models.User{}
	// user.Name = "User3"
	// user.Email = "user3@gmail.com"
	// _, status, err := dbC.AddUser(&user)
	// if err != nil {
	// 	panic(err)
	// }

	// users, err := dbC.GetAllUsers()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(users); i++ {
	// 	fmt.Println(users[i])
	// }

	// fmt.Println("Status: ", status)
}
