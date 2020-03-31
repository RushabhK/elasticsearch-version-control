package model

import "strconv"

type IndexDetails struct {
	Version     int
	Alias       string
	Script      string
	IndexConfig string
}

func (indexDetails IndexDetails) GetName() string {
	return indexDetails.Alias + "_v" + strconv.Itoa(indexDetails.Version)
}
