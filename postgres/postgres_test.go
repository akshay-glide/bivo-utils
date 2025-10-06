package postgres

import (
	"testing"

	"gorm.io/gorm"
)

type TestUser struct {
	gorm.Model
	Name string
	Age  int
	City string
}

func TestPostgresConnectionAndMigration(t *testing.T) {

	// 1. Provide your Postgres configuration here
	config := &ConfigPostgres{
		Host:             ptrString("localhost"),
		Port:             ptrInt(5432),
		Database:         ptrString("test_db"),
		User:             ptrString("root"),
		Password:         ptrString("secret"),
		SSLMode:          ptrString("disable"),
		ConnMaxOpen:      ptrInt(10),
		ConnMaxIdleTime:  ptrInt64(40000),
		ConnMaxIdleConns: ptrInt(5),
	}

	// 2. Test connection with DB
	db, err := NewPostgresDB(config)
	if err != nil {
		t.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// 3. Test migration
	err = Migrate(db, &TestUser{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 4. Test CRUD Operations in DB
	user := TestUser{Name: "Alice", Age: 30, City: "Wonderland"}
	result := db.Create(&user)
	if result.Error != nil {
		t.Errorf("Failed to create user: %v", result.Error)
	}

	var found TestUser
	result = db.First(&found, user.ID)
	if result.Error != nil {
		t.Errorf("Failed to read user: %v", result.Error)
	}

	if found.Name != user.Name {
		t.Errorf("Expected %s, got %s", user.Name, found.Name)
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrInt(i int) *int {
	return &i
}

func ptrString(s string) *string {
	return &s
}
