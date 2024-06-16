package tests

import (
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var testDbInstance *gorm.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase(DbName)
	defer testDB.TearDown()
	testDbInstance = testDB.DbInstance
	err := ExecuteSQLsFromDir(testDbInstance, TestDataDir)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}
