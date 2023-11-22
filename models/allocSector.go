package models

type AllocSector struct {
	AllocSectorKey uint8  `db:"alloc_sector_key" json:"alloc_sector_key"`
	ProductKey     uint64 `db:"product_key" json:"product_key"`
	PeriodeKey     uint64 `db:"periode_key" json:"periode_key"`
	SectorKey      uint64 `db:"sector_key" json:"sector_key"`
	SectorValue    uint64 `db:"sector_value" json:"sector_value"`
}
