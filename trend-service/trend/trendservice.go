package trend

import models "productivity-planner/trend-service/model"

type TrendService struct {
	repo models.Repository
}

func NewTrendService(repo models.Repository) *TrendService {
	return &TrendService{
		repo: repo,
	}
}
