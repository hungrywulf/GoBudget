package main

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Item stuff
type Item struct {
	ID          uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID      uint `gorm:"not null;Column:user_id"`
	Date        time.Time
	Catagory    string
	Description string
	Amount      float64 `gorm:"not null;DEFAULT:0"`
	Reoccuring  string  `gorm:"type:ENUM('f', 'daily', 'monthly', 'yearly');DEFAULT:'f';not null;"`
}

// Budget stuff
type Budget struct {
	ID            uint    `gorm:"primary_key;AUTO_INCREMENT"`
	UserID        uint    `gorm:"not null;Column:user_id"`
	MonthlyBudget float64 `gorm:"not null;DEFAULT:0;Column:monthly_budget"`
}

// StockPrice stuff
type StockPrice struct {
	Date          time.Time
	Symbol        string `gorm:"not null"`
	Open          float64
	High          float64
	Low           float64
	Close         float64
	AdjustedClose float64 `gorm:"Column:adjusted_close"`
	Volume        uint
	Dividend      float64
	Split         float64
}

// User stuff
type User struct {
	ID       uint `gorm:"primary_key;AUTO_INCREMENT"`
	Token    string
	Name     string
	Email    string
	Password string
}

// Stock stuff
type Stock struct {
	Symbol  string
	Updated bool
}

// InvestedStock stuff
type InvestedStock struct {
	UserID   uint `gorm:"Column:user_id"`
	BudgetID uint `gorm:"Column:budget_id"`
	Symbol   string
	Date     time.Time
	Amount   float64
}

// Bond stuff
type Bond struct {
	ID             uint
	UserID         uint      `gorm:"Column:user_id"`
	BudgetID       uint      `gorm:"Column:budget_id"`
	FaceValue      float64   `gorm:"Column:face_value"`
	SettlementDate time.Time `gorm:"Column:settlement_date"`
	MaturityDate   time.Time `gorm:"Column:maturity_date"`
	CouponRate     float64   `gorm:"Column:coupon_rate"`
	MarketRate     float64   `gorm:"Column:market_rate"`
	Yield          float64
	BondValue      float64 `gorm:"Column:bond_value"`
	Redemption     float64
}

// Function stuff
type Function struct {
	UserID      uint    `gorm:"Column:user_id"`
	DailySalary float64 `gorm:"Column:daily_salary"`
	DailyBudget float64 `gorm:"Column:daily_budget"`
}

// TableName stuff

func (StockPrice) TableName() string {
	return "stock_price"
}

// TableName stuff
func (InvestedStock) TableName() string {
	return "invested_stock"
}
