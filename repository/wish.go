package repository

type Wish struct {
	ID              int64  `gorm:"autoIncrement"`
	ListID          int64  `gorm:"not null"`
	List            List   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title           string `gorm:"not null"`
	Url             string
	Description     string
	Price           string
	ReservationFree bool
	ReservedBy      int64
}
