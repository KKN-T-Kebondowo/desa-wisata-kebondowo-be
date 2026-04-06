package e2e

import (
	"fmt"
	"kebondowo/models"
	"testing"
)

// --- SQL Injection Prevention ---

func TestGallerySortByInjection(t *testing.T) {
	// Try to inject SQL via sortby parameter
	resp := doRequest("GET", "/api/galleries/?sortby=1;DROP+TABLE+galleries--&orderedby=desc", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 (whitelist should sanitize), got %d: %s", resp.StatusCode, resp.Body)
	}
	// The server should not crash — the invalid sortby should be replaced with "created_at"
}

func TestGalleryOrderByInjection(t *testing.T) {
	resp := doRequest("GET", "/api/galleries/?sortby=created_at&orderedby=desc;DROP+TABLE+galleries--", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 (whitelist should sanitize), got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestArticleSortByInjection(t *testing.T) {
	resp := doRequest("GET", "/api/articles/?sortby=1;DROP+TABLE+articles--&orderedby=asc", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 (whitelist should sanitize), got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestGalleryValidSortFields(t *testing.T) {
	tests := []struct {
		sortby    string
		orderedby string
	}{
		{"created_at", "asc"},
		{"created_at", "desc"},
		{"updated_at", "asc"},
		{"id", "desc"},
	}
	for _, tt := range tests {
		resp := doRequest("GET", "/api/galleries/?sortby="+tt.sortby+"&orderedby="+tt.orderedby, nil, "")
		if resp.StatusCode != 200 {
			t.Errorf("sortby=%s orderedby=%s: expected 200, got %d", tt.sortby, tt.orderedby, resp.StatusCode)
		}
	}
}

func TestGalleryInvalidLimitOffset(t *testing.T) {
	resp := doRequest("GET", "/api/galleries/?limit=abc", nil, "")
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for invalid limit, got %d", resp.StatusCode)
	}

	resp = doRequest("GET", "/api/galleries/?offset=abc", nil, "")
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for invalid offset, got %d", resp.StatusCode)
	}
}

// --- Visitor Counter Atomicity ---

func TestVisitorCounterIncrement(t *testing.T) {
	// Create an article
	body := map[string]string{
		"title":       "Visitor Test Article",
		"slug":        "visitor-test-e2e",
		"author":      "Tester",
		"content":     "Content for visitor test",
		"picture_url": "https://example.com/visitor.jpg",
	}
	resp := doRequest("POST", "/api/articles/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("failed to create article: %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	article := result["article"].(map[string]interface{})
	articleID := article["id"].(float64)

	// Visit the article multiple times
	for i := 0; i < 5; i++ {
		resp = doRequest("GET", "/api/articles/visitor-test-e2e", nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("visit %d: expected 200, got %d", i, resp.StatusCode)
		}
	}

	// Check visitor count
	resp = doRequest("GET", "/api/articles/visitor-test-e2e", nil, "")
	result = parseJSON(resp.Body)
	visitedArticle := result["article"].(map[string]interface{})
	// 5 visits + this final visit = 6 (but initial GetOne during create response doesn't increment)
	visitorCount := visitedArticle["visitor"].(float64)
	if visitorCount < 5 {
		t.Errorf("expected visitor count >= 5, got %v", visitorCount)
	}

	// Cleanup
	doRequest("DELETE", "/api/articles/"+fmt.Sprintf("%d", int(articleID)), nil, authToken)
}

func TestTourismVisitorCounterIncrement(t *testing.T) {
	// Create a tourism
	body := map[string]interface{}{
		"title":             "Visitor Tourism",
		"slug":              "visitor-tourism-e2e",
		"description":       "For visitor test",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/vt.jpg",
	}
	resp := doRequest("POST", "/api/tourisms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("failed to create tourism: %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	tourism := result["tourism"].(map[string]interface{})
	tourismID := tourism["id"].(float64)

	// Visit 3 times
	for i := 0; i < 3; i++ {
		doRequest("GET", "/api/tourisms/visitor-tourism-e2e", nil, "")
	}

	// Check visitor count
	resp = doRequest("GET", "/api/tourisms/visitor-tourism-e2e", nil, "")
	result = parseJSON(resp.Body)
	visited := result["tourism"].(map[string]interface{})
	visitorCount := visited["visitor"].(float64)
	if visitorCount < 3 {
		t.Errorf("expected visitor count >= 3, got %v", visitorCount)
	}

	// Cleanup
	doRequest("DELETE", "/api/tourisms/"+fmt.Sprintf("%d", int(tourismID)), nil, authToken)
}

// --- Delete Cascading ---

func TestTourismDeleteCascadesPictures(t *testing.T) {
	// Create tourism with pictures
	body := map[string]interface{}{
		"title":             "Cascade Tourism",
		"slug":              "cascade-tourism-e2e",
		"description":       "Test cascade",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/cascade.jpg",
		"pictures": []map[string]string{
			{"picture_url": "https://example.com/c1.jpg"},
			{"picture_url": "https://example.com/c2.jpg"},
			{"picture_url": "https://example.com/c3.jpg"},
		},
	}
	resp := doRequest("POST", "/api/tourisms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("failed to create: %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	tourism := result["tourism"].(map[string]interface{})
	tourismID := tourism["id"].(float64)

	// Verify pictures exist
	var count int64
	testDB.Model(&models.TourismPicture{}).Where("tourism_id = ?", int(tourismID)).Count(&count)
	if count != 3 {
		t.Fatalf("expected 3 pictures before delete, got %d", count)
	}

	// Delete tourism
	resp = doRequest("DELETE", "/api/tourisms/"+fmt.Sprintf("%d", int(tourismID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Verify pictures were deleted
	testDB.Model(&models.TourismPicture{}).Where("tourism_id = ?", int(tourismID)).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 pictures after cascade delete, got %d", count)
	}
}

func TestUMKMDeleteCascadesPictures(t *testing.T) {
	// Create UMKM with pictures
	body := map[string]interface{}{
		"title":             "Cascade UMKM",
		"slug":              "cascade-umkm-e2e",
		"description":       "Test cascade",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/cascade.jpg",
		"contact":           "08123",
		"contact_name":      "Test",
		"pictures": []map[string]string{
			{"picture_url": "https://example.com/u1.jpg"},
			{"picture_url": "https://example.com/u2.jpg"},
		},
	}
	resp := doRequest("POST", "/api/umkms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("failed to create: %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	umkm := result["umkm"].(map[string]interface{})
	umkmID := umkm["id"].(float64)

	// Verify pictures exist
	var count int64
	testDB.Model(&models.UMKMPicture{}).Where("umkm_id = ?", int(umkmID)).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 pictures before delete, got %d", count)
	}

	// Delete UMKM
	resp = doRequest("DELETE", "/api/umkms/"+fmt.Sprintf("%d", int(umkmID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Verify pictures were deleted
	testDB.Model(&models.UMKMPicture{}).Where("umkm_id = ?", int(umkmID)).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 pictures after cascade delete, got %d", count)
	}
}

// --- UMKMPicture Delete Actually Deletes ---

func TestUMKMPictureDeleteActuallyDeletes(t *testing.T) {
	// Create UMKM
	umkmBody := map[string]interface{}{
		"title":             "PicDelete UMKM",
		"slug":              "picdelete-umkm-e2e",
		"description":       "Test",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/cover.jpg",
		"contact":           "08123",
		"contact_name":      "Test",
	}
	resp := doRequest("POST", "/api/umkms/", umkmBody, authToken)
	result := parseJSON(resp.Body)
	umkm := result["umkm"].(map[string]interface{})
	umkmID := uint32(umkm["id"].(float64))

	// Create UMKM picture
	picBody := map[string]interface{}{
		"picture_url": "https://example.com/deleteme.jpg",
		"umkm_id":     umkmID,
	}
	resp = doRequest("POST", "/api/umkm-pictures/", picBody, authToken)
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, resp.Body)
	}
	result = parseJSON(resp.Body)
	pic := result["umkm_picture"].(map[string]interface{})
	picID := pic["id"].(float64)

	// Verify it exists in DB
	var count int64
	testDB.Model(&models.UMKMPicture{}).Where("id = ?", int(picID)).Count(&count)
	if count != 1 {
		t.Fatalf("expected picture to exist before delete")
	}

	// Delete it
	resp = doRequest("DELETE", "/api/umkm-pictures/"+fmt.Sprintf("%d", int(picID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Verify it's actually gone from the database
	testDB.Model(&models.UMKMPicture{}).Where("id = ?", int(picID)).Count(&count)
	if count != 0 {
		t.Errorf("UMKM picture was NOT actually deleted from DB (count=%d)", count)
	}

	// Cleanup
	doRequest("DELETE", "/api/umkms/"+fmt.Sprintf("%d", int(umkmID)), nil, authToken)
}

// --- Auth with invalid token ---

func TestProtectedEndpointWithInvalidToken(t *testing.T) {
	resp := doRequest("POST", "/api/tourisms/", map[string]interface{}{
		"title":             "Should Fail",
		"slug":              "should-fail",
		"description":       "fail",
		"latitude":          0.0,
		"longitude":         0.0,
		"cover_picture_url": "https://example.com/fail.jpg",
	}, "invalid-token-here")

	if resp.StatusCode != 401 {
		t.Errorf("expected 401 for invalid token, got %d", resp.StatusCode)
	}
}

// --- Slug Dedup ---

func TestTourismSlugDedup(t *testing.T) {
	body := map[string]interface{}{
		"title":             "Slug Test",
		"slug":              "slug-dedup-e2e",
		"description":       "First",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/s1.jpg",
	}

	// Create first
	resp := doRequest("POST", "/api/tourisms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("first create failed: %d: %s", resp.StatusCode, resp.Body)
	}
	result1 := parseJSON(resp.Body)
	t1 := result1["tourism"].(map[string]interface{})
	t1ID := t1["id"].(float64)

	// Create second with same slug
	body["description"] = "Second"
	resp = doRequest("POST", "/api/tourisms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("second create failed: %d: %s", resp.StatusCode, resp.Body)
	}
	result2 := parseJSON(resp.Body)
	t2 := result2["tourism"].(map[string]interface{})
	t2ID := t2["id"].(float64)
	t2Slug := t2["slug"].(string)

	if t2Slug == "slug-dedup-e2e" {
		t.Errorf("expected deduplicated slug, but got same slug: %s", t2Slug)
	}

	// Cleanup
	doRequest("DELETE", "/api/tourisms/"+fmt.Sprintf("%d", int(t1ID)), nil, authToken)
	doRequest("DELETE", "/api/tourisms/"+fmt.Sprintf("%d", int(t2ID)), nil, authToken)
}
