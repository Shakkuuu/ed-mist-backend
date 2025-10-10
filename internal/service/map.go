package service

import (
	"fmt"

	"github.com/Shakkuuu/ed-mist-backend/pkg/mistapi"
)

// MapService マップサービス
type MapService struct {
	mistClient *mistapi.Client
}

// NewMapService マップサービスを作成
func NewMapService(mistClient *mistapi.Client) *MapService {
	return &MapService{
		mistClient: mistClient,
	}
}

// GetAll 全マップを取得
func (s *MapService) GetAll() ([]mistapi.MistMap, error) {
	if s.mistClient == nil {
		return nil, fmt.Errorf("mist APIクライアントが初期化されていません")
	}

	mistMaps, err := s.mistClient.GetMaps(s.mistClient.SiteID)
	if err != nil {
		return nil, fmt.Errorf("マップ取得エラー: %w", err)
	}

	return mistMaps, nil
}

// GetImage マップ画像を取得
func (s *MapService) GetImage(mapID string) ([]byte, error) {
	if s.mistClient == nil {
		return nil, fmt.Errorf("mist APIクライアントが初期化されていません")
	}

	imageData, err := s.mistClient.GetMapImage(s.mistClient.SiteID, mapID)
	if err != nil {
		return nil, fmt.Errorf("マップ画像取得エラー: %w", err)
	}

	return imageData, nil
}
