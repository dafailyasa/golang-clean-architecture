package mysql_test

import (
	"auth-service/domain/repository"

	mysql2 "auth-service/infrastructure/persistence/mysql"
	dataTest "auth-service/test"
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepositoryTest struct {
	suite.Suite

	db         *sql.DB
	gdb        *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.Repository
}

func (t *UserRepositoryTest) SetupTest() {
	db, mock, err := sqlmock.New()
	t.NoError(err)

	t.db = db
	t.mock = mock

	dialector := mysql.New(mysql.Config{
		Conn:                      t.db,
		SkipInitializeWithVersion: true,
	})

	t.gdb, err = gorm.Open(dialector, &gorm.Config{})
	t.Require().NoError(err)

	t.repository = mysql2.NewUserRepository(t.gdb)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTest))
}

var userColumns = []string{
	"id",
	"keycloak_uuid",
	"email",
	"first_name",
	"last_name",
	"username",
	"status",
	"is_admin",
	"created_at",
	"updated_at",
}

func userRow(id uint) *sqlmock.Rows {
	return sqlmock.NewRows(userColumns).
		AddRow(
			id,
			dataTest.UserTestData.KeycloakUUID,
			dataTest.UserTestData.Email.String(),
			dataTest.UserTestData.FirstName,
			dataTest.UserTestData.LastName,
			dataTest.UserTestData.Username,
			dataTest.UserTestData.Status,
			false,
			dataTest.UserTestData.CreatedAt,
			dataTest.UserTestData.UpdatedAt,
		)
}

func emptyUserRows() *sqlmock.Rows {
	return sqlmock.NewRows(userColumns)
}

func (t *UserRepositoryTest) TestCreateUser() {
	q := regexp.QuoteMeta(
		"INSERT INTO `users` (`keycloak_uuid`,`email`,`first_name`,`last_name`,`username`,`status`,`is_admin`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)",
	)

	t.Run("it should be create a user", func() {
		userData := dataTest.UserTestData
		userData.ID = 0

		t.mock.ExpectBegin()

		t.mock.ExpectExec(q).
			WithArgs(
				userData.KeycloakUUID,
				userData.Email.String(),
				userData.FirstName,
				userData.LastName,
				userData.Username,
				userData.Status,
				false,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		t.mock.ExpectCommit()

		err := t.repository.Create(context.Background(), &userData)
		t.NoError(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return an error when database execution fails", func() {
		userData := dataTest.UserTestData
		userData.ID = 0

		t.mock.ExpectBegin()

		t.mock.ExpectExec(q).
			WithArgs(
				userData.KeycloakUUID,
				userData.Email.String(),
				userData.FirstName,
				userData.LastName,
				userData.Username,
				userData.Status,
				false,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).WillReturnError(gorm.ErrDuplicatedKey)

		t.mock.ExpectRollback()

		err := t.repository.Create(context.Background(), &userData)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}

func (t *UserRepositoryTest) TestExistsByEmailAndUsername() {
	q := regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE email = ? OR username = ?")

	t.Run("it should return true when email or username exist", func() {
		rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.Email.String(), dataTest.UserTestData.Username).
			WillReturnRows(rows)

		exists, err := t.repository.ExistsByEmailAndUsername(context.Background(), dataTest.UserTestData.Email.String(), dataTest.UserTestData.Username)
		t.NoError(err)
		t.True(exists)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return error when query fails", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.Email.String(), dataTest.UserTestData.Username).
			WillReturnError(gorm.ErrForeignKeyViolated)

		exists, err := t.repository.ExistsByEmailAndUsername(context.Background(), dataTest.UserTestData.Email.String(), dataTest.UserTestData.Username)
		t.Error(err)
		t.False(exists)
	})
}

func (t *UserRepositoryTest) TestFindByID() {
	q := regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? ORDER BY `users`.`id` LIMIT ?")

	t.Run("it should return user when found", func() {
		t.mock.ExpectQuery(q).WithArgs(uint(1), 1).WillReturnRows(userRow(1))

		usr, err := t.repository.FindByID(context.Background(), uint(1))
		t.NoError(err)
		t.NotNil(usr)
		t.Equal(uint(1), usr.ID)
		t.Equal(dataTest.UserTestData.Email.String(), usr.Email.String())
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return business error when user not found", func() {
		t.mock.ExpectQuery(q).WithArgs(uint(1), 1).WillReturnError(gorm.ErrRecordNotFound)

		usr, err := t.repository.FindByID(context.Background(), uint(1))
		t.Nil(usr)
		t.Error(err)
		t.Error(err, "User not found")
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return infra error when query fails", func() {
		t.mock.ExpectQuery(q).WithArgs(uint(1), 1).WillReturnError(gorm.ErrInvalidData)

		usr, err := t.repository.FindByID(context.Background(), uint(1))
		t.Nil(usr)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}

func (t *UserRepositoryTest) TestFindByEmail() {
	q := regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT ?")

	t.Run("it should return user when email found", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.Email.String(), 1).
			WillReturnRows(userRow(1))

		usr, err := t.repository.FindByEmail(context.Background(), dataTest.UserTestData.Email.String())
		t.NoError(err)
		t.NotNil(usr)
		t.Equal(dataTest.UserTestData.Email.String(), usr.Email.String())
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return business error when email not found", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.Email.String(), 1).
			WillReturnRows(emptyUserRows())

		usr, err := t.repository.FindByEmail(context.Background(), dataTest.UserTestData.Email.String())
		t.Nil(usr)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return infra error when query fails", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.Email.String(), 1).
			WillReturnError(sqlmock.ErrCancelled)

		usr, err := t.repository.FindByEmail(context.Background(), dataTest.UserTestData.Email.String())
		t.Nil(usr)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}

func (t *UserRepositoryTest) TestFindByKeycloakUUID() {
	q := regexp.QuoteMeta("SELECT * FROM `users` WHERE keycloak_uuid = ? ORDER BY `users`.`id` LIMIT ?")

	t.Run("it should return user when keycloak uuid found", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.KeycloakUUID, 1).
			WillReturnRows(userRow(1))

		usr, err := t.repository.FindByKeycloakUUID(context.Background(), dataTest.UserTestData.KeycloakUUID)
		t.NoError(err)
		t.NotNil(usr)
		t.Equal(dataTest.UserTestData.KeycloakUUID, usr.KeycloakUUID)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return business error when keycloak uuid not found", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.KeycloakUUID, 1).
			WillReturnRows(emptyUserRows())

		usr, err := t.repository.FindByKeycloakUUID(context.Background(), dataTest.UserTestData.KeycloakUUID)
		t.Nil(usr)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return infra error when query fails", func() {
		t.mock.ExpectQuery(q).
			WithArgs(dataTest.UserTestData.KeycloakUUID, 1).
			WillReturnError(sqlmock.ErrCancelled)

		usr, err := t.repository.FindByKeycloakUUID(context.Background(), dataTest.UserTestData.KeycloakUUID)
		t.Nil(usr)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}

func (t *UserRepositoryTest) TestUpdateUser() {
	q := regexp.QuoteMeta("UPDATE `users` SET")

	t.Run("it should update user successfully", func() {
		t.mock.ExpectBegin()
		t.mock.ExpectExec(q).
			WillReturnResult(sqlmock.NewResult(0, 1))
		t.mock.ExpectCommit()

		err := t.repository.Update(context.Background(), &dataTest.UserTestData)
		t.NoError(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return infra error when update fails", func() {
		t.mock.ExpectBegin()
		t.mock.ExpectExec(q).
			WillReturnError(sqlmock.ErrCancelled)
		t.mock.ExpectRollback()

		err := t.repository.Update(context.Background(), &dataTest.UserTestData)
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}

func (t *UserRepositoryTest) TestExistEmailByUserID() {
	q := regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE id != ? AND email = ?")

	t.Run("it should return true when email exists for another user", func() {
		rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
		t.mock.ExpectQuery(q).
			WithArgs(uint(1), dataTest.UserTestData.Email.String()).
			WillReturnRows(rows)

		exists, err := t.repository.ExistEmailByUserID(context.Background(), uint(1), dataTest.UserTestData.Email.String())
		t.NoError(err)
		t.True(exists)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return infra error when query fails", func() {
		t.mock.ExpectQuery(q).
			WithArgs(uint(1), dataTest.UserTestData.Email.String()).
			WillReturnError(sqlmock.ErrCancelled)

		exists, err := t.repository.ExistEmailByUserID(context.Background(), uint(1), dataTest.UserTestData.Email.String())
		t.Error(err)
		t.False(exists)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}

func (t *UserRepositoryTest) TestDeleteByID() {
	q := regexp.QuoteMeta("DELETE FROM `users` WHERE id = ?")

	t.Run("it should delete user successfully", func() {
		t.mock.ExpectBegin()
		t.mock.ExpectExec(q).
			WithArgs(uint(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		t.mock.ExpectCommit()

		err := t.repository.DeleteByID(context.Background(), uint(1))
		t.NoError(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return business error when user not found", func() {
		t.mock.ExpectBegin()
		t.mock.ExpectExec(q).
			WithArgs(uint(1)).
			WillReturnError(gorm.ErrRecordNotFound)
		t.mock.ExpectRollback()

		err := t.repository.DeleteByID(context.Background(), uint(1))
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})

	t.Run("it should return infra error when delete query fails", func() {
		t.mock.ExpectBegin()
		t.mock.ExpectExec(q).
			WithArgs(uint(1)).
			WillReturnError(sqlmock.ErrCancelled)
		t.mock.ExpectRollback()

		err := t.repository.DeleteByID(context.Background(), uint(1))
		t.Error(err)
		t.NoError(t.mock.ExpectationsWereMet())
	})
}
