package e2e

import (
	"fmt"
	"kebondowo/models"
	"testing"
)

// --- Role CRUD ---

func TestRoleGetAllPublic(t *testing.T) {
	resp := doRequest("GET", "/api/roles/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	roles := result["roles"].([]interface{})
	if len(roles) < 2 {
		t.Errorf("expected at least 2 roles (admin, user), got %d", len(roles))
	}
}

func TestRoleCRUDRequiresAuth(t *testing.T) {
	body := map[string]string{"name": "moderator"}
	resp := doRequest("POST", "/api/roles/", body, "")
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 for creating role without auth, got %d", resp.StatusCode)
	}
}

func TestRoleCRUD(t *testing.T) {
	// Create
	body := map[string]string{"name": "testrole"}
	resp := doRequest("POST", "/api/roles/", body, authToken)
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	role := result["role"].(map[string]interface{})
	roleID := role["ID"].(float64)

	// Get one
	resp = doRequest("GET", fmt.Sprintf("/api/roles/%d", int(roleID)), nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d", resp.StatusCode)
	}

	// Update
	updateBody := map[string]string{"name": "testrole_updated"}
	resp = doRequest("PUT", fmt.Sprintf("/api/roles/%d", int(roleID)), updateBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Update, got %d: %s", resp.StatusCode, resp.Body)
	}
	result = parseJSON(resp.Body)
	role = result["role"].(map[string]interface{})
	if role["Name"] != "testrole_updated" {
		t.Errorf("expected updated name, got %v", role["Name"])
	}

	// Delete
	resp = doRequest("DELETE", fmt.Sprintf("/api/roles/%d", int(roleID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}
}

// --- Gallery CRUD ---

func TestGalleryGetAllPublic(t *testing.T) {
	resp := doRequest("GET", "/api/galleries/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestGalleryCreateRequiresAuth(t *testing.T) {
	body := map[string]string{
		"picture_url": "https://example.com/pic.jpg",
		"caption":     "test",
	}
	resp := doRequest("POST", "/api/galleries/", body, "")
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestGalleryCRUD(t *testing.T) {
	// Create
	body := map[string]string{
		"picture_url": "https://example.com/gallery.jpg",
		"caption":     "Test Gallery",
	}
	resp := doRequest("POST", "/api/galleries/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	gallery := result["gallery"].(map[string]interface{})
	galleryID := gallery["id"].(float64)

	// Get one
	resp = doRequest("GET", fmt.Sprintf("/api/galleries/%d", int(galleryID)), nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d", resp.StatusCode)
	}

	// Get all with pagination
	resp = doRequest("GET", "/api/galleries/?limit=5&offset=0&sortby=created_at&orderedby=desc", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetAll with params, got %d", resp.StatusCode)
	}
	result = parseJSON(resp.Body)
	meta := result["meta"].(map[string]interface{})
	if meta["total"].(float64) < 1 {
		t.Errorf("expected total >= 1")
	}

	// Delete
	resp = doRequest("DELETE", fmt.Sprintf("/api/galleries/%d", int(galleryID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}
}

// --- Article CRUD ---

func TestArticleGetAllPublic(t *testing.T) {
	resp := doRequest("GET", "/api/articles/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestArticleCRUD(t *testing.T) {
	// Create
	body := map[string]string{
		"title":       "Test Article",
		"slug":        "test-article-e2e",
		"author":      "Tester",
		"content":     "This is test content for E2E testing.",
		"picture_url": "https://example.com/article.jpg",
	}
	resp := doRequest("POST", "/api/articles/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	article := result["article"].(map[string]interface{})
	articleID := article["id"].(float64)

	// Get one by slug
	resp = doRequest("GET", "/api/articles/test-article-e2e", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Update
	updateBody := map[string]string{
		"title":   "Updated Article",
		"slug":    "test-article-e2e",
		"author":  "Tester",
		"content": "Updated content",
	}
	resp = doRequest("PUT", fmt.Sprintf("/api/articles/%d", int(articleID)), updateBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Update, got %d: %s", resp.StatusCode, resp.Body)
	}
	result = parseJSON(resp.Body)
	updatedArticle := result["article"].(map[string]interface{})
	if updatedArticle["title"] != "Updated Article" {
		t.Errorf("expected title 'Updated Article', got %v", updatedArticle["title"])
	}

	// Delete
	resp = doRequest("DELETE", fmt.Sprintf("/api/articles/%d", int(articleID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}
}

// --- Tourism CRUD ---

func TestTourismGetAllPublic(t *testing.T) {
	resp := doRequest("GET", "/api/tourisms/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestTourismCRUD(t *testing.T) {
	// Create
	body := map[string]interface{}{
		"title":             "Test Tourism",
		"slug":              "test-tourism-e2e",
		"description":       "A beautiful place",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/tourism.jpg",
		"pictures": []map[string]string{
			{"picture_url": "https://example.com/pic1.jpg"},
			{"picture_url": "https://example.com/pic2.jpg"},
		},
	}
	resp := doRequest("POST", "/api/tourisms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	tourism := result["tourism"].(map[string]interface{})
	tourismID := tourism["id"].(float64)

	// Get one by slug
	resp = doRequest("GET", "/api/tourisms/test-tourism-e2e", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Update
	updateBody := map[string]interface{}{
		"title":       "Updated Tourism",
		"slug":        "test-tourism-e2e",
		"description": "An updated place",
		"latitude":    -7.6,
		"longitude":   110.5,
	}
	resp = doRequest("PUT", fmt.Sprintf("/api/tourisms/%d", int(tourismID)), updateBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Update, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Delete (should cascade delete pictures)
	resp = doRequest("DELETE", fmt.Sprintf("/api/tourisms/%d", int(tourismID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Verify pictures were cascaded
	var count int64
	testDB.Model(&models.TourismPicture{}).Where("tourism_id = ?", int(tourismID)).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 tourism pictures after cascade delete, got %d", count)
	}
}

// --- UMKM CRUD ---

func TestUMKMGetAllPublic(t *testing.T) {
	resp := doRequest("GET", "/api/umkms/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestUMKMCRUD(t *testing.T) {
	// Create
	body := map[string]interface{}{
		"title":             "Test UMKM",
		"slug":              "test-umkm-e2e",
		"description":       "A test UMKM",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/umkm.jpg",
		"contact":           "081234567890",
		"contact_name":      "Test Contact",
		"pictures": []map[string]string{
			{"picture_url": "https://example.com/umkm1.jpg"},
		},
	}
	resp := doRequest("POST", "/api/umkms/", body, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	umkm := result["umkm"].(map[string]interface{})
	umkmID := umkm["id"].(float64)

	// Get one by slug
	resp = doRequest("GET", "/api/umkms/test-umkm-e2e", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Update
	updateBody := map[string]interface{}{
		"title":        "Updated UMKM",
		"slug":         "test-umkm-e2e",
		"description":  "Updated desc",
		"latitude":     -7.6,
		"longitude":    110.5,
		"contact":      "081234567890",
		"contact_name": "Updated Contact",
	}
	resp = doRequest("PUT", fmt.Sprintf("/api/umkms/%d", int(umkmID)), updateBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Update, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Delete (should cascade delete pictures)
	resp = doRequest("DELETE", fmt.Sprintf("/api/umkms/%d", int(umkmID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Verify pictures were cascaded
	var count int64
	testDB.Model(&models.UMKMPicture{}).Where("umkm_id = ?", int(umkmID)).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 UMKM pictures after cascade delete, got %d", count)
	}
}

// --- Tourism Picture CRUD ---

func TestTourismPictureCRUD(t *testing.T) {
	// Create a tourism first
	tourismBody := map[string]interface{}{
		"title":             "Tourism For Pic",
		"slug":              "tourism-for-pic-e2e",
		"description":       "For picture test",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/cover.jpg",
	}
	resp := doRequest("POST", "/api/tourisms/", tourismBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("setup: failed to create tourism: %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	tourism := result["tourism"].(map[string]interface{})
	tourismID := uint32(tourism["id"].(float64))

	// Create tourism picture
	picBody := map[string]interface{}{
		"picture_url": "https://example.com/tpic.jpg",
		"caption":     "Nice view",
		"tourism_id":  tourismID,
	}
	resp = doRequest("POST", "/api/tourism-pictures/", picBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result = parseJSON(resp.Body)
	pic := result["tourism_picture"].(map[string]interface{})
	picID := pic["id"].(float64)

	// Get all
	resp = doRequest("GET", "/api/tourism-pictures/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetAll, got %d", resp.StatusCode)
	}

	// Get one
	resp = doRequest("GET", fmt.Sprintf("/api/tourism-pictures/%d", int(picID)), nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d", resp.StatusCode)
	}

	// Delete picture
	resp = doRequest("DELETE", fmt.Sprintf("/api/tourism-pictures/%d", int(picID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Cleanup tourism
	doRequest("DELETE", fmt.Sprintf("/api/tourisms/%d", int(tourismID)), nil, authToken)
}

// --- UMKM Picture CRUD ---

func TestUMKMPictureCRUD(t *testing.T) {
	// Create a UMKM first
	umkmBody := map[string]interface{}{
		"title":             "UMKM For Pic",
		"slug":              "umkm-for-pic-e2e",
		"description":       "For picture test",
		"latitude":          -7.5,
		"longitude":         110.4,
		"cover_picture_url": "https://example.com/cover.jpg",
		"contact":           "08123",
		"contact_name":      "Tester",
	}
	resp := doRequest("POST", "/api/umkms/", umkmBody, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("setup: failed to create UMKM: %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	umkm := result["umkm"].(map[string]interface{})
	umkmID := uint32(umkm["id"].(float64))

	// Create UMKM picture
	picBody := map[string]interface{}{
		"picture_url": "https://example.com/upic.jpg",
		"caption":     "Product shot",
		"umkm_id":     umkmID,
	}
	resp = doRequest("POST", "/api/umkm-pictures/", picBody, authToken)
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, resp.Body)
	}
	result = parseJSON(resp.Body)
	pic := result["umkm_picture"].(map[string]interface{})
	picID := pic["id"].(float64)

	// Get all
	resp = doRequest("GET", "/api/umkm-pictures/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetAll, got %d", resp.StatusCode)
	}

	// Get one
	resp = doRequest("GET", fmt.Sprintf("/api/umkm-pictures/%d", int(picID)), nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GetOne, got %d", resp.StatusCode)
	}

	// Delete picture
	resp = doRequest("DELETE", fmt.Sprintf("/api/umkm-pictures/%d", int(picID)), nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on Delete, got %d: %s", resp.StatusCode, resp.Body)
	}

	// Cleanup UMKM
	doRequest("DELETE", fmt.Sprintf("/api/umkms/%d", int(umkmID)), nil, authToken)
}

// --- Dashboard ---

func TestDashboard(t *testing.T) {
	resp := doRequest("GET", "/api/dashboard/", nil, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	// Dashboard returns a valid JSON response.
	// Fields with zero values may be omitted due to omitempty tags,
	// so we just verify the response is parseable and the status is 200.
	result := parseJSON(resp.Body)
	if result == nil {
		t.Error("expected parseable JSON response from dashboard")
	}
}
