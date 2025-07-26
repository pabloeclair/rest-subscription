package repositories

import (
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"gorm.io/gorm"
)

var ErrBadRequest error

type SubscribeRepository interface {
	Create(subscribe *models.Subscribe) error
	FindAll() ([]*models.Subscribe, error)
	FindByID(id uint) (*models.Subscribe, error)
	FindByUserId(userId string) ([]*models.Subscribe, error)
	FindByServiceName(serviceName string) ([]*models.Subscribe, error)
	Update(id uint, subscribe *models.Subscribe) error
	Delete(id uint) error
}

type GormSubscribeRepository struct {
	Db *gorm.DB
}

func (r *GormSubscribeRepository) Create(subscribe *models.Subscribe) error {
	return r.Db.Create(subscribe).Error
}

func (r *GormSubscribeRepository) FindAll() ([]*models.Subscribe, error) {
	if err := r.Db.Find(&[]*models.Subscribe{}).Error; err != nil {
		return nil, err
	}
	return []*models.Subscribe{}, nil
}

func (r *GormSubscribeRepository) FindByID(id uint) (*models.Subscribe, error) {
	subscribe := &models.Subscribe{}
	if err := r.Db.First(subscribe, id).Error; err != nil {
		return nil, err
	}
	return subscribe, nil
}

func (r *GormSubscribeRepository) FindByUserId(userId string) ([]*models.Subscribe, error) {
	subscribes := []*models.Subscribe{}
	if userId == "" {
		return nil, ErrBadRequest
	}
	if err := r.Db.Where(&models.Subscribe{UserId: userId}).Find(&subscribes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return subscribes, nil
}

func (r *GormSubscribeRepository) FindByServiceName(serviceName string) ([]*models.Subscribe, error) {
	subscribes := []*models.Subscribe{}
	if serviceName == "" {
		return nil, ErrBadRequest
	}
	if err := r.Db.Where(&models.Subscribe{ServiceName: serviceName}).Find(&subscribes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return subscribes, nil
}

func (r *GormSubscribeRepository) Update(id uint, subscribe *models.Subscribe) error {
	res := r.Db.Model(&models.Subscribe{}).Omit("id").Where("id = ?", id).Updates(subscribe)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormSubscribeRepository) Delete(id uint) error {
	res := r.Db.Delete(&models.Subscribe{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
