package product

type CreateProductReq struct {
	Title        string  `json:"title"`
	Subtitle     *string `json:"subtitle"`
	Description  *string `json:"description"`
	Handle       string  `json:"handle"`
	Thumbnail    *string `json:"thumbnail"`
	Status       string  `json:"status"`
	CollectionID *string `json:"collection_id"`
}

type CreateVariantReq struct {
	Title           string  `json:"title"`
	SKU             *string `json:"sku"`
	ManageInventory bool    `json:"manage_inventory"`
	AllowBackorder  bool    `json:"allow_backorder"`
}
