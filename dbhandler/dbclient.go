package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/qasim-sajid/clockify-api/models"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "dbpass444"
	DB_NAME     = "ClockifyApp"
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
func SetupDB() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(fmt.Errorf("SetupDB: %v", err))
	}

	dbConnection = db
}

func (db *dbClient) RunInsertQuery(query string) (sql.Result, error) {
	if dbConnection == nil {
		SetupDB()
	}

	result, err := dbConnection.Exec(query)

	if err != nil {
		return nil, fmt.Errorf("RunInsertQuery: %v", err)
	}

	return result, nil
}

func (db *dbClient) RunSelectQuery(query string) (*sql.Rows, error) {
	if dbConnection == nil {
		SetupDB()
	}

	rows, err := dbConnection.Query(query)

	if err != nil {
		return nil, fmt.Errorf("RunSelectQuery: %v", err)
	}

	return rows, nil
}

func (db *dbClient) RunUpdateQuery(query string) (sql.Result, error) {
	if dbConnection == nil {
		SetupDB()
	}

	result, err := dbConnection.Exec(query)

	if err != nil {
		return nil, fmt.Errorf("RunUpdateQuery: %v", err)
	}

	return result, nil
}

func (db *dbClient) RunDeleteQuery(query string) (sql.Result, error) {
	if dbConnection == nil {
		SetupDB()
	}

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
		if project.Client != nil {
			query = fmt.Sprintf(`%s, '%s'`, query, project.Client.ID)
		} else {
			query = fmt.Sprintf(`%s, %v`, query, "null")
		}
		if project.Workspace != nil {
			query = fmt.Sprintf(`%s, '%s')`, query, project.Workspace.ID)
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
			tableName, db.GetColumnNamesForStruct(task), task.ID, task.Description, task.Billable, task.StartTime.Format(time.RFC850),
			task.EndTime.Format(time.RFC850), task.Date.Format(time.RFC850), task.IsActive)
		if task.Project != nil {
			query = fmt.Sprintf(`%s, '%s')`, query, task.Project.ID)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "TeamGroup":
		teamGroup := structType.(models.TeamGroup)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s'`,
			tableName, db.GetColumnNamesForStruct(teamGroup), teamGroup.ID, teamGroup.Name)
		if teamGroup.Workspace != nil {
			query = fmt.Sprintf(`%s, '%s')`, query, teamGroup.Workspace.ID)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "TeamMember":
		teamMember := structType.(models.TeamMember)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', %f)`,
			tableName, db.GetColumnNamesForStruct(teamMember), teamMember.ID, teamMember.BillableRate)
		if teamMember.Workspace != nil {
			query = fmt.Sprintf(`%s, '%s'`, query, teamMember.Workspace.ID)
		} else {
			query = fmt.Sprintf(`%s, %v`, query, "null")
		}
		if teamMember.User != nil {
			query = fmt.Sprintf(`%s, '%s'`, query, teamMember.User.Email)
		} else {
			query = fmt.Sprintf(`%s, %v`, query, "null")
		}
		if teamMember.TeamRole != nil {
			query = fmt.Sprintf(`%s, '%s')`, query, teamMember.TeamRole.ID)
		} else {
			query = fmt.Sprintf(`%s, %v)`, query, "null")
		}
	case "TeamRole":
		teamRole := structType.(models.TeamRole)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(teamRole), teamRole.ID, teamRole.Role)
	case "User":
		user := structType.(models.User)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(user), user.ID, user.Email, user.Name)
	case "Workspace":
		workspace := structType.(models.Workspace)
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES ('%s', '%s')`,
			tableName, db.GetColumnNamesForStruct(workspace), workspace.ID, workspace.Name)
	default:
		return ``, fmt.Errorf("GetInsertQuery: %v",
			errors.New("Insert query generation error!"))
	}

	return query, nil
}

func (db *dbClient) GetInsertQueryForStruct(structType interface{}) (string, error) {
	if reflect.ValueOf(structType).Kind() == reflect.Struct {
		tableName, err := db.GetTableNameForStruct(structType)
		if err != nil {
			return ``, fmt.Errorf("GetInsertQueryForStruct: %v", err)
		}

		query := fmt.Sprintf("INSERT INTO %s values(", tableName)

		v := reflect.ValueOf(structType)
		for i := 0; i < v.NumField(); i++ {

			columnName, err := db.GetColumnNameForStructField(v.Field(i).Type().Field(i))

			if err != nil {
				return ``, fmt.Errorf("GetInsertQueryForStruct: %v", err)
			}

			if columnName == "IgnoreField" {
				continue
			}

			kind := v.Field(i).Kind()
			if kind == reflect.Int || kind == reflect.Int64 {
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %d", query, v.Field(i).Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s%f", query, v.Field(i).Float())
				} else {
					query = fmt.Sprintf("%s, %f", query, v.Field(i).Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s%t", query, v.Field(i).Bool())
				} else {
					query = fmt.Sprintf("%s, %t", query, v.Field(i).Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, \"%s\"", query, v.Field(i).String())
				}
			}
		}

		query = fmt.Sprintf("%s)", query)
		return query, nil
	}

	return ``, fmt.Errorf("GetInsertQueryForStruct: %v %v", errors.New("Insert query generation error: "), reflect.ValueOf(structType).Kind())
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
					query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s && %s = %d", query, k, val.Field(i).Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Field(i).Float())
				} else {
					query = fmt.Sprintf("%s && %s = %f", query, k, val.Field(i).Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Field(i).Bool())
				} else {
					query = fmt.Sprintf("%s && %s = %t", query, k, val.Field(i).Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = \"%s\"", k, query, val.Field(i).String())
				} else {
					query = fmt.Sprintf("%s && %s = \"%s\"", k, query, val.Field(i).String())
				}
			}

			i++
		}

		return query, nil
	}

	return ``, fmt.Errorf("GetSelectQueryForStruct: %v", errors.New("Select query generation error!"))
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
					query = fmt.Sprintf("%s SET %s = %d", query, k, val.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %s = %d", query, k, val.Field(i).Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = %f", query, k, val.Field(i).Float())
				} else {
					query = fmt.Sprintf("%s, %s = %f", query, k, val.Field(i).Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = %t", query, k, val.Field(i).Bool())
				} else {
					query = fmt.Sprintf("%s, %s = %t", query, k, val.Field(i).Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s SET %s = \"%s\"", k, query, val.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, %s = \"%s\"", k, query, val.Field(i).String())
				}
			}

			i++
		}

		query = fmt.Sprintf("%s WHERE _id = \"%s\")", query, itemID)
		return query, nil
	}

	return ``, fmt.Errorf("GetUpdateQueryForStruct: %v", errors.New("Update query generation error!"))
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
					query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s && %s = %d", query, k, val.Field(i).Int())
				}
			} else if kind == reflect.Float32 || kind == reflect.Float64 {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Field(i).Float())
				} else {
					query = fmt.Sprintf("%s && %s = %f", query, k, val.Field(i).Float())
				}
			} else if kind == reflect.Bool {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Field(i).Bool())
				} else {
					query = fmt.Sprintf("%s && %s = %t", query, k, val.Field(i).Bool())
				}
			} else {
				if i == 0 {
					query = fmt.Sprintf("%s WHERE %s = \"%s\"", k, query, val.Field(i).String())
				} else {
					query = fmt.Sprintf("%s && %s = \"%s\"", k, query, val.Field(i).String())
				}
			}

			i++
		}

		return query, nil
	}

	return ``, fmt.Errorf("GetSelectQueryForStruct: %v", errors.New("Select query generation error!"))
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
		return "user", nil
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

func (db *dbClient) GetColumnNameForStructField(t reflect.StructField) (string, error) {
	switch jsonTag := t.Tag.Get("json"); jsonTag {
	case "-":
		return "IgnoreField", nil
	case "":
		return "", fmt.Errorf("GetColumnNameForStructField: %v", errors.New("Struct field json tag not found!"))
	default:
		parts := strings.Split(jsonTag, ",")
		columnName := parts[0]
		if columnName == "" {
			return "", fmt.Errorf("GetColumnNameForStructField: %v", errors.New("Struct field json tag not found!"))
		}
		return columnName, nil
	}
}

//Names for composite tables in database
const (
	PROJECT_TEAM_MEMBER    = "project_team_member"
	PROJECT_TEAM_GROUP     = "group"
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
		return "", errors.New("GetInsertQueryForCompositeTable: No values given!")

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
				query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Field(i).Int())
			} else {
				query = fmt.Sprintf("%s && %s = %d", query, k, val.Field(i).Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Field(i).Float())
			} else {
				query = fmt.Sprintf("%s && %s = %f", query, k, val.Field(i).Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Field(i).Bool())
			} else {
				query = fmt.Sprintf("%s && %s = %t", query, k, val.Field(i).Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = \"%s\"", k, query, val.Field(i).String())
			} else {
				query = fmt.Sprintf("%s && %s = \"%s\"", k, query, val.Field(i).String())
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
				query = fmt.Sprintf("%s SET %s = %d", query, k, val.Field(i).Int())
			} else {
				query = fmt.Sprintf("%s, %s = %d", query, k, val.Field(i).Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = %f", query, k, val.Field(i).Float())
			} else {
				query = fmt.Sprintf("%s, %s = %f", query, k, val.Field(i).Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = %t", query, k, val.Field(i).Bool())
			} else {
				query = fmt.Sprintf("%s, %s = %t", query, k, val.Field(i).Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s SET %s = \"%s\"", k, query, val.Field(i).String())
			} else {
				query = fmt.Sprintf("%s, %s = \"%s\"", k, query, val.Field(i).String())
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
				query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Field(i).Int())
			} else {
				query = fmt.Sprintf("%s && %s = %d", query, k, val.Field(i).Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Field(i).Float())
			} else {
				query = fmt.Sprintf("%s && %s = %f", query, k, val.Field(i).Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Field(i).Bool())
			} else {
				query = fmt.Sprintf("%s && %s = %t", query, k, val.Field(i).Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = \"%s\"", k, query, val.Field(i).String())
			} else {
				query = fmt.Sprintf("%s && %s = \"%s\"", k, query, val.Field(i).String())
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
				query = fmt.Sprintf("%s WHERE %s = %d", query, k, val.Field(i).Int())
			} else {
				query = fmt.Sprintf("%s && %s = %d", query, k, val.Field(i).Int())
			}
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %f", query, k, val.Field(i).Float())
			} else {
				query = fmt.Sprintf("%s && %s = %f", query, k, val.Field(i).Float())
			}
		} else if kind == reflect.Bool {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = %t", query, k, val.Field(i).Bool())
			} else {
				query = fmt.Sprintf("%s && %s = %t", query, k, val.Field(i).Bool())
			}
		} else {
			if i == 0 {
				query = fmt.Sprintf("%s WHERE %s = \"%s\"", k, query, val.Field(i).String())
			} else {
				query = fmt.Sprintf("%s && %s = \"%s\"", k, query, val.Field(i).String())
			}
		}

		i++
	}

	return query, nil
}
