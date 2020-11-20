package infrastructure

import (
	"coredns_api/internal/interface/repository"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Domain struct {
	gorm.Model
	Uuid string `gorm:"index"`
	Name string `gorm:"unique"`
}

type SQLite struct {
	db *gorm.DB
}

func NewSQLite() repository.IDatabase {
	dbPath := os.Getenv("DB_PATH")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Domain{})
	if err != nil {
		panic(err)
	}

	return &SQLite{db: db}
}

func (s *SQLite) SelectDomain(uuid string) (repository.Domain, error) {
	var domain Domain
	result := s.db.First(&domain, "uuid = ?", uuid)
	domainResult := repository.Domain{Uuid: domain.Uuid, Name: domain.Name}
	if result.Error != nil {
		return domainResult, result.Error
	}

	return domainResult, nil
}

func (s *SQLite) GetDomainsList() ([]repository.Domain, error) {
	var domains []Domain
	result := s.db.Find(&domains)

	var domList []repository.Domain
	if result.Error != nil {
		return domList, result.Error
	}

	for _, d := range domains {
		domain := repository.Domain{Uuid: d.Uuid, Name: d.Name}
		domList = append(domList, domain)
	}

	return domList, nil
}

func (s *SQLite) InsertDomain(uuid, name string) error {
	result := s.db.Create(&Domain{Uuid: uuid, Name: name})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SQLite) DeleteDomain(uuid string) error {
	result := s.db.Delete(&Domain{}, "uuid = ?", uuid)
	if result.RowsAffected == 0 {
		return errors.New(repository.GormNotFound)
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}
