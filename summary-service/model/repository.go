package models

type Repository interface {
	FindAllSessionsBetweenDates(summaryDao *Summary) ([]Session, error)
}
