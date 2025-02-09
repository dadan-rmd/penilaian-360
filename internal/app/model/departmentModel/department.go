package departmentModel

import "time"

/* Table Definition */
type Department struct {
	Id              int64     `json:"id" gorm:"id"`
	DepartmentName  string    `json:"department_name" gorm:"DepartmentName"`
	CalenderId      string    `json:"calender_id" gorm:"CalenderId"`
	UrlViewCalender string    `json:"url_view_calender" gorm:"UrlViewCalender"`
	RowStatus       string    `json:"row_status" gorm:"RowStatus"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Department) TableName() string {
	return "master_department"
}

// DTO
type (
	DepartmentResponse struct {
		Id   int64  `json:"id" gorm:"id"`
		Name string `json:"name"`
	}
)

func (entity Department) ToDepartmentResponse() DepartmentResponse {
	return DepartmentResponse{
		Id:   entity.Id,
		Name: entity.DepartmentName,
	}
}
