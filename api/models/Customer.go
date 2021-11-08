package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Customer struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:50;not null;" json:"name"`
	Author    User      `json:"author"`
	AuthorID  uint32    `sql:"type:int REFERENCES users(id)" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Customer) Prepare() {
	c.ID = 0
	c.Name = c.Name
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Customer) Validate() error {

	if c.Name == "" {
		return errors.New("Required Name")
	}
	return nil
}

func (c *Customer) SaveCustomer(db *gorm.DB) (*Customer, error) {
	var err error

	err = db.Debug().Model(&Customer{}).Create(&c).Error
	if err != nil {
		return &Customer{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.AuthorID).Take(&c.Author).Error
		if err != nil {
			return &Customer{}, err
		}
	}

	return c, nil
}

func (c *Customer) FindAllCustomer(db *gorm.DB) (*[]Customer, error) {
	var err error
	customers := []Customer{}
	err = db.Debug().Model(&Customer{}).Limit(100).Find(&customers).Error
	if err != nil {
		return &[]Customer{}, err
	}

	if len(customers) > 0 {
		for i, _ := range customers {
			err := db.Debug().Model(&User{}).Where("id = ?", customers[i].AuthorID).Take(&customers[i].Author).Error
			if err != nil {
				return &[]Customer{}, err
			}
		}
	}

	return &customers, nil
}

func (c *Customer) FindCustomerByID(db *gorm.DB, cid uint64) (*Customer, error) {
	var err error
	err = db.Debug().Model(&Customer{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Customer{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.AuthorID).Take(&c.Author).Error
		if err != nil {
			return &Customer{}, err
		}
	}

	return c, nil
}

func (c *Customer) UpdateAXCustomer(db *gorm.DB) (*Customer, error) {

	var err error
	// db = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&Post{}).UpdateColumns(
	// 	map[string]interface{}{
	// 		"title":      p.Title,
	// 		"content":    p.Content,
	// 		"updated_at": time.Now(),
	// 	},
	// )
	// err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	// if err != nil {
	// 	return &Post{}, err
	// }
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &Post{}, err
	// 	}
	// }
	err = db.Debug().Model(&Customer{}).Where("id = ?", c.ID).Updates(Customer{Name: c.Name, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Customer{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.AuthorID).Take(&c.Author).Error
		if err != nil {
			return &Customer{}, err
		}
	}
	return c, nil
}

func (c *Customer) DeleteACustomer(db *gorm.DB, cid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Customer{}).Where("id = ? and author_id = ?", cid, uid).Take(&Customer{}).Delete(&Customer{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Customer not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (c *Customer) IsCheckCustExist(db *gorm.DB, name string) int {
	var count int
	db.Debug().Model(&Customer{}).Where("name = ? ", c.Name).Count(&count)
	return count
}
