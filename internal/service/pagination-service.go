package service

import (
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
)

func Paginate[T any](
	items []T,
	limit int,
	offset int,
) (dto.Pagination[T], error) {
	var pagination dto.Pagination[T]

	pagination.Total = len(items)

	pagination.Limit = limit

	pagination.Offset = offset

	if (offset)*limit < pagination.Total {
		pagination.HasNext = true
	} else {
		pagination.HasNext = false
	}

	if (offset-1)*limit > 0 {
		pagination.HasPrevious = true
	} else {
		pagination.HasPrevious = false
	}
	start := (offset - 1) * limit

	end := offset * limit
	if end > pagination.Total {
		end = pagination.Total
	}
	if start >= pagination.Total {
		return dto.Pagination[T]{}, fmt.Errorf(
			"страница %d не существует: всего элементов %d",
			offset,
			pagination.Total,
		)
	}

	pagination.Items = items[start:end]
	return pagination, nil
}
