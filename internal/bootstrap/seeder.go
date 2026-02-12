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
		`INSERT IGNORE INTO user_groups (uuid,name,status) VALUES
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

		// ticket types
		`INSERT IGNORE INTO ticket_types (uuid, code, name, description, created_by, status, created_at, updated_at) VALUES
		('a1b2c3d4-e5f6-4a1b-8c2d-1e2f3a4b5c6d', 'FEAT', 'Feature Request', 'New functionality or improvements to existing tools.', 1, 'active', NOW(), NOW()),
		('b2c3d4e5-f6a1-4b2c-9d3e-2f3a4b5c6d7e', 'BUG', 'Bug Report', 'Technical errors or unexpected behavior in the system.', 1, 'active', NOW(), NOW()),
		('c3d4e5f6-a1b2-4c3d-0e4f-3a4b5c6d7e8f', 'ENH', 'Enhancement', 'Optimization of current features for better performance or UX.', 1, 'active', NOW(), NOW()),
		('d4e5f6a1-b2c3-4d4e-1f5a-4b5c6d7e8f9a', 'SEC', 'Security', 'Vulnerability reports or access control issues.',1, 'active', NOW(), NOW()),
		('e5f6a1b2-c3d4-4e5f-2a6b-5c6d7e8f9a0b', 'TASK', 'Maintenance', 'Routine updates, server patches, or database cleanup.', 1, 'active', NOW(), NOW());`,
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
