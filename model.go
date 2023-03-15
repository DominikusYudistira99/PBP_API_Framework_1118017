package PBP_API_Framework_1118017

import (
	"encoding/xml"
)

type Album struct {
	XMLName xml.Name `json:"-" xml:"album"`
	Id      int      `json:"id" xml:"id,attr"`
	Band    string   `json:"band" xml:"band"`
	Title   string   `json:"title" xml:"title"`
	Year    int      `json:"year" xml:"year"`
}