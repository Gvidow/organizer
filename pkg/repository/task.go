package repository

type Task struct {
	SessionId string  `json:"sessionId"`
	Title     string  `json:"title"`
	Desc      string  `json:"desc"`
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	Day       int     `json:"day"`
	Hour      int     `json:"hour"`
	Min       int     `json:"min"`
	Subj      string  `json:"subj"`
	Exam      string  `json:"exam"`
	Mark      float64 `json:"mark"`
}
