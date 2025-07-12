package dto

type MovieCreateDto struct {
	Title      string  `form:"title" binding:"required"`
	Actor      string  `form:"actor"`
	Filename   string  `form:"filename" binding:"required"`
	Duration   int16   `form:"duration"`
	Categories []int16 `form:"categories"`
	ActorId    int     `form:"actorId"`
}
