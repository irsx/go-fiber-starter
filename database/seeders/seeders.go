package seeders

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Seed(db *gorm.DB) error
}

type SeederRunner struct {
	db      *gorm.DB
	seeders []Seeder
}

func (r *SeederRunner) Run() error {
	for _, seeder := range r.seeders {
		if err := seeder.Seed(r.db); err != nil {
			return err
		}
	}
	return nil
}

func NewSeederRunner(db *gorm.DB, seeders ...Seeder) *SeederRunner {
	return &SeederRunner{
		db:      db,
		seeders: seeders,
	}
}

func All(db *gorm.DB) *SeederRunner {
	return NewSeederRunner(db,
		&UserSeeder{},
		&NewsSeeder{},
		&ProductSeeder{},
	)
}
