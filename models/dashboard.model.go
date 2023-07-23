package models

type DashboardResponse struct {
	TotalArticle    int64   `json:"total_article,omitempty"`
	TotalGallery    int64   `json:"total_gallery,omitempty"`
	TotalTourism    int64   `json:"total_tourism,omitempty"`
	ArticlePerMonth []int64 `json:"article_per_month,omitempty"`
	GalleryPerMonth []int64 `json:"gallery_per_month,omitempty"`
	TourismPerMonth []int64 `json:"tourism_per_month,omitempty"`
}
