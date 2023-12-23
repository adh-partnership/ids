package weather

type Weather struct {
	METAR string `json:"metar"`
	TAF   string `json:"taf"`
}

type response struct {
	METARs []METAR `xml:"data>METAR"`
	TAFs   []TAF   `xml:"data>TAF"`
}

type METAR struct {
	StationID string `xml:"station_id"`
	RawText   string `xml:"raw_text"`
}

type TAF struct {
	StationID string `xml:"station_id"`
	RawText   string `xml:"raw_text"`
}
