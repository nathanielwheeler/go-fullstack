package models

import (
	"github.com/jinzhu/gorm"
)

// #region ERRORS

const (
	// ErrUserIDRequired indicates that there is a missing user ID
	ErrUserIDRequired modelError = "models: user ID is required"
	// ErrNameRequired indicates that there is a missing name
	ErrNameRequired modelError = "models: name is required"
)

// #endregion

// Value will hold all of the information needed for a blog value.
type Value struct {
	gorm.Model
	Name string `gorm:"not_null"`
}

// #region SERVICE

// ValuesService will handle business rules for values.
type ValuesService interface {
	ValuesDB
}

type valuesService struct {
	ValuesDB
}

// NewValuesService is
func NewValuesService(db *gorm.DB) ValuesService {
	return &valuesService{
		ValuesDB: &valuesValidator{
			ValuesDB: &valuesGorm{
				db: db,
			},
		},
	}
}

// #endregion

// #region GORM

//    #region GORM CONFIG

// ValuesDB will handle database interaction for values.
type ValuesDB interface {
	Get(id uint) (*Value, error)
	GetAll() ([]Value, error)
	Create(value *Value) error
	Update(value *Value) error
	Delete(id uint) error
}

type valuesGorm struct {
	db *gorm.DB
}

// Ensure that valuesGorm always implements ValuesDB interface
var _ ValuesDB = &valuesGorm{}

//    #endregion

//    #region GORM METHODS

// ByID will search the values database for a value using input ID.
func (pg *valuesGorm) Get(id uint) (*Value, error) {
	var value Value
	db := pg.db.Where("id = ?", id)
	err := first(db, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// ByYearAndTitle will search the values database for input URL-friendly year and title.
func (pg *valuesGorm) ByYearAndTitle(year int, urlTitle string) (*Value, error) {
	var value Value
	db := pg.db.Where("url_title = ? AND year = ?", urlTitle, year)
	err := first(db, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// GetAll will return all values
func (pg *valuesGorm) GetAll() ([]Value, error) {
	var values []Value
	if err := pg.db.Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

// Create will add a value to the database
func (pg *valuesGorm) Create(value *Value) error {
	return pg.db.Create(value).Error
}

// Update will edit a value in a database
func (pg *valuesGorm) Update(value *Value) error {
	return pg.db.Save(value).Error
}

// Delete will remove a value from default queries.
/* Really, it will add a timestamp for deleted_at, which will exclude the value from normal queries. */
func (pg *valuesGorm) Delete(id uint) error {
	value := Value{Model: gorm.Model{ID: id}}
	return pg.db.Delete(&value).Error
}

//    #endregion

// #endregion

// #region VALIDATOR

type valuesValidator struct {
	ValuesDB
}

//    #region DB VALIDATORS

func (pv *valuesValidator) Create(value *Value) error {
	err := runValuesValFns(value,
		pv.nameRequired)
	if err != nil {
		return err
	}
	return pv.ValuesDB.Create(value)
}

func (pv *valuesValidator) Update(value *Value) error {
	err := runValuesValFns(value,
		pv.nameRequired)
	if err != nil {
		return err
	}
	return pv.ValuesDB.Update(value)
}

func (pv *valuesValidator) Delete(id uint) error {
	var value Value
	value.ID = id
	if err := runValuesValFns(&value, pv.nonZeroID); err != nil {
		return err
	}
	return pv.ValuesDB.Delete(value.ID)
}

//    #endregion

//    #region VAL METHODS

type valuesValFn func(*Value) error

func runValuesValFns(value *Value, fns ...valuesValFn) error {
	for _, fn := range fns {
		if err := fn(value); err != nil {
			return err
		}
	}
	return nil
}

func (pv *valuesValidator) nameRequired(p *Value) error {
	if p.Name == "" {
		return ErrNameRequired
	}
	return nil
}

func (pv *valuesValidator) nonZeroID(value *Value) error {
	if value.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

//    #endregion

// #endregion
