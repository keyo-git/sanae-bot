package bot

// Cookie ...
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Config stores Sanae initialization parameters
type Config struct {
	Token   string   `json:"token"`
	ConnStr string   `json:"conn_str"`
	Cookies []Cookie `json:"cookies"`
}
