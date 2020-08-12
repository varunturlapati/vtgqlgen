package db

type Fruit struct {
	Id int
	Name string
	Quantity int
	Color string
	Level string
	X Detail
	Rack ServerRack
}

type Detail struct {
	Id int64
	Name string
	Color string
	Taste string
}

type Level struct {
	Color string
	Level string
}

type ServerRack struct { // Yes! Trying to demo different datasources that may have relevant data of an object
	Id           int64
	Name         string
	CustomFields CustomFields `json:"custom_fields"`
	Created      string
}
type CustomFields struct {
	RblxRackId     int64  `json:"rblx_rack_id"`
	DesignRevision string `json:"design_revision"`
	CageId         string `json:"cage_id"`
}
