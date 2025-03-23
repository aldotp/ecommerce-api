package dto

type OrderRequest struct {
	ID int `uri:"id" binding:"required"`
}
