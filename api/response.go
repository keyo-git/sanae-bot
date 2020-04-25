package api

type torrent struct {
	Hash  string
	Added string
	Name  string
	TSize string
	FSize string
}

type GalleryMetadata struct {
	Gid          int       `json:"gid"`
	Token        string    `json:"token"`
	ArchiverKey  string    `json:"archiver_key"`
	Title        string    `json:"title"`
	TitleJPN     string    `json:"title_jpn"`
	Category     string    `json:"category"`
	Thumb        string    `json:"thumb"`
	Uploader     string    `json:"uploader"`
	Posted       int64     `json:"posted,string"`
	FileCount    int       `json:"filecount,string"`
	FileSize     int       `json:"filesize"`
	Expunged     bool      `json:"expunged"`
	Rating       float32   `json:"rating,string"`
	TorrentCount string    `json:"torrentcount"`
	Torrents     []torrent `json:"torrents"`
	Tags         []string  `json:"tags"`
	Error        string    `json:"error"`
}

type apiResponse struct {
	GMetaData []GalleryMetadata
}
