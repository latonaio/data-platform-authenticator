package models

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"data-platform-authenticator/pkg/db"

	validation "github.com/go-ozzo/ozzo-validation"
)

func NewUser() UserIF {
	return &User{}
}

func (u *User) Register() error {
	result := db.ConPool.Con.Create(u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) Update() error {
	result := db.ConPool.Con.Model(u).Updates(u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) Login() error {
	result := db.ConPool.Con.Model(u).UpdateColumn("LastLoginAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) GetByEmailAddress(EmailAddress string) (*User, error) {
	result := db.ConPool.Con.Model(u).Where("EmailAddress = ?", EmailAddress).First(u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (u *User) TableName() string {
	return db.ConPool.Info.TableName
}

func (u *User) User() *User {
	return u
}

func (u *User) SetUser(user *User) {
	u.EmailAddress = user.EmailAddress
	u.BusinessPartner = user.BusinessPartner
	u.BusinessPartnerName = user.BusinessPartnerName
	u.Password = user.Password
	u.BusinessUserFirstName = user.BusinessUserFirstName
	u.BusinessUserLastName = user.BusinessUserLastName
	u.BusinessUserFullName = user.BusinessUserFullName
	u.Language = user.Language
	u.Qos = user.Qos
	u.IsEncrypt = user.IsEncrypt
	//u.CreatedAt = user.CreatedAt
	//u.UpdatedAt = user.UpdatedAt
	//u.DeletedAt = user.DeletedAt
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func (u *User) NeedsValidation() bool {
	return u.Qos == QosDefault
}

func (u User) Validate() error {
	const minEmailAddressLength = 3
	const maxEmailAddressLength = 200
	const minPasswordLength = 8
	const maxPasswordLength = 30

	return validation.ValidateStruct(&u,
		validation.Field(&u.EmailAddress,
			validation.Required,
			validation.Length(minEmailAddressLength, maxEmailAddressLength),
			validation.By(UsableEmailAddress),
		),
		validation.Field(&u.Password,
			validation.Required,
			validation.Length(minPasswordLength, maxPasswordLength),
			validation.By(UsableString),
			validation.By(ContainsUppercase),
			validation.By(ContainsLowercase),
			validation.By(notInclude(u.EmailAddress)),
		),
	)
}

// notInclude checks that value contains str.
func notInclude(str string) validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("failed to cast string")
		}
		if strings.Contains(s, str) {
			return errors.New("contains an invalid string")
		}
		return nil
	}
}

// ContainsUppercase checks that value contains uppercase characters.
func ContainsUppercase(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to cast string")
	}
	if !containsUppercase(str) {
		return errors.New("uppercase is not contain")
	}
	return nil
}

// ContainsUppercase checks that value contains lowercase characters.
func ContainsLowercase(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to cast string")
	}
	if !containsLowercase(str) {
		return errors.New("lowercase is not contain")
	}
	return nil
}

// UsableString checks that value consists of only usable characters.
func UsableString(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to cast string")
	}
	if !usableString(str) {
		return errors.New("contains unusable characters")
	}
	return nil
}

func UsableEmailAddress(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to cast string")
	}
	if !usableEmailAddress(str) {
		return errors.New("contains unusable characters")
	}
	return nil
}

func containsUppercase(str string) bool {
	for _, r := range str {
		if 'A' <= r && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(str string) bool {
	for _, r := range str {
		if 'a' <= r && r <= 'z' {
			return true
		}
	}
	return false
}

func usableString(str string) bool {
	for _, r := range str {
		if 'a' <= r && r <= 'z' {
			continue
		}
		if 'A' <= r && r <= 'Z' {
			continue
		}
		if '0' <= r && r <= '9' {
			continue
		}
		if '-' == r {
			continue
		}
		if '_' == r {
			continue
		}
		if '.' == r {
			continue
		}
		if '\'' == r {
			continue
		}
		return false
	}
	return true
}

func usableEmailAddress(emailAddress string) bool {
	_, err := mail.ParseAddress(emailAddress)
	return err == nil
}
