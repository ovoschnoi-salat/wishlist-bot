package repository

type List struct {
	ID      int64  `gorm:"autoIncrement"`
	OwnerID int64  `gorm:"not null"`
	Title   string `gorm:"not null"`
	Open    bool
	Access  []User `gorm:"many2many:list_access;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
