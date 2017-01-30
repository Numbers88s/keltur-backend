package models

// Item is each of the item.
type Item struct {
	Id       string
	Title    string
	Vendor   string
	ImageURL string
	Category string
	Price    float32
}

func (db *DB) AllItems() ([]*Item, error) {
    rows, err := db.Query("SELECT * FROM items")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := make([]*Item, 0)
    for rows.Next() {
        item := new(Item)
        err := rows.Scan(&item.Id, &item.Title, &item.Vendor, &item.ImageURL, &item.Category, &item.Price)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return items, nil
}