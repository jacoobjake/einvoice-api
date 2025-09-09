package seeders

import (
	"context"
	"log"

	"github.com/aarondl/opt/omit"
	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/stephenafamo/bob"
	"golang.org/x/crypto/bcrypt"
)

// Implement user seeding logic here
func SeedUsers(db *bob.DB) error {
	if err := seedSuperAdmin(db); err != nil {
		log.Fatalf("failed to seed super admin: %v", err)
		return err
	}

	return nil
}

// Implement super admin seeding logic here
func seedSuperAdmin(db *bob.DB) error {
	ctx := context.Background()
	pw, err := bcrypt.GenerateFromPassword(
		[]byte("superadminpassword"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	_, insErr := models.Users.Insert(&models.UserSetter{
		FirstName: omit.From("super"),
		LastName:  omit.From("admin"),
		Email:     omit.From("superadmin@example.com"),
		Password:  omit.From(string(pw)),
	}).One(ctx, db)

	if insErr != nil {
		return insErr
	}

	return nil
}
