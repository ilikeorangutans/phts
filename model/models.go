package model

type Photo struct {
	//Timestamps
	Record
	Renditions   []Rendition
	CollectionID int64  `db:"collection_id"`
	Description  string `db:"description"`
}

type Rendition struct {
	//Timestamps
	Record
	Original bool
	PhotoID  int64
	Data     []byte
}
