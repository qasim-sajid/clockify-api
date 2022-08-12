package main

import (
	"fmt"

	"github.com/qasim-sajid/clockify-api/dbhandler"
	"github.com/qasim-sajid/clockify-api/models"
)

// Main function
func main() {
	dbC, err := dbhandler.NewDBClient("DB")
	if err != nil {
		panic(err)
	}

	// task := models.Task{}
	// task.Billable = true
	// task.Description = "Description2"
	// task.StartTime = time.Now()
	// task.EndTime = time.Now()
	// task.Date = time.Now()
	// task.IsActive = false
	// task.Project = &models.Project{ID: "p_99544cf2-bba4-4274-82a3-7f0707a5404c"}
	// tags := make([]*models.Tag, 0)
	// tags = append(tags, &models.Tag{ID: "t_b51b6a64-5456-4ea5-8722-52d6515491ab"})
	// tags = append(tags, &models.Tag{ID: "t_8ccf4a82-53cf-49d4-8c87-dbeed1859a16"})
	// task.Tags = tags
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
	// project.Name = "Project4"
	// project.ColorTag = "Silver"
	// project.IsPublic = false
	// project.TrackedHours = 10
	// project.ProgressPercentage = 1
	// project.Client = &models.Client{ID: "c_dd736ca0-e21d-4d5e-b8b1-f73987c7f0d6"}
	// project.Workspace = &models.Workspace{ID: "w_3ef2b6ed-6b36-43f9-825f-b6396cb4bd56"}
	// teamMembers := make([]*models.TeamMember, 0)
	// teamMembers = append(teamMembers, &models.TeamMember{ID: "tm_cb605827-789b-41d6-96cf-7af79f4ad810"})
	// project.TeamMembers = teamMembers
	// teamGroups := make([]*models.TeamGroup, 0)
	// teamGroups = append(teamGroups, &models.TeamGroup{ID: "tg_b861b2aa-3a05-4423-91d1-fec0d24949eb"})
	// project.TeamGroups = teamGroups
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

	// project := models.Project{ID: "p_c337bb5a-b0a0-4f6a-b7a1-f1056e2c02ea"}
	// updates := make(map[string]interface{})
	// updates["tracked_hours"] = 77
	// updates["progress_percentage"] = 100
	// // teamMembers := make([]*models.TeamMember, 0)
	// // teamMembers = append(teamMembers, &models.TeamMember{ID: "tm_b44ee6f5-250b-487e-b89f-3af2c1873352"})
	// // teamMembers = append(teamMembers, &models.TeamMember{ID: "tm_a0a5c6c4-9610-40ac-a254-0c75c77b0bd6"})
	// // updates["project_team_members"] = teamMembers
	// _, err = dbC.UpdateProject(project.ID, updates)
	// if err != nil {
	// 	panic(err)
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

	// teamGroup := models.TeamGroup{ID: "tg_b861b2aa-3a05-4423-91d1-fec0d24949eb"}
	// updates := make(map[string]interface{})
	// teamMembers := make([]*models.TeamMember, 0)
	// teamMembers = append(teamMembers, &models.TeamMember{ID: "tm_b44ee6f5-250b-487e-b89f-3af2c1873352"})
	// teamMembers = append(teamMembers, &models.TeamMember{ID: "tm_a0a5c6c4-9610-40ac-a254-0c75c77b0bd6"})
	// updates["team_group_team_members"] = teamMembers
	// _, err = dbC.UpdateTeamGroup(teamGroup.ID, updates)
	// if err != nil {
	// 	panic(err)
	// }

	// teamMember := models.TeamMember{}
	// teamMember.BillableRate = 50
	// teamMember.User = &models.User{Email: "user4@gmail.com"}
	// teamMember.Workspace = &models.Workspace{ID: "w_3ef2b6ed-6b36-43f9-825f-b6396cb4bd56"}
	// teamMember.TeamRole = &models.TeamRole{ID: "tr_f46cb68a-5587-47a6-ac04-48c9d0c7557c"}
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

	// teamMember := models.TeamMember{ID: "tm_a0a5c6c4-9610-40ac-a254-0c75c77b0bd6"}
	// updates := make(map[string]interface{})
	// teamGroups := make([]*models.TeamGroup, 0)
	// teamGroups = append(teamGroups, &models.TeamGroup{ID: "tg_b861b2aa-3a05-4423-91d1-fec0d24949eb"})
	// updates["team_member_team_groups"] = teamGroups
	// _, err = dbC.UpdateTeamMember(teamMember.ID, updates)
	// if err != nil {
	// 	panic(err)
	// }

	teamRole := models.TeamRole{}
	teamRole.Role = "TeamRole3"
	_, status, err := dbC.AddTeamRole(&teamRole)
	if err != nil {
		panic(err)
	}

	// teamRoles, err := dbC.GetAllTeamRoles()
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < len(teamRoles); i++ {
	// 	fmt.Println(teamRoles[i])
	// }

	// user := models.User{}
	// user.Name = "User4"
	// user.Email = "user4@gmail.com"
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

	fmt.Println("Status: ", status)
}
