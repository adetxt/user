package domain

type PaginationInfo struct {
	Page      int32
	PageSize  int32
	TotalData int64
}

func MakePaginationInfo(page int32, pageSize int32) *PaginationInfo {
	return &PaginationInfo{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *PaginationInfo) GetOffset() int32 {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationInfo) GetTotalPage() int64 {
	return int64((p.TotalData / int64(p.PageSize)) + 1)
}

func (p *PaginationInfo) SetTotalData(v int64) {
	p.TotalData = v
}
