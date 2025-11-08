package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

type TestRepository struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) *TestRepository {
	return &TestRepository{db}
}

func (r *TestRepository) CreateTestFull(ctx context.Context, newTest CreateTest, tagIDs []int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	//Check Tags
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM tags WHERE tag_id = ANY($1)`, pq.Array(tagIDs)).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count != len(tagIDs) {
		return 0, fmt.Errorf("one or more tags do not exist")
	}

	//Create Test
	var testID int
	err = tx.QueryRow(
		`INSERT INTO tests (title, time_limit, type) VALUES ($1, $2, $3) RETURNING test_id`,
		newTest.Title, newTest.TimeLimit, newTest.Type,
	).Scan(&testID)
	if err != nil {
		return 0, err
	}

	//Link Tags
	for _, tagID := range tagIDs {
		_, err = tx.Exec(`INSERT INTO test_tags (test_id, tag_id) VALUES ($1, $2)`, testID, tagID)
		if err != nil {
			return 0, err
		}
	}

	//Create Tasts
	for _, task := range newTest.Tasks {
		optionsJSON, err := json.Marshal(task.Options)
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(
			`INSERT INTO tasks (test_id, type, data_json) VALUES ($1, $2, $3)`,
			testID, task.Type, string(optionsJSON),
		)
		if err != nil {
			return 0, err
		}
	}

	return testID, nil
}

func (r *TestRepository) GetAllTests(ctx context.Context) ([]Test, error) {
	rows, err := r.db.Query(`
		SELECT 
			t.test_id, t.title, t.time_limit, t.type,
			COALESCE(json_agg(DISTINCT jsonb_build_object('tag_id', tg.tag_id, 'name', tg.name)) 
			         FILTER (WHERE tg.tag_id IS NOT NULL), '[]') AS tags,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'task_id', ts.task_id,
				'type', ts.type, 
				'options', ts.data_json
			)) FILTER (WHERE ts.task_id IS NOT NULL), '[]') AS tasks
		FROM tests t
		LEFT JOIN test_tags tt ON t.test_id = tt.test_id
		LEFT JOIN tags tg ON tt.tag_id = tg.tag_id
		LEFT JOIN tasks ts ON t.test_id = ts.test_id
		GROUP BY t.test_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []Test
	for rows.Next() {
		var test Test
		var tagsJSON, tasksJSON []byte

		if err := rows.Scan(&test.TestID, &test.Title, &test.TimeLimit, &test.Type, &tagsJSON, &tasksJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tagsJSON, &test.Tags); err != nil {
			return nil, err
		}

		// Parse tasks with proper handling of the options JSON string
		var tasksRaw []struct {
			TaskID  int    `json:"task_id"`
			Type    string `json:"type"`
			Options string `json:"options"` // This is stored as a JSON string in the database
		}
		if err := json.Unmarshal(tasksJSON, &tasksRaw); err != nil {
			return nil, err
		}

		for _, t := range tasksRaw {
			var options string
			// First unmarshal the JSON string to get the actual options array
			if err := json.Unmarshal([]byte(t.Options), &options); err != nil {
				return nil, err
			}
			
			task := Task{
				CreateTask: CreateTask{
					TestID:  test.TestID,
					Type:    t.Type,
					Options: options,
				},
				TaskID: t.TaskID,
			}
			test.Tasks = append(test.Tasks, task.CreateTask)
		}

		tests = append(tests, test)
	}
	if tests == nil {
		return []Test{}, nil
	}
	return tests, nil
}

func (r *TestRepository) GetTestsByTags(ctx context.Context, tagIDs []int) ([]Test, error) {
    if len(tagIDs) == 0 {
        return r.GetAllTests(ctx)
    }

    rows, err := r.db.Query(`
        SELECT 
            t.test_id, t.title, t.time_limit, t.type,
            COALESCE(json_agg(DISTINCT jsonb_build_object('tag_id', tg.tag_id, 'name', tg.name)) 
                     FILTER (WHERE tg.tag_id IS NOT NULL), '[]') AS tags,
            COALESCE(json_agg(DISTINCT jsonb_build_object(
                'task_id', ts.task_id,
                'type', ts.type, 
                'options', ts.data_json
            )) FILTER (WHERE ts.task_id IS NOT NULL), '[]') AS tasks
        FROM tests t
        LEFT JOIN test_tags tt ON t.test_id = tt.test_id
        LEFT JOIN tags tg ON tt.tag_id = tg.tag_id
        LEFT JOIN tasks ts ON t.test_id = ts.test_id
        WHERE t.test_id IN (
            SELECT test_id FROM test_tags 
            WHERE tag_id = ANY($1) 
            GROUP BY test_id 
            HAVING COUNT(DISTINCT tag_id) = $2
        )
        GROUP BY t.test_id
    `, pq.Array(tagIDs), len(tagIDs))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tests []Test
    for rows.Next() {
        var test Test
        var tagsJSON, tasksJSON []byte

        if err := rows.Scan(&test.TestID, &test.Title, &test.TimeLimit, &test.Type, &tagsJSON, &tasksJSON); err != nil {
            return nil, err
        }

        if err := json.Unmarshal(tagsJSON, &test.Tags); err != nil {
            return nil, err
        }

        // Parse tasks
        var tasksRaw []struct {
            TaskID  int    `json:"task_id"`
            Type    string `json:"type"`
            Options string `json:"options"`
        }
        if err := json.Unmarshal(tasksJSON, &tasksRaw); err != nil {
            return nil, err
        }

        for _, t := range tasksRaw {
            var options string
            if err := json.Unmarshal([]byte(t.Options), &options); err != nil {
                return nil, err
            }
            
            task := Task{
                CreateTask: CreateTask{
                    TestID:  test.TestID,
                    Type:    t.Type,
                    Options: options,
                },
                TaskID: t.TaskID,
            }
            test.Tasks = append(test.Tasks, task.CreateTask)
        }

        tests = append(tests, test)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	if tests == nil {
		return []Test{}, nil
	}
    return tests, nil
}

func (r *TestRepository) GetTestByID(ctx context.Context, testID int) (*Test, error) {
    if testID <= 0 {
        return nil, fmt.Errorf("invalid test ID: %d", testID)
    }

    query := `
        SELECT 
            t.test_id, t.title, t.time_limit, t.type,
            COALESCE(json_agg(DISTINCT jsonb_build_object('tag_id', tg.tag_id, 'name', tg.name)) 
                     FILTER (WHERE tg.tag_id IS NOT NULL), '[]') AS tags,
            COALESCE(json_agg(DISTINCT jsonb_build_object(
                'task_id', ts.task_id,
                'test_id', ts.test_id,
                'type', ts.type, 
                'options', ts.data_json
            )) FILTER (WHERE ts.task_id IS NOT NULL), '[]') AS tasks
        FROM tests t
        LEFT JOIN test_tags tt ON t.test_id = tt.test_id
        LEFT JOIN tags tg ON tt.tag_id = tg.tag_id
        LEFT JOIN tasks ts ON t.test_id = ts.test_id
        WHERE t.test_id = $1
        GROUP BY t.test_id
    `

    var test Test
    var tagsJSON, tasksJSON []byte

    err := r.db.QueryRowContext(ctx, query, testID).Scan(
        &test.TestID, &test.Title, &test.TimeLimit, &test.Type, 
        &tagsJSON, &tasksJSON,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("test with ID %d not found", testID)
        }
        return nil, fmt.Errorf("failed to get test: %w", err)
    }

    // Unmarshal tags
    if err := json.Unmarshal(tagsJSON, &test.Tags); err != nil {
        return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
    }

    // Unmarshal tasks
    var tasksRaw []struct {
        TaskID  int    `json:"task_id"`
        TestID  int    `json:"test_id"`
        Type    string `json:"type"`
        Options string `json:"options"`
    }
    if err := json.Unmarshal(tasksJSON, &tasksRaw); err != nil {
        return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
    }

    for _, t := range tasksRaw {
        var options string
        if err := json.Unmarshal([]byte(t.Options), &options); err != nil {
            return nil, fmt.Errorf("failed to unmarshal task options: %w", err)
        }
        
        task := Task{
            CreateTask: CreateTask{
                TestID:  t.TestID,
                Type:    t.Type,
                Options: options,
            },
            TaskID: t.TaskID,
        }
        test.Tasks = append(test.Tasks, task.CreateTask)
    }

    return &test, nil
}

func (r *TestRepository) DeleteTest(ctx context.Context, testID int) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    // First, check if test exists
    var exists bool
    err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tests WHERE test_id = $1)`, testID).Scan(&exists)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("test with ID %d not found", testID)
    }

    // Delete related records in test_tags (many-to-many relationship)
    _, err = tx.ExecContext(ctx, `DELETE FROM test_tags WHERE test_id = $1`, testID)
    if err != nil {
        return err
    }

    // Delete tasks associated with the test
    _, err = tx.ExecContext(ctx, `DELETE FROM tasks WHERE test_id = $1`, testID)
    if err != nil {
        return err
    }

    // Finally, delete the test itself
    result, err := tx.ExecContext(ctx, `DELETE FROM tests WHERE test_id = $1`, testID)
    if err != nil {
        return err
    }

    // Check if any row was actually deleted
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("test with ID %d not found", testID)
    }

    return nil
}