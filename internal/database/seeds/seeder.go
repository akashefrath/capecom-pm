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

		`INSERT IGNORE INTO attendance_policies ( uuid, name,  min_work_hours_minutes ,  half_day_minutes , late_grace_minutes, early_exit_grace_minutes, max_break_minutes, auto_checkout_time, status, is_default) 
		 VALUES ('330e8400-e29b-41d4-a716-446655440001', 'basic', '480', '240', '15', '10', '60', '1200', 'active', '1')`,

		`INSERT IGNORE INTO attendance_policy_groups ( uuid, name,  attendance_policy_id, status,created_by) 
			VALUES ('330e8400-e29b-41d4-a716-446655440001', 'basic', 1,'active', '1')`,

		`INSERT IGNORE INTO shift_system ( uuid, name, start_time, end_time, checkin_early, checkin_late, checkout_early, checkout_late, is_overnight, is_default, status) 
		 VALUES ( '340e8400-e29b-41d4-a716-446655440001', 'Morning Shift', '09:00:00', '17:00:00', '30', '15', '15', '30', '0', '1', 'active')`,

		`INSERT IGNORE INTO shift_system_groups ( uuid, name,  shift_system_id, status,created_by) 
			VALUES ('330e8400-e29b-41d4-a716-446655440001', 'basic', 1,'active', '1')`,
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
