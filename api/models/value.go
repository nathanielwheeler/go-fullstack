package models

import (
	"github.com/jinzhu/gorm"
)

// Value is a struct
type Value struct {
  Title string 
}

// ValuesS handles business logic involving Values and implements ValuesDB
type ValuesS interface{
  ValuesDB
}

type valuesS struct {
  ValuesDB
}

// NewValuesS is a constructor for ValuesS
func NewValuesS(db *gorm.DB) ValuesS {
  return &valuesS{
    ValuesDB: &valuesVal{
      ValuesDB: &valuesGorm{
        db: db,
      },
    },
  }
}

// ValuesDB will handle database interactions for values
type ValuesDB interface {
  GetAll() ([]Value, error)
  Create(value *Value) error
}

type valuesGorm struct {
  db *gorm.DB
}

// GetAll will return all values.
func (vg *valuesGorm) GetAll() ([]Value, error) {
  var values []Value
  if err := vg.db.Find(&values).Error; err != nil {
    return nil, err
  }
  return values, nil
}

func (vg *valuesGorm) Create(value *Value) error {
  return vg.db.Create(value).Error
}

type valuesVal struct {
  ValuesDB
}

func (vv *valuesVal) Create(value *Value) error {
  err := runValuesValFns(value, vv.titleRequired)
  if err != nil {
    return err
  }
  return vv.ValuesDB.Create(value)
}

type valuesValFn func(*Value) error

func runValuesValFns(value *Value, fns ...valuesValFn) error {
  for _, fn := range fns {
    if err := fn(value); err != nil {
      return err
    }
  }
  return nil
}

func (vv *valuesVal) titleRequired(v *Value) error {
  if v.Title == "" {
    return errTitleRequired
  }
  return nil
}