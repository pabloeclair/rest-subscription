package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn: mockDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}

	return db, mock, nil
}

func TestSubscribeCreate(t *testing.T) {
	db, mock, err := NewMock()
	assert.NoError(t, err)

	repo := GormSubscribeRepository{Db: db}
	endDate := time.Date(2025, time.August, 26, 0, 0, 0, 0, time.Local)
	subscribeTest := &models.Subscribe{
		ServiceName: "Kinopoisk",
		Price:       399,
		UserId:      "6061fee-2bf1-aef6f-763675gre",
		StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
		EndDate:     &endDate,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "subscribes"`).
		WithArgs(
			subscribeTest.ServiceName,
			subscribeTest.Price,
			subscribeTest.UserId,
			subscribeTest.StartDate,
			subscribeTest.EndDate,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(subscribeTest)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscribeFindAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)

		repo := GormSubscribeRepository{Db: db}
		subscribesExpected := []*models.Subscribe{
			{
				ID:          1,
				ServiceName: "Kinopoisk",
				Price:       399,
				UserId:      "6061fee-2bf1-aef6f-763675gre",
				StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
				EndDate:     nil,
			},
			{
				ID:          2,
				ServiceName: "Kinopoisk",
				Price:       199,
				UserId:      "708gr-26896-agrfrf-fr5655gre",
				StartDate:   time.Date(2025, time.July, 15, 0, 0, 0, 0, time.Local),
				EndDate:     nil,
			},
		}

		mock.ExpectQuery(`SELECT \* FROM "subscribes"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
				AddRow(
					subscribesExpected[0].ID,
					subscribesExpected[0].ServiceName,
					subscribesExpected[0].Price,
					subscribesExpected[0].UserId,
					subscribesExpected[0].StartDate,
					subscribesExpected[0].EndDate,
				).
				AddRow(
					subscribesExpected[1].ID,
					subscribesExpected[1].ServiceName,
					subscribesExpected[1].Price,
					subscribesExpected[1].UserId,
					subscribesExpected[1].StartDate,
					subscribesExpected[1].EndDate,
				))

		subscribes, err := repo.FindAll()
		assert.NoError(t, err)
		for i, subscribe := range subscribes {
			assert.Equal(t, subscribesExpected[i], subscribe)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscribeFindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)

		repo := GormSubscribeRepository{Db: db}
		endDate := time.Date(2025, time.August, 26, 0, 0, 0, 0, time.Local)
		subscribeExpected := &models.Subscribe{
			ID:          1,
			ServiceName: "Kinopoisk",
			Price:       399,
			UserId:      "6061fee-2bf1-aef6f-763675gre",
			StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
			EndDate:     &endDate,
		}

		mock.ExpectQuery(`SELECT \* FROM "subscribes" WHERE "subscribes"."id" = \$1 ORDER BY "subscribes"."id" LIMIT \$2`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
				AddRow(
					subscribeExpected.ID,
					subscribeExpected.ServiceName,
					subscribeExpected.Price,
					subscribeExpected.UserId,
					subscribeExpected.StartDate,
					subscribeExpected.EndDate,
				))

		subscribe, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, subscribe)
		assert.Equal(t, subscribeExpected, subscribe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrRecordNotFound", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		mock.ExpectQuery(`SELECT \* FROM "subscribes" WHERE "subscribes"."id" = \$1 ORDER BY "subscribes"."id" LIMIT \$2`).
			WithArgs(1, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		subscribe, err := repo.FindByID(1)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, subscribe)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscribeFindByServiceName(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		subscribesExpected := []*models.Subscribe{
			{
				ID:          1,
				ServiceName: "Kinopoisk",
				Price:       399,
				UserId:      "6061fee-2bf1-aef6f-763675gre",
				StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
				EndDate:     nil,
			},
			{
				ID:          12,
				ServiceName: "Kinopoisk",
				Price:       199,
				UserId:      "708gr-26896-agrfrf-fr5655gre",
				StartDate:   time.Date(2025, time.July, 15, 0, 0, 0, 0, time.Local),
				EndDate:     nil,
			},
		}

		mock.ExpectQuery(`SELECT \* FROM "subscribes" WHERE "subscribes"."service_name" = \$1`).
			WithArgs("Kinopoisk").
			WillReturnRows(sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
				AddRow(
					subscribesExpected[0].ID,
					subscribesExpected[0].ServiceName,
					subscribesExpected[0].Price,
					subscribesExpected[0].UserId,
					subscribesExpected[0].StartDate,
					subscribesExpected[0].EndDate,
				).
				AddRow(
					subscribesExpected[1].ID,
					subscribesExpected[1].ServiceName,
					subscribesExpected[1].Price,
					subscribesExpected[1].UserId,
					subscribesExpected[1].StartDate,
					subscribesExpected[1].EndDate,
				))

		subescribes, err := repo.FindByServiceName("Kinopoisk")
		assert.NoError(t, err)
		for i, sub := range subescribes {
			assert.Equal(t, subscribesExpected[i], sub)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessEmpty", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		mock.ExpectQuery(`SELECT \* FROM "subscribes" WHERE "subscribes"."service_name" = \$1`).
			WithArgs("Kinopoisk").
			WillReturnError(gorm.ErrRecordNotFound)

		subescribes, err := repo.FindByServiceName("Kinopoisk")
		assert.NoError(t, err)
		assert.Len(t, subescribes, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscribeFindByUserId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		subscribesActual := []*models.Subscribe{
			{
				ID:          1,
				ServiceName: "Kinopoisk",
				Price:       399,
				UserId:      "6061fee-2bf1-aef6f-763675gre",
				StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
				EndDate:     nil,
			},
			{
				ID:          12,
				ServiceName: "Kinopoisk",
				Price:       199,
				UserId:      "6061fee-2bf1-aef6f-763675gre",
				StartDate:   time.Date(2025, time.July, 15, 0, 0, 0, 0, time.Local),
				EndDate:     nil,
			},
		}

		mock.ExpectQuery(`SELECT \* FROM "subscribes" WHERE "subscribes"."user_id" = \$1`).
			WithArgs("6061fee-2bf1-aef6f-763675gre").
			WillReturnRows(sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
				AddRow(
					subscribesActual[0].ID,
					subscribesActual[0].ServiceName,
					subscribesActual[0].Price,
					subscribesActual[0].UserId,
					subscribesActual[0].StartDate,
					subscribesActual[0].EndDate,
				).
				AddRow(
					subscribesActual[1].ID,
					subscribesActual[1].ServiceName,
					subscribesActual[1].Price,
					subscribesActual[1].UserId,
					subscribesActual[1].StartDate,
					subscribesActual[1].EndDate,
				))

		subescribes, err := repo.FindByUserId("6061fee-2bf1-aef6f-763675gre")
		assert.NoError(t, err)
		for i, sub := range subescribes {
			assert.Equal(t, subscribesActual[i], sub)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessEmpty", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		mock.ExpectQuery(`SELECT \* FROM "subscribes" WHERE "subscribes"."user_id" = \$1`).
			WithArgs("6061fee-2bf1-aef6f-763675gre").
			WillReturnError(gorm.ErrRecordNotFound)

		subescribes, err := repo.FindByUserId("6061fee-2bf1-aef6f-763675gre")
		assert.NoError(t, err)
		assert.Len(t, subescribes, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscribeUpdate(t *testing.T) {

	t.Run("SuccessWithEndDate", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}
		endDate := time.Date(2025, time.August, 26, 0, 0, 0, 0, time.Local)
		subscribeTest := &models.Subscribe{
			ID:          1,
			ServiceName: "Kinopoisk",
			Price:       399,
			UserId:      "6061fee-2bf1-aef6f-763675gre",
			StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
			EndDate:     &endDate,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "subscribes" 
			SET "service_name"=\$1,"price"=\$2,"user_id"=\$3,"start_date"=\$4,"end_date"=\$5 
			WHERE id = \$6`).
			WithArgs(
				subscribeTest.ServiceName,
				subscribeTest.Price,
				subscribeTest.UserId,
				subscribeTest.StartDate,
				subscribeTest.EndDate,
				subscribeTest.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Update(1, subscribeTest)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessWithoutEndDate", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}
		subscribeTest := &models.Subscribe{
			ID:          1,
			ServiceName: "Kinopoisk",
			Price:       399,
			UserId:      "6061fee-2bf1-aef6f-763675gre",
			StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
			EndDate:     nil,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "subscribes" 
			SET "service_name"=\$1,"price"=\$2,"user_id"=\$3,"start_date"=\$4 
			WHERE id = \$5`).
			WithArgs(
				subscribeTest.ServiceName,
				subscribeTest.Price,
				subscribeTest.UserId,
				subscribeTest.StartDate,
				subscribeTest.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Update(1, subscribeTest)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrRecordNotFound", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}
		subscribeTest := &models.Subscribe{
			ID:          1,
			ServiceName: "Kinopoisk",
			Price:       399,
			UserId:      "6061fee-2bf1-aef6f-763675gre",
			StartDate:   time.Date(2025, time.July, 26, 0, 0, 0, 0, time.Local),
			EndDate:     nil,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "subscribes" 
			SET "service_name"=\$1,"price"=\$2,"user_id"=\$3,"start_date"=\$4 
			WHERE id = \$5`).
			WithArgs(
				subscribeTest.ServiceName,
				subscribeTest.Price,
				subscribeTest.UserId,
				subscribeTest.StartDate,
				subscribeTest.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err = repo.Update(1, subscribeTest)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscribeDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscribes" WHERE "subscribes"."id" = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrRecordNotFound", func(t *testing.T) {
		db, mock, err := NewMock()
		assert.NoError(t, err)
		repo := GormSubscribeRepository{Db: db}

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscribes" WHERE "subscribes"."id" = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err = repo.Delete(1)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
