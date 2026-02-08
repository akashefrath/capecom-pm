package bootstrap

import (
	"log"

	"gorm.io/gorm"
)

func SeedMasterData(db *gorm.DB) error {
	log.Println("Starting master data seeding...")

	queries := []string{

		// roles
		`INSERT IGNORE INTO roles (uuid,name,status) VALUES
		('550e8400-e29b-41d4-a716-446655440001','Super Admin','active'),
		('550e8400-e29b-41d4-a716-446655440002','Admin','active'),
		('550e8400-e29b-41d4-a716-446655440003','Manager','active'),
		('550e8400-e29b-41d4-a716-446655440004','Employee','active')`,

		// groups
		`INSERT IGNORE INTO groups (uuid,name,status) VALUES
		('650e8400-e29b-41d4-a716-446655440001','Head Office','active'),
		('650e8400-e29b-41d4-a716-446655440002','Branch A','active'),
		('650e8400-e29b-41d4-a716-446655440003','Branch B','active')`,

		// designations
		`INSERT IGNORE INTO designations (uuid,designation,status) VALUES
		('750e8400-e29b-41d4-a716-446655440001','CEO','active'),
		('750e8400-e29b-41d4-a716-446655440002','Manager','active'),
		('750e8400-e29b-41d4-a716-446655440003','Team Leader','active'),
		('750e8400-e29b-41d4-a716-446655440004','Employee','active')`,

		// departments
		`INSERT IGNORE INTO departments (uuid,department,status) VALUES
		('850e8400-e29b-41d4-a716-446655440001','HR','active'),
		('850e8400-e29b-41d4-a716-446655440002','Finance','active'),
		('850e8400-e29b-41d4-a716-446655440003','Operations','active'),
		('850e8400-e29b-41d4-a716-446655440004','Support','active'),
		('850e8400-e29b-41d4-a716-446655440005','Developer','active')`,
	}

	for i, q := range queries {
		log.Printf("Executing query %d...\n", i+1)
		result := db.Exec(q)
		if result.Error != nil {
			log.Printf("Error executing query %d: %v\n", i+1, result.Error)
			return result.Error
		}
		log.Printf("Query %d executed successfully. Rows affected: %d\n", i+1, result.RowsAffected)
	}

	log.Println("Master data seeding completed successfully!")
	return nil
}
