package database

type Item struct {
	Id          string `json:"id"`
	Name        string `json:"title"`
	Description string `json:"description"`
	ImageUri    string `json:"imageUri"`
	VideoUri    string `json:"videoUri"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	ImageUri string `json:"imageUri"`
}

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
