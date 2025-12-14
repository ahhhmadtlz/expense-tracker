package category

import "time"

type CategoryType string

const (
	TypeIncome CategoryType="income"
	TypeExpense CategoryType="expense"
)

func (c CategoryType)String()string{
	return string(c)
}


func MapToCategoryType(catType string)CategoryType{
	switch catType {
	case "income":
		return TypeIncome
	case "expense":
		return TypeExpense
	default:
		return TypeExpense
	}
}

type Category struct {
	ID uint
	UserID uint
	Name string
	Type CategoryType
	Icon string
	Color string
	CreatedAt time.Time
}