package pagination

// Pagination value object cho phân trang
type Pagination struct {
	Page     int   `json:"page" validate:"min=1"`
	PageSize int   `json:"page_size" validate:"min=1,max=100"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

// NewPagination tạo pagination mới với validation
func NewPagination(page, pageSize int) *Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset tính offset cho database query
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// SetTotal set total records và tính số pages
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	if total > 0 {
		p.Pages = int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	} else {
		p.Pages = 0
	}
}
