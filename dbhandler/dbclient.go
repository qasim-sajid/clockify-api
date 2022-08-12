package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/models"
	"github.com/qasim-sajid/clockify-api/queries"
)

type dbClient struct {
	dbName string
}

// NewDBClient returns ref to a new dbClient object
func NewDBClient(dbName string) (h DbHandler, err error) {
	client := &dbClient{
		dbName: dbName,
	}

	return client, nil
}

var dbConnection *sql.DB

// DB set up
func (db *dbClient) SetupDB() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conf.Configs.DBUser, conf.Configs.DBPassword, conf.Configs.DBName, conf.Configs.DBHost, conf.Configs.DBPort)
	dbC, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(fmt.Errorf("SetupDB: %v", err))
	}

	dbConnection = dbC

	initializeTablesIfNotExist()
}

func (db *dbClient) CloseDB() {
	if dbConnection != nil {
		dbConnection.Close()
	}
}

func initializeTablesIfNotExist() {
	_, _ = dbConnection.Exec(queries.CREATE_TABLES)
}

func (db *dbClient) RunInsertQuery(query string) (sql.Result, error) {
	result, err := dbConnection.Exec(query)

	if err != nil {
		return nil, fmt.Errorf("RunInsertQuery: %v", err)
	}

	return result, nil
}

func (db *dbClient) RunSelectQuery(query string) (*sql.Rows, error) {
	rows, err := dbConnection.Query(query)

	if err != nil {
		return nil, fmt.Errorf("RunSelectQuery: %v", err)
	}

	return rows, nil
}

func (db *dbClient) RunUpdateQuery(query string) (sql.Result, error) {
	result, err := dbConnection.Exec(query)

	if err != nil {
		return nil, fmt.Errorf("RunUpdateQuery: %v", err)
	}

	return result, nil
}

func (db *dbClient) RunDeleteQuery(query string) (sql.Result, error) {
	result, err := dbConnection.Exec(query)

	if err != nil {
		return nil, fmt.Errorf("RunDeleteQuery: %v", err)
	}

	return result, nil
}

func (db *dbClient) GetInsertQuery(structType interface{}) (string, error) {
	tableName, err := db.GetTableNameForStruct(structType)
	if err != nil {
		return ``, fmt.Errorf("GetInsertQuery: %v", err)
	}

	query := ``
	switch reflect.TypeOf(structType).Name() {
	case "Client":
		client := structType.(models.Client)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s', '%s', '%s', %t)`,
			tableName, db.GetColumnNamesForStruct(client), client.ID, client.Name, client.Address, client.Note, client.IsArchived)
	case "Project":
		project := structType.(models.Project)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s', '%s', %t, %f, %f, %f`,
			tableName, db.GetColumnNamesForStruct(project), project.ID, project.Name, project.ColorTag, project.IsPublic,
			project.TrackedHours, project.TrackedAmount, project.ProgressPercentage)
		if project.Client != "" {
			query = fmt.Sprintf(`%s, '%s'`, query, project.Client)
		} else {
			query = fmt.Sprintf(`%s, %v`, query, "null")
		}
		if project.Workspace != "" {
			query = fmt.Sprintf(`%s, '%s')`, query, project.Workspace)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "Tag":
		tag := structType.(models.Tag)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(tag), tag.ID, tag.Name)
	case "Task":
		task := structType.(models.Task)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s', %t, '%s', '%s', '%s', %t`,
			tableName, db.GetColumnNamesForStruct(task), task.ID, task.Description, task.Billable, task.StartTime.Format(conf.TIME_LAYOUT),
			task.EndTime.Format(conf.TIME_LAYOUT), task.Date.Format(conf.TIME_LAYOUT), task.IsActive)
		if task.Project != "" {
			query = fmt.Sprintf(`%s, '%s')`, query, task.Project)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "TeamGroup":
		teamGroup := structType.(models.TeamGroup)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s'`,
			tableName, db.GetColumnNamesForStruct(teamGroup), teamGroup.ID, teamGroup.Name)
		if teamGroup.Workspace != "" {
			query = fmt.Sprintf(`%s, '%s')`, query, teamGroup.Workspace)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "TeamMember":
		teamMember := structType.(models.TeamMember)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', %f`, tableName, db.GetColumnNamesForStruct(teamMember),
			teamMember.ID, teamMember.BillableRate)
		if teamMember.Workspace != "" {
			query = fmt.Sprintf(`%s, '%s'`, query, teamMember.Workspace)
		} else {
			query = fmt.Sprintf(`%s, %v`, query, "null")
		}
		if teamMember.User != "" {
			query = fmt.Sprintf(`%s, '%s'`, query, teamMember.User)
		} else {
			query = fmt.Sprintf(`%s, %v`, query, "null")
		}
		if teamMember.TeamRole != "" {
			query = fmt.Sprintf(`%s, '%s')`, query, teamMember.TeamRole)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "TeamRole":
		teamRole := structType.(models.TeamRole)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(teamRole), teamRole.ID, teamRole.Role)
	case "User":
		user := structType.(models.User)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s', '%s', '%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(user), user.ID, user.Name, user.Email, user.Username, user.Password)
	case "Workspace":
		workspace := structType.(models.Workspace)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(workspace), workspace.ID, workspace.Name)
	default:
		return ``, fmt.Errorf("GetInsertQuery: %v",
			errors.New("insert query generation error"))
	}

	return query, nil
}

func (db *dbClient) GetSelectQueryForStruct(structType interface{}, searchParams map[string]interface{}) (string, error) {
	if reflect.ValueOf(structType).Kind() == reflect.Struct {
		tableName, err := db.GetTableNameForStruct(structType)
		if err != nil {
			return ``, fmt.Errorf("GetSelectQueryForStruct: %v", err)
		}

		query := fmt.Sprintf("SELECT * FROM %s", tableName)

		i := 0
		for k, v := range searchParams {
			val := reflect.ValueOf(v)
			kind := val.Kind()
			if kind == reflect.Int || kind == reflect.Int64 {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Int())
				} else {
					query = fmt.Sprintf("%s && %s = %d", query, k, val.Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Float())
				} else {
					query = fmt.Sprintf("%s && %s = %f", query, k, val.Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Bool())
				} else {
					query = fmt.Sprintf("%s && %s = %t", query, k, val.Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = '%s'", query, k, val.String())
				} else {
					query = fmt.Sprintf("%s && %s = '%s'", query, k, val.String())
				}
			}

			i++
		}

		return query, nil
	}

	return ``, fmt.Errorf("GetSelectQueryForStruct: %v", errors.New("select query generation error"))
}

func (db *dbClient) GetUpdateQueryForStruct(structType interface{}, itemID string, updates map[string]interface{}) (string, error) {
	if reflect.ValueOf(structType).Kind() == reflect.Struct {
		tableName, err := db.GetTableNameForStruct(structType)
		if err != nil {
			return ``, fmt.Errorf("GetUpdateQueryForStruct: %v", err)
		}

		query := fmt.Sprintf("UPDATE %s", tableName)

		i := 0
		for k, v := range updates {
			kind := reflect.ValueOf(v).Kind()
			val := reflect.ValueOf(v)
			if kind == reflect.Int || kind == reflect.Int64 {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = %d", query, k, val.Int())
				} else {
					query = fmt.Sprintf("%s, %s = %d", query, k, val.Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = %f", query, k, val.Float())
				} else {
					query = fmt.Sprintf("%s, %s = %f", query, k, val.Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = %t", query, k, val.Bool())
				} else {
					query = fmt.Sprintf("%s, %s = %t", query, k, val.Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = '%s'", query, k, val.String())
				} else {
					query = fmt.Sprintf("%s, %s = '%s'", query, k, val.String())
				}
			}

			i++
		}

		query = fmt.Sprintf("%s WHERE _id = '%s'", query, itemID)
		return query, nil
	}

	return ``, fmt.Errorf("GetUpdateQueryForStruct: %v", errors.New("update query generation error"))
}

func (db *dbClient) GetDeleteQueryForStruct(structType interface{}, columnParams map[string]interface{}) (string, error) {
	if reflect.ValueOf(structType).Kind() == reflect.Struct {
		tableName, err := db.GetTableNameForStruct(structType)
		if err != nil {
			return ``, fmt.Errorf("GetSelectQueryForStruct: %v", err)
		}

		query := fmt.Sprintf("DELETE FROM %s", tableName)

		i := 0
		for k, v := range columnParams {
			kind := reflect.ValueOf(v).Kind()
			val := reflect.ValueOf(v)
			if kind == reflect.Int || kind == reflect.Int64 {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Int())
				} else {
					query = fmt.Sprintf("%s && %s = %d", query, k, val.Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Float())
				} else {
					query = fmt.Sprintf("%s && %s = %f", query, k, val.Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Bool())
				} else {
					query = fmt.Sprintf("%s && %s = %t", query, k, val.Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = '%s'", query, k, val.String())
				} else {
					query = fmt.Sprintf("%s && %s = '%s'", query, k, val.String())
				}
			}

			i++
		}

		return query, nil
	}

	return ``, fmt.Errorf("GetSelectQueryForStruct: %v", errors.New("select query generation error"))
}

func (db *dbClient) GetTableNameForStruct(t interface{}) (string, error) {
	switch reflect.TypeOf(t).Name() {
	case "Client":
		return "client", nil
	case "Project":
		return "project", nil
	case "Tag":
		return "tag", nil
	case "Task":
		return "task", nil
	case "TeamGroup":
		return "team_group", nil
	case "TeamMember":
		return "team_member", nil
	case "TeamRole":
		return "team_role", nil
	case "User":
		return `"user"`, nil
	case "Workspace":
		return "workspace", nil
	}

	return "", fmt.Errorf("GetTableNameForStruct: %v", errors.New("Struct not found: "+reflect.TypeOf(t).Name()))
}

func (db *dbClient) GetColumnNamesForStruct(structType interface{}) string {
	s := reflect.TypeOf(structType)
	columnNames := ""
	for i := 0; i < s.NumField(); i++ {
		r := s.Field(i)
		if r.Type.Kind() == reflect.Pointer {
			r = reflect.Indirect(reflect.ValueOf(structType)).Type().Field(i)
		}
		if r.Type.Kind() == reflect.Slice {
			continue
		}
		if r.Type.Kind() == reflect.Array {
			continue
		}

		switch jsonTag := r.Tag.Get("json"); jsonTag {
		case "-":
			continue
		case "":
			continue
		default:
			parts := strings.Split(jsonTag, ",")
			columnName := parts[0]
			if columnName == "" {
				continue
			}

			if i == 0 {
				columnNames = columnName
			} else {
				columnNames = fmt.Sprintf("%s, %s", columnNames, columnName)
			}
		}
	}

	return columnNames
}

//Names for composite tables in database
const (
	PROJECT_TEAM_MEMBER    = "project_team_member"
	PROJECT_TEAM_GROUP     = "project_team_group"
	TEAM_GROUP_TEAM_MEMBER = "team_group_team_member"
	TASK_TAG               = "task_tag"
)

func (db *dbClient) GetInsertQueryForCompositeTable(tableName string, valuesMap map[string]interface{}) (string, error) {
	query := ""
	columnNames := ""
	values := ""

	i := 0
	for k, v := range valuesMap {
		if i > 0 {
			columnNames = fmt.Sprintf("%s, ", columnNames)
			values = fmt.Sprintf("%s, ", values)
		}

		columnNames = fmt.Sprintf("%s%s", columnNames, k)

		kind := reflect.ValueOf(v).Kind()
		if kind == reflect.Int || kind == reflect.Int64 {
			values = fmt.Sprintf("%s%d", values, v)
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			values = fmt.Sprintf("%s%f", values, v)
		} else if kind == reflect.Bool {
			values = fmt.Sprintf("%s%t", values, v)
		} else {
			values = fmt.Sprintf("%s'%s'", values, v)
		}

		i++
	}

	if i <= 0 {
		return "", errors.New("getInsertqueryforcompositetable: no values given")

	} else {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnNames, values)
	}

	return query, nil
}

func (db *dbClient) GetSelectQueryForCompositeTable(tableName string, searchParams map[string]interface{}) (string, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	i := 0
	for k, v := range searchParams {
		val := reflect.ValueOf(v)
		kind := val.Kind()
		if kind == reflect.Int || kind == reflect.Int64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Int())
			} else {
				query = fmt.Sprintf("%s && %s = %d", query, k, val.Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Float())
			} else {
				query = fmt.Sprintf("%s && %s = %f", query, k, val.Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Bool())
			} else {
				query = fmt.Sprintf("%s && %s = %t", query, k, val.Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = '%s'", query, k, val.String())
			} else {
				query = fmt.Sprintf("%s && %s = '%s'", query, k, val.String())
			}
		}

		i++
	}

	return query, nil
}

func (db *dbClient) GetUpdateQueryForCompositeTable(tableName string, searchParams map[string]interface{}, updates map[string]interface{}) (string, error) {
	query := fmt.Sprintf("UPDATE %s", tableName)

	i := 0
	for k, v := range updates {
		kind := reflect.ValueOf(v).Kind()
		val := reflect.ValueOf(v)
		if kind == reflect.Int || kind == reflect.Int64 {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = %d", query, k, val.Int())
			} else {
				query = fmt.Sprintf("%s, %s = %d", query, k, val.Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = %f", query, k, val.Float())
			} else {
				query = fmt.Sprintf("%s, %s = %f", query, k, val.Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = %t", query, k, val.Bool())
			} else {
				query = fmt.Sprintf("%s, %s = %t", query, k, val.Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = '%s'", query, k, val.String())
			} else {
				query = fmt.Sprintf("%s, %s = '%s'", query, k, val.String())
			}
		}

		i++
	}

	i = 0
	for k, v := range searchParams {
		val := reflect.ValueOf(v)
		kind := val.Kind()
		if kind == reflect.Int || kind == reflect.Int64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Int())
			} else {
				query = fmt.Sprintf("%s && %s = %d", query, k, val.Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Float())
			} else {
				query = fmt.Sprintf("%s && %s = %f", query, k, val.Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Bool())
			} else {
				query = fmt.Sprintf("%s && %s = %t", query, k, val.Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = '%s'", query, k, val.String())
			} else {
				query = fmt.Sprintf("%s && %s = '%s'", query, k, val.String())
			}
		}

		i++
	}

	return query, nil
}

func (db *dbClient) GetDeleteQueryForCompositeTable(tableName string, searchParams map[string]interface{}) (string, error) {
	query := fmt.Sprintf("DELETE FROM %s", tableName)

	i := 0
	for k, v := range searchParams {
		kind := reflect.ValueOf(v).Kind()
		val := reflect.ValueOf(v)
		if kind == reflect.Int || kind == reflect.Int64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Int())
			} else {
				query = fmt.Sprintf("%s && %s = %d", query, k, val.Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Float())
			} else {
				query = fmt.Sprintf("%s && %s = %f", query, k, val.Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Bool())
			} else {
				query = fmt.Sprintf("%s && %s = %t", query, k, val.Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = '%s'", query, k, val.String())
			} else {
				query = fmt.Sprintf("%s && %s = '%s'", query, k, val.String())
			}
		}

		i++
	}

	return query, nil
}
