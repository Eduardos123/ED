package repository

import (
	"test/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddUser(t *testing.T) {
	NormalUser := models.User{1, "Test", "18", nil}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	users := CreateRepo(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	t.Logf("Testing insert query for adding user")
	mock.ExpectExec("INSERT user").WithArgs(NormalUser.Name, NormalUser.Age).WillReturnResult(sqlmock.NewResult(1, 1))
	t.Logf("Testing select query for adding user")
	mock.ExpectQuery("SELECT id FROM user WHERE name = (.+) AND age = (.+)").WithArgs(NormalUser.Name, NormalUser.Age).WillReturnRows(rows)

	if _, err = users.AddUser(NormalUser); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}
