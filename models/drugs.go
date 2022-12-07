package models

type Drug struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"size:30;unique" json:"name"`
	Stock        uint          `json:"stock"`
	Price        uint          `json:"price"`
	Transactions []Transaction `gorm:"foreignKey:DrugID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transactions"`
}

type drugModel struct{}

func NewDrugModel() drugModel {
	return drugModel{}
}

// GetAllDrug
// return all drug
func (d drugModel) GetAllDrug() []Drug {
	var drugs []Drug
	db.Model(&Drug{}).Preload("Transactions").Find(&drugs)
	return drugs
}

// GetDrugByID
// return drug by passing the `id`
func (d drugModel) GetDrugByID(id uint) (Drug, error) {
	var drug Drug
	err := db.Model(&Drug{}).Where("id = ?", id).Preload("Transactions").First(&drug).Error
	return drug, err
}

// CreateNewDrug
// create new drug
func (d drugModel) CreateNewDrug(drug *Drug) error {
	return db.Create(drug).Error
}

// UpdateDrug
// updateing drug
func (d drugModel) UpdateDrug(drug *Drug) error {
	err := db.Save(drug).Error
	return err
}

// DeleteDrug
// delete drug
func (d drugModel) DeleteDrug(drug *Drug) error {
	err := db.Delete(drug).Error
	return err
}
