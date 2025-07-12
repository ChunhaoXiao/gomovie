package dto

type ActorCreateDto struct {
	Name        string   `form:"name"`
	Pictures    []string `form:"pictures[]"`
	IsRecommand uint     `form:"isRecommand"`
}
