package models

import (
	"github.com/jinzhu/gorm"
	// This is implicitly needed by gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Services holds information about services used in models package and the database used by them.
type Services struct {
	Values ValuesS
	db     *gorm.DB
}

// NewServices will accept a list of config functions to run.  Each config function will accept a pointer to the Service object, manipulate it, and return an error if there is one.
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
  var s Services
  for _, cfg := range cfgs {
    if err := cfg(&s); err != nil {
      return nil, err
    }
  }
  return &s, nil
}

// ServicesConfig is a type of functional option which return an error
type ServicesConfig func(*Services) error

// WithValues is a functional option that will construct ValuesS.
func WithValues() ServicesConfig {
	return func(s *Services) error {
		s.Values = NewValuesS(s.db)
		return nil
	}
}

// WithGorm is a functional option that will open a connection to GORM, returning an error if something goes wrong.
func WithGorm(dialect, connectionStr string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionStr)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

// WithLogMode is a functional option that configure log mode with the database.
func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

// HELPERS

// Close shuts down the connection to the database
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will attempt to automatically migrate tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Post{}).Error
}

// DestructiveReset will drop tables and call AutoMigrate
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Post{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
