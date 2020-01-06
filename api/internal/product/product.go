package product

//Product information
type Product struct {
	ProductName string
	Ingredients []string
	Serving     string
	Nutrition   map[string]float32
	Version     int
	ID          string `bson:"_id"`
	Trust       map[string]points
	Changed     string
	Error       string
	Changes     []Product
}
type points struct {
	User    string
	Confirm int
	Deny    int
	Version int
}
