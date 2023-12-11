package charts

import "time"

type Chart struct {
	ID        string    `json:"-" gorm:"primary_key"`
	AirportID string    `json:"arpt_id" example:"FAI" gorm:"index"`
	Cycle     int       `json:"cycle" example:"2213" gorm:"index"`
	FromDate  time.Time `json:"from_date" example:"2021-09-01T00:00:00Z"`
	ToDate    time.Time `json:"to_date" example:"2021-09-30T00:00:00Z"`
	ChartCode string    `json:"chart_code" example:"DP"`
	ChartName string    `json:"chart_name" example:"RDFLG FOUR (RNAV)"`
	ChartURL  string    `json:"chart_url" example:"https://aeronav.faa.gov/d-tpp/2212/01234RDFLG.PDF"`
	CreatedAt time.Time `json:"created_at" example:"2021-09-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-09-01T00:00:00Z"`
}
