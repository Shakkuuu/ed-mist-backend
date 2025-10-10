package mistapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIToken   string
	SiteID     string
	HTTPClient *http.Client
}

// 新しいMist APIクライアントを作成
func NewClient(baseURL, apiToken, siteID string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIToken:   apiToken,
		SiteID:     siteID,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Mistのマップ情報
type MistMap struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	WidthM       float64 `json:"width_m"`
	HeightM      float64 `json:"height_m"`
	PPM          float64 `json:"ppm"`
	URL          string  `json:"url"`
	ThumbnailURL string  `json:"thumbnail_url"`
}

type Vertices struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Mistのゾーン情報
type MistZone struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Vertices   []Vertices `json:"vertices"`
	VerticesM  []Vertices `json:"vertices_m"`
	MapID      string     `json:"map_id"`
	Clients    []string   `json:"clients"`
	SDKClients []string   `json:"sdkclients"`
}

// Mistのクライアント情報
type MistSDKClient struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	UUID     string  `json:"uuid"`
	LastSeen float64 `json:"last_seen"`
}
type MistWirelessClient struct {
	Annotation string  `json:"annotation"`
	HostName   string  `json:"hostname"`
	Mac        string  `json:"mac"`
	OS         string  `json:"os"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	XM         float64 `json:"x_m"`
	YM         float64 `json:"y_m"`
	LastSeen   float64 `json:"last_seen"`
}

// BLEデバイス情報
// type MistBLEDevice struct {
// 	Mac         string  `json:"mac"`
// 	CurrSite    string  `json:"curr_site"`
// 	Timestamp   int64   `json:"_timestamp"`
// 	Manufacture string  `json:"manufacture"`
// 	LastSeen    float64 `json:"last_seen"`
// 	MapID       string  `json:"map_id"`
// 	TTL         int     `json:"_ttl"`
// 	ID          string  `json:"_id"`
// 	X           float64 `json:"x,omitempty"`
// 	Y           float64 `json:"y,omitempty"`
// }

type MistBLEDevice struct {
	Manufacture string  `json:"manufacture"`
	Mac         string  `json:"mac"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	XM          float64 `json:"x_m"`
	YM          float64 `json:"y_m"`
	LastSeen    float64 `json:"last_seen"`
}

// Mistのマップ一覧を取得
func (c *Client) GetMaps(siteID string) ([]MistMap, error) {
	url := fmt.Sprintf("%s/api/v1/sites/%s/maps", c.BaseURL, siteID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.SetAuthHeader(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mist APIマップ取得エラー: %s, %s", resp.Status, string(body))
	}

	// デバッグ用：レスポンスボディをログ出力
	body, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Mist API マップ取得レスポンス: %s\n", string(body))

	var maps []MistMap
	if err := json.Unmarshal(body, &maps); err != nil {
		return nil, err
	}
	return maps, nil
}

// Mistのゾーン一覧を取得
func (c *Client) GetZones(siteID string) ([]MistZone, error) {
	url := fmt.Sprintf("%s/api/v1/sites/%s/zones", c.BaseURL, siteID)
	// fmt.Printf("Mist API ゾーン取得URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.SetAuthHeader(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()

	// レスポンスボディを読み取り
	body, _ := io.ReadAll(resp.Body)
	// fmt.Printf("mist API ゾーン取得レスポンス: Status=%s, Body=%s\n", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mist APIゾーン取得エラー: %s, %s", resp.Status, string(body))
	}

	// まず配列として解析を試みる
	var zones []MistZone
	if err := json.Unmarshal(body, &zones); err != nil {
		// 配列として解析できない場合、オブジェクトとして解析を試みる
		var responseObj map[string]any
		if err := json.Unmarshal(body, &responseObj); err != nil {
			return nil, fmt.Errorf("ゾーンデータのJSON解析エラー: %w, Body: %s", err, string(body))
		}

		// オブジェクト内のzonesフィールドを探す
		if zonesData, ok := responseObj["zones"]; ok {
			zonesJSON, err := json.Marshal(zonesData)
			if err != nil {
				return nil, fmt.Errorf("ゾーンデータの変換エラー: %w", err)
			}
			if err := json.Unmarshal(zonesJSON, &zones); err != nil {
				return nil, fmt.Errorf("ゾーンデータのJSON解析エラー: %w, Body: %s", err, string(body))
			}
		} else {
			return nil, fmt.Errorf("レスポンスにzonesフィールドが見つかりません: %s", string(body))
		}
	}

	log.Printf("取得したゾーン数: %d\n", len(zones))
	return zones, nil
}

// マップ画像を取得
func (c *Client) GetMapImage(siteID, mapID string) ([]byte, error) {
	// まずマップ一覧を取得して、指定されたマップのURLを取得
	maps, err := c.GetMaps(siteID)
	if err != nil {
		return nil, fmt.Errorf("マップ一覧取得エラー: %w", err)
	}

	var targetMap *MistMap
	for _, m := range maps {
		if m.ID == mapID {
			targetMap = &m
			break
		}
	}

	if targetMap == nil {
		return nil, fmt.Errorf("マップが見つかりません: %s", mapID)
	}

	// URLフィールドから直接画像を取得
	if targetMap.URL == "" {
		return nil, fmt.Errorf("マップ画像URLが利用できません")
	}

	req, err := http.NewRequest("GET", targetMap.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("マップ画像取得エラー: %s, %s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}

// 認証ヘッダーをセット
func (c *Client) SetAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Token "+c.APIToken)
}

// Mistのゾーン内クライアント一覧を取得
func (c *Client) GetZoneClients(siteID, zoneID string) (sdkClients []string, Clients []string, err error) {
	url := fmt.Sprintf("%s/api/v1/sites/%s/stats/zones/%s", c.BaseURL, siteID, zoneID)
	// fmt.Printf("Mist API ゾーン統計取得URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	c.SetAuthHeader(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()

	// レスポンスボディを読み取り
	body, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Mist API ゾーン統計取得レスポンス: Status=%s, Body=%s\n", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		// 404エラーの場合は空配列を返す（クライアントが存在しない）
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("ゾーン統計が見つかりません（404）: %s\n", string(body))
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("mist APIゾーン統計取得エラー: %s, %s", resp.Status, string(body))
	}

	var zone MistZone
	if err := json.Unmarshal(body, &zone); err != nil {
		return nil, nil, fmt.Errorf("レスポンスにzonesフィールドが見つかりません: %s", string(body))
	}

	// fmt.Printf("取得したゾーンクライアント数(sdk): %d\n", len(zone.SDKClients))
	return zone.SDKClients, zone.Clients, nil
}

// SDKクライアント詳細取得
func (c *Client) GetSDKClient(siteID, clientID string) (*MistSDKClient, error) {
	url := fmt.Sprintf("%s/api/v1/sites/%s/stats/sdkclients/%s", c.BaseURL, siteID, clientID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.SetAuthHeader(req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mist API SDKクライアント取得エラー: %s, %s", resp.Status, string(body))
	}
	var client MistSDKClient
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, err
	}
	return &client, nil
}

// WiFiクライアント詳細取得
func (c *Client) GetWirelessClient(siteID, clientID string) (*MistWirelessClient, error) {
	url := fmt.Sprintf("%s/api/v1/sites/%s/stats/clients/%s", c.BaseURL, siteID, clientID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.SetAuthHeader(req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mist API WiFiクライアント取得エラー: %s, %s", resp.Status, string(body))
	}
	var client MistWirelessClient
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, err
	}
	return &client, nil
}

// 新しいBLEデバイス一覧を取得
func (c *Client) GetBLEDevices(siteID, mapID string) ([]MistBLEDevice, error) {
	url := fmt.Sprintf("%s/api/v1/sites/%s/stats/maps/%s/discovered_assets", c.BaseURL, siteID, mapID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.SetAuthHeader(req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close error: %v\n", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mist API BLEデバイス取得エラー: %s, %s", resp.Status, string(body))
	}
	var devices []MistBLEDevice
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

// ゾーン内のBLEデバイスを抽出
func (c *Client) GetZoneBLEDevices(siteID, zoneID, mapID string) ([]MistBLEDevice, error) {
	// まずゾーン情報を取得
	zones, err := c.GetZones(siteID)
	if err != nil {
		return nil, fmt.Errorf("ゾーン情報取得エラー: %w", err)
	}

	// 指定されたゾーンを探す
	var targetZone *MistZone
	for _, zone := range zones {
		if zone.ID == zoneID {
			targetZone = &zone
			break
		}
	}

	if targetZone == nil {
		return nil, fmt.Errorf("ゾーンが見つかりません: %s", zoneID)
	}

	// BLEデバイス一覧を取得
	bleDevices, err := c.GetBLEDevices(siteID, mapID)
	if err != nil {
		return nil, fmt.Errorf("BLEデバイス取得エラー: %w", err)
	}

	// ゾーン内のBLEデバイスを抽出
	var zoneBLEDevices []MistBLEDevice
	for _, device := range bleDevices {
		if isPointInZone(device.X, device.Y, targetZone.Vertices) {
			zoneBLEDevices = append(zoneBLEDevices, device)
		}
	}

	return zoneBLEDevices, nil
}

// 点がゾーン内にあるかどうかを判定（ポリゴン内判定）
func isPointInZone(x, y float64, vertices []Vertices) bool {
	if len(vertices) < 3 {
		return false
	}

	inside := false
	j := len(vertices) - 1

	for i := range len(vertices) {
		if ((vertices[i].Y > y) != (vertices[j].Y > y)) &&
			(x < (vertices[j].X-vertices[i].X)*(y-vertices[i].Y)/(vertices[j].Y-vertices[i].Y)+vertices[i].X) {
			inside = !inside
		}
		j = i
	}

	return inside
}
