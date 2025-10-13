package tmx

type TrackmaniaExchangeAPI interface {
	GetMapInfo(mapID string) (*MapInfo, error)
}
