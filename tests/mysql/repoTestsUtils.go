package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	DbName = "saladRecipes"
	DbUser = "user"
	DbPass = "pass"
	Port   = "5432/tcp"
	Image  = "mysql:8.4"
)

const (
	TestDataDir = "../test_data/"
)

type TestDatabase struct {
	DbInstance *gorm.DB
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase(dbName string) *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := CreateContainer(ctx, dbName)
	if err != nil {
		log.Fatal("failed to setup test: ", err)
	}

	err = MigrateDb(dbAddr, dbName, "../migrations")
	if err != nil {
		log.Fatal("failed to perform db migration: ", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	instance, _ := tdb.DbInstance.DB()
	instance.Close()
	_ = tdb.container.Terminate(context.Background())
}

func CreateContainer(ctx context.Context, dbName string) (testcontainers.Container, *gorm.DB, string, error) {
	var env = map[string]string{
		"MYSQL_ROOT_PASSWORD": "admin",
		"MYSQL_DATABASE":      "saladRecipes",
	}

	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage(Image),
		testcontainers.WithEnv(env),
		//testcontainers.WithHostPortAccess(3306),
		//mysql.WithDatabase(dbName),
		mysql.WithUsername(DbUser),
		mysql.WithPassword(DbPass),
	)

	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("mysql container ready and running at port: ", p.Port())
	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	mysqlCredentials := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPass, "tcp", dbAddr, dbName)
	db, err := gorm.Open(gormMysql.Open(mysqlCredentials), &gorm.Config{})
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, dbAddr, nil
}

func MigrateDb(dbAddr string, dbName string, pathToMigrationFiles string) error {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path")
	}
	_ = path

	databaseURL := fmt.Sprintf("mysql://%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPass, "tcp", dbAddr, dbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func ExecuteSQLsFromDir(db *gorm.DB, dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error while reading dir with test data: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dir, file.Name())
			err = executeSQL(db, filePath)
			if err != nil {
				return fmt.Errorf("executing sql %s: %w", file.Name(), err)
			}
		}
	}
	return nil
}

func executeSQL(db *gorm.DB, filePath string) error {
	scriptContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading sql script: %w", err)
	}

	err = db.WithContext(context.Background()).Raw(string(scriptContent)).Error
	if err != nil {
		return fmt.Errorf("execution sql script: %w", err)
	}
	return nil
}
