package seeds

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func SeedMasterData(db *sqlx.DB) error {
	log.Println("Starting master data seeding...")

	queries := []string{
		`INSERT IGNORE INTO roles (uuid,name,status) VALUES
		('550e8400-e29b-41d4-a716-446655440001','Super Admin','active'),
		('550e8400-e29b-41d4-a716-446655440002','Admin','active'),
		('550e8400-e29b-41d4-a716-446655440003','Employee','active')`,

		`INSERT IGNORE INTO users (uuid,name,email,password_hash,status) VALUES
         ('660e8400-e29b-41d4-a716-446655440001','Super Admin','admin@mail.com','$2a$06$1iycMKBOrR8d69HWBjg.de2URqDiLElgHoUq5eXrsvaKwetH8sD0K','active')
         `,
		`INSERT IGNORE INTO user_roles (uuid,user_id,role_id,status) VALUES
         ('770e8400-e29b-41d4-a716-446655440001','1','1','active')
         `,
	}

	for i, q := range queries {
		log.Printf("Executing query %d...\n", i+1)
		result, err := db.Exec(q)
		if err != nil {
			log.Printf("Error executing query %d: %v\n", i+1, err)
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error executing query %d: %v\n", i+1, err)
		}
		log.Printf("Query %d executed successfully. Rows affected: %d\n", i+1, rows)
	}

	log.Println("Master data seeding completed successfully!")
	return nil
}
