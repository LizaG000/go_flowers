package service

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
)

func Paginate[T any](
	items []T,
	limit int,
	offset int,
) dto.Pagination[T] {
	var pagination dto.Pagination[T]

	pagination.Total = len(items)

	pagination.Limit = limit

	pagination.Offset = offset

	if (offset-1)*limit < pagination.Total {
		pagination.HasNext = true
	} else {
		pagination.HasNext = false
	}

	if (offset-1)*limit > 0 {
		pagination.HasPrevious = true
	} else {
		pagination.HasPrevious = false
	}

	pagination.Items = items[(offset-1)*limit : (offset)*limit]
	return pagination
}
