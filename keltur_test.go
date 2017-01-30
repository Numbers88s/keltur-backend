package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/Numbers88s/keltur-backend/models"
)

type mockDB struct{}

func (mdb *mockDB) AllItems() ([]*models.Item, error) {
    items := make([]*models.Item, 0)
    items = append(items, &models.Item{"978-1503261969", "Emma", "Crate&Barrel", "/assets/images/rug.jpg", "carpet", 755.95})
    items = append(items, &models.Item{"978-1503261988", "Platinum rug", "William Sonoma", "/assets/images/silver_rug.jpg", "carpet", 1055.95})
    return items, nil
}

func TestBooksIndex(t *testing.T) {
    rec := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/items", nil)

    env := Env{db: &mockDB{}}
    http.HandlerFunc(env.itemsIndex).ServeHTTP(rec, req)

    expected := "978-1503261969, Emma, Crate&Barrel, /assets/images/rug.jpg, carpet, $755.95\n978-1503261988, Platinum rug, William Sonoma, /assets/images/silver_rug.jpg, carpet, $1055.95\n"
    if expected != rec.Body.String() {
        t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
    }
}