package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Customers struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	FirstName     string
	LastName      string
	FullName      string
	MotherName    string
	Nationality   string
	Income        decimal.Decimal
	Type          PeopleType
	StatusMarital MaritalStatus
	Employments   []CustomersEmployment   `gorm:"foreignKey:CustomerID"`
	Address       []CustomersAddress      `gorm:"foreignKey:CustomerID"`
	Documents     []CustomersDocuments    `gorm:"foreignKey:CustomerID"`
	BankAccounts  []CustomersBankAccounts `gorm:"foreignKey:CustomerID"`
}

type CustomersAddress struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CustomerID int `gorm:"index:,option:CONCURRENTLY;not null"`
	Number     string
	Complement string
	Address    string
	State      string
	Country    string
	Type       AddressTypes
}

type CustomersDocuments struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerID    int       `gorm:"index:,option:CONCURRENTLY;not null"`
	Type          DocumentsTypes
	Number        string
	IssueDate     time.Time
	Issue         time.Time
	IssueLocality string
}

type CustomersBankAccounts struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerID   int       `gorm:"index:,option:CONCURRENTLY;not null"`
	BankCode     string
	Account      string
	AccountDigit string
	Agency       string
	AgencyDigit  string
	AccountType  AccountBankType
}

type CustomersEmployment struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerID          int       `gorm:"index:,option:CONCURRENTLY;not null"`
	Kind                string
	position            string
	employer            string
	PhoneNumber         string
	Address             string
	MonthlyIncomeRange  string
	MonthlyIncome       int
	main                decimal.Decimal
	ActualMonthlyIncome int
	PayDay              int
}
