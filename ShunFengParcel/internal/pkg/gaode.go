package pkg

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/go-kratos/kratos/v2/log"
)

const matrixURL = "https://restapi.amap.com/v3/distance"

// Client 高德客户端
type Client struct {
	key     string
	client  *http.Client
	timeout time.Duration
	log     *log.Helper
}

// NewClient 通过 Kratos 容器注入
func NewClient(c *http.Client, key string, timeout time.Duration, logger log.Logger) *Client {
	return &Client{
		key:     key,
		client:  c,
		timeout: timeout,
		log:     log.NewHelper(logger),
	}
}

// Pair 结果对
type Pair struct {
	OrderID  int64
	Distance int // 米
	Duration int // 秒
}

// BatchDistance 批量骑行距离+时长
func (c *Client) BatchDistance(ctx context.Context, riderLng, riderLat float64, orders []OrderCoord) ([]Pair, error) {
	origins := fmt.Sprintf("%.6f,%.6f", riderLng, riderLat)
	dst := make([]string, len(orders))
	for i, o := range orders {
		dst[i] = fmt.Sprintf("%.6f,%.6f", o.Lng, o.Lat)
	}
	val := url.Values{}
	val.Set("key", c.key)
	val.Set("origins", origins)
	val.Set("destinations", strings.Join(dst, "|"))
	val.Set("type", "3") // 骑行

	reqURL := matrixURL + "?" + val.Encode()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

    // 兼容旧版 Go：使用 NewRequest + WithContext
    req, err := http.NewRequest("GET", reqURL, nil)
    if err != nil {
        return nil, fmt.Errorf("gaode matrix build request err: %w", err)
    }
    req = req.WithContext(ctx)
    resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gaode matrix http err: %w", err)
	}
	defer resp.Body.Close()

	var body matrixResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("gaode matrix json err: %w", err)
	}
	if body.Status != "1" {
		return nil, fmt.Errorf("gaode matrix api fail, info=%s", body.Info)
	}

	pairs := make([]Pair, len(orders))
	for i, d := range body.Distance {
		dist, _ := strconv.Atoi(d.Distance)
		dur, _ := strconv.Atoi(d.Duration)
		pairs[i] = Pair{OrderID: orders[i].ID, Distance: dist, Duration: dur}
	}
	c.log.Infof("gaode matrix success, count=%d", len(pairs))
	return pairs, nil
}

// ---- 内部类型 ----
type OrderCoord struct {
	ID  int64
	Lng float64
	Lat float64
}

type matrixResp struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Distance []struct {
		Distance string `json:"distance"`
		Duration string `json:"duration"`
	} `json:"distance"`
}
