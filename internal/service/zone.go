package service

import (
	"fmt"

	"github.com/Shakkuuu/ed-mist-backend/pkg/mistapi"
)

// ZoneService ゾーンサービス
type ZoneService struct {
	mistClient *mistapi.Client
}

// NewZoneService ゾーンサービスを作成
func NewZoneService(mistClient *mistapi.Client) *ZoneService {
	return &ZoneService{
		mistClient: mistClient,
	}
}

// GetAll 全ゾーンを取得
func (s *ZoneService) GetAll() ([]mistapi.MistZone, error) {
	if s.mistClient == nil {
		return nil, fmt.Errorf("mist APIクライアントが初期化されていません")
	}

	// Mist APIからゾーン一覧を取得
	mistZones, err := s.mistClient.GetZones(s.mistClient.SiteID)
	if err != nil {
		return nil, fmt.Errorf("mist APIゾーン取得エラー: %w", err)
	}

	return mistZones, nil
}

// GetClientsByZoneID ゾーンIDでクライアント一覧を取得
func (s *ZoneService) GetClientsByZoneID(zoneID string) (*ZoneClientsResponse, error) {
	if s.mistClient == nil {
		return nil, fmt.Errorf("mist APIクライアントが初期化されていません")
	}

	sdkClientIDs, clientIDs, err := s.mistClient.GetZoneClients(s.mistClient.SiteID, zoneID)
	if err != nil {
		return nil, fmt.Errorf("mist APIゾーン取得エラー: %w", err)
	}

	var sdkClients []interface{}
	if len(sdkClientIDs) > 0 {
		for _, id := range sdkClientIDs {
			sc, err := s.mistClient.GetSDKClient(s.mistClient.SiteID, id)
			if err != nil {
				return nil, err
			}
			sdkClients = append(sdkClients, sc)
		}
	}

	var wirelessClients []interface{}
	if len(clientIDs) > 0 {
		for _, id := range clientIDs {
			wc, err := s.mistClient.GetWirelessClient(s.mistClient.SiteID, id)
			if err != nil {
				return nil, err
			}
			wirelessClients = append(wirelessClients, wc)
		}
	}

	return &ZoneClientsResponse{
		SDKClients: sdkClients,
		Clients:    wirelessClients,
	}, nil
}

// ZoneClientsResponse ゾーンのクライアントレスポンス
type ZoneClientsResponse struct {
	SDKClients []interface{} `json:"sdkclients"`
	Clients    []interface{} `json:"clients"`
}
