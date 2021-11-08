package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Currency struct {
	ID           uint64    `gorm:"primary_key;auto_increment" json:"id"`
	CurrencyFrom int       `gorm:"size:4;not null;" json:"currency_from"`
	CurrencyTo   int       `gorm:"size:4;not null;" json:"currency_to"`
	Rate         float64   `gorm:"default:0" json:"rate"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type RequestCalCurrency struct {
	CurrencyFrom int     `json:"currency_from"`
	CurrencyTo   int     `json:"currency_to"`
	Amount       float64 `json:"amount"`
}
type ResponseCalCurrency struct {
	CurrencyFrom int     `json:"currency_from"`
	CurrencyTo   int     `json:"currency_to"`
	Amount       float64 `json:"amount"`
	Result       float64 `json:"result"`
}

func (c *Currency) Prepare() {
	c.ID = 0
	c.CurrencyFrom = c.CurrencyFrom
	c.CurrencyTo = c.CurrencyTo
	c.Rate = c.Rate
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}
func (c *Currency) Validate() error {

	if c.CurrencyFrom == 0 {
		return errors.New("Required CurrencyFrom")
	}
	if c.CurrencyTo == 0 {
		return errors.New("Required CurrencyTo")
	}
	if c.Rate == 0 {
		return errors.New("Required Rate")
	}
	return nil
}

func (c *Currency) SaveCurrency(db *gorm.DB) (*Currency, error) {
	var err error
	err = db.Debug().Model(&Currency{}).Create(&c).Error
	if err != nil {
		return &Currency{}, err
	}

	return c, nil
}

func (c *Currency) FindAllCurrency(db *gorm.DB) (*[]Currency, error) {
	var err error
	currencies := []Currency{}
	err = db.Debug().Model(&Currency{}).Limit(100).Find(&currencies).Error
	if err != nil {
		return &[]Currency{}, err
	}

	return &currencies, nil
}

func (c *Currency) FindCurrencyByID(db *gorm.DB, cid uint64) (*Currency, error) {
	var err error
	err = db.Debug().Model(&Currency{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Currency{}, err
	}
	return c, nil
}

func (c *Currency) FindCurrencyByCurrencyID(db *gorm.DB, currencyFrom, currencyTo int) (*Currency, error) {
	var err error
	err = db.Debug().Model(&Currency{}).Where("currency_from = ? AND currency_to = ?", currencyFrom, currencyTo).Take(&c).Error
	if err != nil {
		return &Currency{}, err
	}
	return c, nil
}

func (c *Currency) UpdateACurrency(db *gorm.DB) (*Currency, error) {

	var err error

	err = db.Debug().Model(&Currency{}).Where("id = ?", c.ID).Updates(Currency{CurrencyFrom: c.CurrencyFrom, CurrencyTo: c.CurrencyTo, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Currency{}, err
	}
	return c, nil
}

func (c *Currency) DeleteACurrency(db *gorm.DB, cid uint64) (int64, error) {

	db = db.Debug().Model(&Currency{}).Where("id = ?", cid).Take(&Currency{}).Delete(&Currency{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Customer not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Currency) IsCheckExist(db *gorm.DB, currencyfrom, currencyto int) int {
	var count int
	db.Debug().Model(&Currency{}).Where("currency_from = ? AND currency_to = ? ", currencyfrom, currencyto).Count(&count)
	return count
}
