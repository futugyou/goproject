package service

type CountTagRequest struct {
	State uint8 `form:"state,default=1" bingding:"oneof=0 1"`
}

type TagListRequest struct {
	Name  string `form:"name" bindging:"max=100"`
	State uint8  `form:"state,default=1" bingding:"oneof=0 1"`
}

type CreateTagRequest struct {
	Name     string `form:"name" bindging:"required,min=3,max=100"`
	CreateBy string `form:"created_by" bindging:"required,min=3,max=100"`
	State    uint8  `form:"state,default=1" bingding:"oneof=0 1"`
}

type UpdateTagRequest struct {
	Id         uint32 `form:"id" bindging:"required,gte=1"`
	Name       string `form:"name" bindging:"required,min=3,max=100"`
	ModifiedBy string `form:"modified_by" bindging:"required,min=3,max=100"`
	State      uint8  `form:"state,default=1" bingding:"oneof=0 1"`
}

type DeleteTagRequest struct {
	Id uint32 `form:"id" bindging:"required,gte=1"`
}
