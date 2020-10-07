package Struct

type News struct {
	OldCreatedAt string `json:"OldCreatedAt"`
	Title        string `json:"Title"`
	ImageURL     string `json:"ImageURL"`
	RealURL      string `json:"RealURL"`
	Content      string `json:"Content"`
}

type Hoax struct {
	OldCreatedAt string
	Title        string
	ImageURL     string
	RealURL      string
	Content      string `gorm:"type:MEDIUMTEXT"`
}

type Question struct {
	Name   string
	Answer string `gorm:"type:MEDIUMTEXT"`
}

type Protocol struct {
	OldCreatedAt string
	Title        string
	ImageURL     string
	RealURL      string
	DownloadURL  string
	Content      string `gorm:"type:MEDIUMTEXT"`
}

type Education struct {
	OldCreatedAt string
	Title        string
	ImageURL     string
	RealURL      string
	DownloadURL  string
	Category     string
	Content      string `gorm:"type:MEDIUMTEXT"`
}

type ResponseJSON struct {
	Message      string      `json:"Message"`
	ErrorMessage string      `json:"ErrorMessage"`
	Status       string      `json:"Status"`
	Data         interface{} `json:"Data"`
}
