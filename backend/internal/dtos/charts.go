package dtos

import "github.com/adh-partnership/ids/backend/internal/domain/charts"

type ChartsResponse struct {
	ChartCode string `json:"chart_code" example:"DP"`
	ChartName string `json:"chart_name" example:"RDFLG FOUR (RNAV)"`
	ChartURL  string `json:"chart_url" example:"https://aeronav.faa.gov/d-tpp/2212/01234RDFLG.PDF"`
}

func ChartResponseFromEntity(chart *charts.Chart) ChartsResponse {
	return ChartsResponse{
		ChartCode: chart.ChartCode,
		ChartName: chart.ChartName,
		ChartURL:  chart.ChartURL,
	}
}

func ChartResponsesFromEntities(charts []*charts.Chart) []ChartsResponse {
	responses := make([]ChartsResponse, len(charts))
	for i, chart := range charts {
		responses[i] = ChartResponseFromEntity(chart)
	}
	return responses
}
