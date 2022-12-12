package models

import "time"

type UserIF interface {
	User() *User       // getter
	SetUser(*User)     // setter
	TableName() string // gormで使用するテーブルを強制指定
	Register() error
	Update() error
	Login() error
	GetByEmailAddress(EmailAddress string) (*User, error)
	NeedsValidation() bool
}

type User struct {
	EmailAddress          string     `gorm:"primaryKey; column:EmailAddress"`
	BusinessPartner       int        `gorm:"column:BusinessPartner"`
	BusinessPartnerName   string     `gorm:"column:BusinessPartnerName"`
	Password              string     `gorm:"column:Password"`
	LastLoginAt           *time.Time `gorm:"column:LastLoginAt"`
	Qos                   Qos        `gorm:"column:Qos"`
	IsEncrypt             *bool      `gorm:"column:IsEncrypt"`
	CreatedAt             time.Time  `gorm:"column:CreatedAt"`
	UpdatedAt             time.Time  `gorm:"column:UpdatedAt"`
	DeletedAt             *time.Time `gorm:"column:DeletedAt"`
	BusinessUserFirstName string     `gorm:"column:BusinessUserFirstName"`
	BusinessUserLastName  string     `gorm:"column:BusinessUserLastName"`
	BusinessUserFullName  string     `gorm:"column:BusinessUserFullName"`
	Language              string     `gorm:"column:Language"`
}

const (
	QosDefault = Qos("default")
	QosRaw     = Qos("raw")
)

// Qos defines type of quality of service
type Qos string

func ToQos(s string) Qos {
	if Qos(s) == QosRaw {
		return QosRaw
	}
	return QosDefault
}
