package dto

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"github.com/google/uuid"
)

type RequestUpdateFlower struct {
	ID          uuid.UUID           `json:"id"`
	Update_data entity.UpdateFlower `json:"update_data"`
}
