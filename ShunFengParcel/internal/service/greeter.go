package service

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "math/rand"
    "net/http"
    "net/url"
    "sort"
    "strings"
    "time"
    "unicode/utf8"

	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/basic/config"
	"ShunFengParcel/internal/biz"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// -------- 轨迹存储键 --------
func trajKey(courierID int64) string { return fmt.Sprintf("traj:%d", courierID) }
func lastKey(courierID int64) string { return fmt.Sprintf("last:%d", courierID) }

// -------- 纠偏：WGS84 -> GCJ-02（中国境内） --------
// 参考实现（简化版）：仅当点位处于中国范围时应用偏移
func outOfChina(lat, lng float64) bool {
	return lng < 72.004 || lng > 137.8347 || lat < 0.8293 || lat > 55.8271
}
func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*math.Pi) + 320.0*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}
func transformLng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}
func wgs84ToGcj02(lat, lng float64) (float64, float64) {
	if outOfChina(lat, lng) {
		return lat, lng
	}
	a := 6378245.0
	ee := 0.00669342162296594323
	dLat := transformLat(lng-105.0, lat-35.0)
	dLng := transformLng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * math.Pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * math.Pi)
	dLng = (dLng * 180.0) / (a / sqrtMagic * math.Cos(radLat) * math.Pi)
	mgLat := lat + dLat
	mgLng := lng + dLng
	return mgLat, mgLng
}

// -------- 平滑：指数平滑与速度门限 --------
type lastPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Ts  int64   `json:"ts"` // 秒级时间戳
}

// 对新点进行指数平滑；若速度过大（疑似跳点），采用前一点坐标
func smoothPoint(prev *lastPoint, lat, lng float64, ts int64) (float64, float64) {
	if prev == nil || prev.Ts <= 0 {
		return lat, lng
	}
	// 估算速度（米/秒），超过 60 m/s 则判定跳点
	dist := haversine(prev.Lat, prev.Lng, lat, lng)
	dt := float64(ts - prev.Ts)
	if dt <= 0 {
		dt = 1
	}
	speed := dist / dt
	if speed > 60 { // 约 216 km/h
		return prev.Lat, prev.Lng
	}
	alpha := 0.6
	return alpha*lat + (1-alpha)*prev.Lat, alpha*lng + (1-alpha)*prev.Lng
}

// 球面距离（米）
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// -------- 读取快递员ID（X-Courier-ID） --------
func courierIDFromCtx(ctx context.Context, fallback int64) int64 {
	if fallback > 0 {
		return fallback
	}
	if tr, ok := transport.FromServerContext(ctx); ok {
		if h := tr.RequestHeader().Get("X-Courier-ID"); h != "" {
			var id int64
			_, _ = fmt.Sscan(h, &id)
			if id > 0 {
				return id
			}
		}
	}
	if v := ctx.Value("courier_id"); v != nil {
		if id, ok := v.(int64); ok && id > 0 {
			return id
		}
	}
	return 0
}

// -------- 快递员位置上报 --------
func (s *GreeterService) ReportCourierLocation(ctx context.Context, req *v1.ReportCourierLocationRequest) (*v1.ReportCourierLocationReply, error) {
    courierID := courierIDFromCtx(ctx, req.CourierId)
    if courierID == 0 {
        return nil, status.Error(codes.InvalidArgument, "缺少快递员ID")
    }

    // 纠偏：假定客户端上报 WGS84，转换到 GCJ-02（中国境内）；境外保持原值
    lat, lng := wgs84ToGcj02(req.Lat, req.Lng)
    // 请求体无时间戳字段，统一以服务器当前时间
    ts := timeNow()

    // 调试/联调开关：通过请求头控制是否跳过平滑或重置上一点
    var ignoreSmooth, resetLast bool
    if tr, ok := transport.FromServerContext(ctx); ok {
        if strings.EqualFold(tr.RequestHeader().Get("X-Ignore-Smooth"), "true") {
            ignoreSmooth = true
        }
        if strings.EqualFold(tr.RequestHeader().Get("X-Reset-Last"), "true") {
            resetLast = true
        }
    }

    // 获取上一点用于平滑
    var prev *lastPoint
    if config.RDB != nil {
        if !resetLast {
            if val, err := config.RDB.Get(ctx, lastKey(courierID)).Result(); err == nil && val != "" {
                var p lastPoint
                if json.Unmarshal([]byte(val), &p) == nil {
                    prev = &p
                }
            }
        }
    }
    // 指数平滑与速度门限
    latSm, lngSm := lat, lng
    if !ignoreSmooth {
        latSm, lngSm = smoothPoint(prev, lat, lng, ts)
    }

    // 存储：last + 有序集合轨迹
    if config.RDB != nil {
        b, _ := json.Marshal(lastPoint{Lat: latSm, Lng: lngSm, Ts: ts})
        pipe := config.RDB.Pipeline()
        pipe.Set(ctx, lastKey(courierID), string(b), 0)
        pipe.ZAdd(ctx, trajKey(courierID), redis.Z{Score: float64(ts), Member: fmt.Sprintf("%f,%f,%d", latSm, lngSm, ts)})
        // GEO 实时位置：用于附近搜索/实时展示
        pipe.GeoAdd(ctx, "geo:couriers", &redis.GeoLocation{Longitude: lngSm, Latitude: latSm, Name: fmt.Sprintf("%d", courierID)})
        // 仅保留最近 30 分钟的轨迹，避免集合无限膨胀
        pipe.ZRemRangeByScore(ctx, trajKey(courierID), "-inf", fmt.Sprintf("%f", float64(ts-1800)))
        _, _ = pipe.Exec(ctx)

        // 即时广播：按订单通道推送位置（X-Order-ID）
        if tr, ok := transport.FromServerContext(ctx); ok {
            orderID := tr.RequestHeader().Get("X-Order-ID")
            if orderID != "" {
                msg := map[string]interface{}{
                    "order_id":   orderID,
                    "courier_id": courierID,
                    "lat":        latSm,
                    "lng":        lngSm,
                    "ts":         ts,
                }
                if payload, err := json.Marshal(msg); err == nil {
                    _ = config.RDB.Publish(ctx, "order:location:"+orderID, string(payload)).Err()
                }
            }
        }
    }
	//每上来一个 GPS 点，先纠偏→平滑→存最新点→追加轨迹→自动清旧点，全程 Redis 内存操作，毫秒级完成，为后续派单、轨迹回放、实时位置推送提供干净数据
	return &v1.ReportCourierLocationReply{}, nil
}

// -------- 轨迹查询 --------
func (s *GreeterService) GetCourierTrajectory(ctx context.Context, req *v1.GetCourierTrajectoryRequest) (*v1.GetCourierTrajectoryReply, error) {
	courierID := courierIDFromCtx(ctx, req.CourierId)
	if courierID == 0 {
		return nil, status.Error(codes.InvalidArgument, "缺少快递员ID")
	}
	lastSec := req.LastSeconds
	if lastSec <= 0 {
		lastSec = 300
	}
	now := timeNow()
	from := now - int64(lastSec)

	points := make([]*v1.Point, 0)
	if config.RDB != nil {
		members, err := config.RDB.ZRangeByScore(ctx, trajKey(courierID), &redis.ZRangeBy{Min: fmt.Sprintf("%f", float64(from)), Max: "+inf"}).Result()
		if err == nil {
			// 稀释：最多返回 300 个点，按固定步长抽样
			step := 1
			if n := len(members); n > 300 {
				step = int(math.Ceil(float64(n) / 300.0))
			}
			for i := 0; i < len(members); i += step {
				var lat, lng float64
				var ts int64
				fmt.Sscanf(members[i], "%f,%f,%d", &lat, &lng, &ts)
				points = append(points, &v1.Point{Lat: lat, Lng: lng})
			}
		}
	}
	return &v1.GetCourierTrajectoryReply{Points: points}, nil
}

// 简单封装当前秒级时间戳
func timeNow() int64 { return time.Now().Unix() }

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// 高德地图API配置
const (
	AmapAPIKey = "a4de4b70bc1b6eee04c555cc460906ad"
	AmapGeoURL = "https://restapi.amap.com/v3/geocode/geo"
)

// 高德地图API响应结构
type AmapGeoResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Geocodes []struct {
		FormattedAddress string `json:"formatted_address"`
		Province         string `json:"province"`
		City             string `json:"city"`
		District         string `json:"district"`
		Location         string `json:"location"`
	} `json:"geocodes"`
}

// 地址到经纬度的映射表，用于模拟地址解析
var addressLocationMap = map[string]*v1.LocationInfo{
	// 北京地区
	"北京市朝阳区三里屯": {
		Lng: 116.447, Lat: 39.936,
	},
	"北京市海淀区中关村": {
		Lng: 116.298, Lat: 39.959,
	},
	"北京市东城区王府井": {
		Lng: 116.418, Lat: 39.914,
	},

	// 上海地区
	"上海市浦东新区陆家嘴": {
		Lng: 121.505, Lat: 31.245,
	},
	"上海市黄浦区外滩": {
		Lng: 121.490, Lat: 31.240,
	},
	"上海市徐汇区徐家汇": {
		Lng: 121.436, Lat: 31.188,
	},

	// 广州地区
	"广州市天河区珠江新城": {
		Lng: 113.322, Lat: 23.120,
	},
	"广州市越秀区北京路": {
		Lng: 113.270, Lat: 23.129,
	},

	// 深圳地区
	"深圳市南山区科技园": {
		Lng: 113.953, Lat: 22.537,
	},
	"深圳市福田区华强北": {
		Lng: 114.088, Lat: 22.548,
	},

	// 杭州地区
	"杭州市西湖区西湖": {
		Lng: 120.139, Lat: 30.259,
	},
	"杭州市滨江区网商路": {
		Lng: 120.210, Lat: 30.210,
	},

	// 青岛地区 - 新增
	"青岛市市南区香港中路": {
		Lng: 120.382, Lat: 36.067,
	},
	"青岛市市北区台东": {
		Lng: 120.374, Lat: 36.087,
	},
	"青岛市崂山区": {
		Lng: 120.469, Lat: 36.107,
	},

	// 天津地区
	"天津市和平区": {
		Lng: 117.195, Lat: 39.118,
	},
	"天津市河西区": {
		Lng: 117.223, Lat: 39.101,
	},

	// 重庆地区
	"重庆市渝中区解放碑": {
		Lng: 106.569, Lat: 29.559,
	},
	"重庆市江北区观音桥": {
		Lng: 106.532, Lat: 29.606,
	},

	// 成都地区
	"成都市锦江区春熙路": {
		Lng: 104.081, Lat: 30.660,
	},
	"成都市武侯区": {
		Lng: 104.043, Lat: 30.641,
	},

	// 武汉地区
	"武汉市武昌区": {
		Lng: 114.315, Lat: 30.554,
	},
	"武汉市汉口区": {
		Lng: 114.273, Lat: 30.593,
	},

	// 西安地区
	"西安市雁塔区": {
		Lng: 108.928, Lat: 34.222,
	},
	"西安市碑林区": {
		Lng: 108.953, Lat: 34.230,
	},

	// 南京地区
	"南京市鼓楼区": {
		Lng: 118.763, Lat: 32.066,
	},
	"南京市玄武区": {
		Lng: 118.797, Lat: 32.048,
	},

	// 苏州地区
	"苏州市姑苏区": {
		Lng: 120.619, Lat: 31.299,
	},
	"苏州市工业园区": {
		Lng: 120.728, Lat: 31.289,
	},

	// 大连地区
	"大连市中山区": {
		Lng: 121.618, Lat: 38.920,
	},
	"大连市沙河口区": {
		Lng: 121.594, Lat: 38.904,
	},

	// 宁波地区
	"宁波市海曙区": {
		Lng: 121.549, Lat: 29.868,
	},
	"宁波市江北区": {
		Lng: 121.550, Lat: 29.888,
	},

	// 厦门地区
	"厦门市思明区": {
		Lng: 118.082, Lat: 24.445,
	},
	"厦门市湖里区": {
		Lng: 118.113, Lat: 24.512,
	},

	// 福州地区
	"福州市鼓楼区": {
		Lng: 119.302, Lat: 26.081,
	},
	"福州市台江区": {
		Lng: 119.313, Lat: 26.061,
	},

	// 济南地区
	"济南市历下区": {
		Lng: 117.073, Lat: 36.674,
	},
	"济南市市中区": {
		Lng: 117.021, Lat: 36.651,
	},

	// 郑州地区
	"郑州市金水区": {
		Lng: 113.660, Lat: 34.757,
	},
	"郑州市中原区": {
		Lng: 113.613, Lat: 34.748,
	},

	// 长沙地区
	"长沙市芙蓉区": {
		Lng: 112.988, Lat: 28.195,
	},
	"长沙市岳麓区": {
		Lng: 112.931, Lat: 28.235,
	},
}

// getSimulatedLocation 根据地址获取模拟的位置信息
func (s *GreeterService) getSimulatedLocation(address string) *v1.LocationInfo {
	log.Printf("正在处理地址: %s (长度: %d, 字节: %v)", address, len(address), []byte(address))

	// 首先尝试精确匹配
	if location, exists := addressLocationMap[address]; exists {
		log.Printf("精确匹配成功: %s", address)
		return &v1.LocationInfo{
			Lng: location.Lng,
			Lat: location.Lat,
		}
	}

	// 如果没有精确匹配，尝试模糊匹配
	addressCity := extractCityFromAddress(address)
	log.Printf("从地址 '%s' 提取的城市: '%s'", address, addressCity)

	if addressCity != "" {
		// 根据城市找到第一个匹配的位置
		for key, location := range addressLocationMap {
			keyCity := extractCityFromKey(key)
			log.Printf("检查映射键 '%s' 的城市: '%s'", key, keyCity)
			if keyCity == addressCity {
				log.Printf("模糊匹配成功: 地址城市 '%s' 匹配映射城市 '%s'", addressCity, keyCity)
				// 创建一个新的位置信息，只返回经纬度
				return &v1.LocationInfo{
					Lng: location.Lng + (rand.Float64()-0.5)*0.01, // 添加小幅随机偏移
					Lat: location.Lat + (rand.Float64()-0.5)*0.01,
				}
			}
		}
	}

	// 特殊处理：如果地址包含特定字符数量，可能是编码问题导致的中文地址
	// 根据地址长度和内容推测可能的城市
	var guessedLocation *v1.LocationInfo
	addressBytes := []byte(address)
	log.Printf("地址字节分析: %v", addressBytes)

	// 改进的降级逻辑：基于地址内容智能推测位置
	guessedLocation = s.intelligentLocationGuess(address)

	if guessedLocation != nil {
		log.Printf("智能推测位置成功: %s", address)
		return guessedLocation
	}

	log.Printf("未找到匹配，使用默认位置: %s", address)
	// 如果都没有匹配，返回默认的北京位置
	return &v1.LocationInfo{
		Lng: 116.397428 + (rand.Float64()-0.5)*0.01,
		Lat: 39.90923 + (rand.Float64()-0.5)*0.01,
	}
}

// intelligentLocationGuess 智能推测地址位置
func (s *GreeterService) intelligentLocationGuess(address string) *v1.LocationInfo {
	addressBytes := []byte(address)

	// 定义主要城市的市中心坐标
	cityCenters := map[string]*v1.LocationInfo{
		"北京": {Lng: 116.397428, Lat: 39.90923},  // 天安门
		"上海": {Lng: 121.473701, Lat: 31.230416}, // 人民广场
		"广州": {Lng: 113.280637, Lat: 23.125178}, // 天河区
		"深圳": {Lng: 114.085947, Lat: 22.547},    // 福田区
		"杭州": {Lng: 120.153576, Lat: 30.287459}, // 西湖区
		"南京": {Lng: 118.767413, Lat: 32.041544}, // 新街口
		"苏州": {Lng: 120.619585, Lat: 31.299379}, // 姑苏区
		"青岛": {Lng: 120.382639, Lat: 36.067082}, // 市南区
		"天津": {Lng: 117.190182, Lat: 39.125596}, // 和平区
		"重庆": {Lng: 106.504962, Lat: 29.533155}, // 渝中区
		"成都": {Lng: 104.065735, Lat: 30.659462}, // 锦江区
		"武汉": {Lng: 114.298572, Lat: 30.584355}, // 江汉区
		"西安": {Lng: 108.948024, Lat: 34.263161}, // 新城区
		"大连": {Lng: 121.618622, Lat: 38.91459},  // 中山区
		"宁波": {Lng: 121.549792, Lat: 29.868388}, // 海曙区
		"厦门": {Lng: 118.11022, Lat: 24.490474},  // 思明区
		"福州": {Lng: 119.306239, Lat: 26.075302}, // 鼓楼区
		"济南": {Lng: 117.000923, Lat: 36.675807}, // 历下区
		"郑州": {Lng: 113.665412, Lat: 34.757975}, // 金水区
		"长沙": {Lng: 112.982279, Lat: 28.19409},  // 芙蓉区
	}

	// 1. 优先检查地址中是否包含具体城市名
	for cityName, location := range cityCenters {
		if strings.Contains(address, cityName) {
			log.Printf("在地址中找到城市: %s，返回市中心坐标", cityName)
			return &v1.LocationInfo{
				Lng: location.Lng,
				Lat: location.Lat,
			}
		}
	}

	// 2. 检查是否包含中文字符的UTF-8编码模式
	if len(address) >= 4 {
		if containsChinesePattern(addressBytes, []byte("青岛")) || containsChinesePattern(addressBytes, []byte("市南")) || containsChinesePattern(addressBytes, []byte("香港中路")) {
			log.Printf("通过字节模式识别为青岛")
			return cityCenters["青岛"]
		} else if containsChinesePattern(addressBytes, []byte("广州")) || containsChinesePattern(addressBytes, []byte("天河")) {
			log.Printf("通过字节模式识别为广州")
			return cityCenters["广州"]
		} else if containsChinesePattern(addressBytes, []byte("深圳")) || containsChinesePattern(addressBytes, []byte("南山")) {
			log.Printf("通过字节模式识别为深圳")
			return cityCenters["深圳"]
		} else if containsChinesePattern(addressBytes, []byte("北京")) || containsChinesePattern(addressBytes, []byte("朝阳")) {
			log.Printf("通过字节模式识别为北京")
			return cityCenters["北京"]
		} else if containsChinesePattern(addressBytes, []byte("上海")) || containsChinesePattern(addressBytes, []byte("浦东")) {
			log.Printf("通过字节模式识别为上海")
			return cityCenters["上海"]
		}
	}

	// 3. 基于省份信息推测主要城市
	if strings.Contains(address, "山东") {
		log.Printf("识别为山东省，返回青岛市中心")
		return cityCenters["青岛"]
	} else if strings.Contains(address, "广东") {
		log.Printf("识别为广东省，返回广州市中心")
		return cityCenters["广州"]
	} else if strings.Contains(address, "江苏") {
		log.Printf("识别为江苏省，返回南京市中心")
		return cityCenters["南京"]
	} else if strings.Contains(address, "浙江") {
		log.Printf("识别为浙江省，返回杭州市中心")
		return cityCenters["杭州"]
	} else if strings.Contains(address, "四川") {
		log.Printf("识别为四川省，返回成都市中心")
		return cityCenters["成都"]
	} else if strings.Contains(address, "湖北") {
		log.Printf("识别为湖北省，返回武汉市中心")
		return cityCenters["武汉"]
	} else if strings.Contains(address, "陕西") {
		log.Printf("识别为陕西省，返回西安市中心")
		return cityCenters["西安"]
	}

	// 4. 基于地址长度的一致性推测（移除随机性）
	if len(address) >= 10 {
		// 使用地址哈希确保相同地址返回相同结果
		hash := 0
		for _, b := range []byte(address) {
			hash = hash*31 + int(b)
		}

		cities := []string{"北京", "上海", "广州", "深圳"}
		cityIndex := hash % len(cities)
		if cityIndex < 0 {
			cityIndex = -cityIndex
		}

		selectedCity := cities[cityIndex]
		log.Printf("基于地址哈希选择城市: %s", selectedCity)
		return cityCenters[selectedCity]
	}

	return nil
}

// containsChinesePattern 检查字节数组是否包含中文字符模式
func containsChinesePattern(haystack, needle []byte) bool {
	if len(needle) == 0 {
		return false
	}

	// 简单的字节匹配
	for i := 0; i <= len(haystack)-len(needle); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// fixEncodingIssue 尝试修复编码问题
func (s *GreeterService) fixEncodingIssue(address string) string {
	// 如果地址全是问号，可能是编码问题
	if strings.Trim(address, "?") == "" && len(address) > 0 {
		log.Printf("检测到全问号地址，可能是编码问题: %s", address)
		return address // 暂时返回原地址，后续可以添加更复杂的修复逻辑
	}

	// 检查是否包含问号字符
	if strings.Contains(address, "?") {
		log.Printf("地址包含问号字符，可能存在编码问题: %s", address)
	}

	return address
}

// extractCityFromKey 从映射表的键中提取城市名
func extractCityFromKey(key string) string {
	if strings.Contains(key, "北京市") {
		return "北京"
	} else if strings.Contains(key, "上海市") {
		return "上海"
	} else if strings.Contains(key, "广州市") {
		return "广州"
	} else if strings.Contains(key, "深圳市") {
		return "深圳"
	} else if strings.Contains(key, "杭州市") {
		return "杭州"
	}
	return ""
}

// extractCityFromAddress 从地址中提取城市名
func extractCityFromAddress(address string) string {
	// 扩展城市列表，包含更多主要城市
	cities := []string{
		"北京", "上海", "广州", "深圳", "杭州", "天津", "重庆", "成都", "武汉", "西安",
		"南京", "苏州", "青岛", "大连", "宁波", "厦门", "福州", "济南", "郑州", "长沙",
		"合肥", "南昌", "太原", "石家庄", "哈尔滨", "长春", "沈阳", "呼和浩特",
		"银川", "西宁", "兰州", "乌鲁木齐", "拉萨", "昆明", "贵阳", "南宁", "海口",
		"三亚", "温州", "嘉兴", "绍兴", "金华", "台州", "湖州", "丽水", "衢州",
		"舟山", "无锡", "常州", "徐州", "扬州", "泰州", "镇江", "盐城", "淮安",
		"连云港", "宿迁", "佛山", "东莞", "中山", "珠海", "汕头", "江门", "湛江",
		"茂名", "肇庆", "惠州", "梅州", "汕尾", "河源", "阳江", "清远", "潮州",
		"揭阳", "云浮", "韶关", "唐山", "秦皇岛", "邯郸", "邢台", "保定", "张家口",
		"承德", "沧州", "廊坊", "衡水", "包头", "鞍山", "抚顺", "本溪", "丹东",
		"锦州", "营口", "阜新", "辽阳", "盘锦", "铁岭", "朝阳", "葫芦岛", "吉林",
		"四平", "辽源", "通化", "白山", "松原", "白城", "延边", "齐齐哈尔", "鸡西",
		"鹤岗", "双鸭山", "大庆", "伊春", "佳木斯", "七台河", "牡丹江", "黑河",
		"绥化", "大兴安岭", "洛阳", "开封", "安阳", "鹤壁", "新乡", "焦作", "濮阳",
		"许昌", "漯河", "三门峡", "南阳", "商丘", "信阳", "周口", "驻马店", "荆门",
		"鄂州", "荆州", "黄冈", "孝感", "黄石", "咸宁", "随州", "恩施", "仙桃",
		"潜江", "天门", "神农架", "株洲", "湘潭", "衡阳", "邵阳", "岳阳", "常德",
		"张家界", "益阳", "郴州", "永州", "怀化", "娄底", "湘西",
	}

	// 按长度排序，优先匹配较长的城市名（避免"长春"被"春"匹配）
	sort.Slice(cities, func(i, j int) bool {
		return len(cities[i]) > len(cities[j])
	})

	for _, city := range cities {
		if strings.Contains(address, city) {
			log.Printf("在地址 '%s' 中找到城市: '%s'", address, city)
			return city
		}
	}

	log.Printf("未在地址 '%s' 中找到匹配的城市", address)
	return ""
}

// GeocodeAddress 地址解析接口，将寄件地址和收件地址转换为经纬度
func (s *GreeterService) GeocodeAddress(ctx context.Context, req *v1.GeocodeAddressRequest) (*v1.GeocodeAddressResponse, error) {
	resp := &v1.GeocodeAddressResponse{}

	log.Printf("收到地址解析请求 - 寄件地址: '%s' (UTF-8有效: %v, 字节: %v), 收件地址: '%s' (UTF-8有效: %v, 字节: %v)",
		req.SenderAddress, utf8.ValidString(req.SenderAddress), []byte(req.SenderAddress),
		req.ReceiverAddress, utf8.ValidString(req.ReceiverAddress), []byte(req.ReceiverAddress))

	// 尝试检测和修复编码问题
	senderAddr := s.fixEncodingIssue(req.SenderAddress)
	receiverAddr := s.fixEncodingIssue(req.ReceiverAddress)

	log.Printf("修复编码后 - 寄件地址: '%s', 收件地址: '%s'", senderAddr, receiverAddr)

	// 解析寄件地址
	if senderAddr != "" {
		log.Printf("开始解析寄件地址: %s", senderAddr)
		senderLocation, err := s.geocodeAddress(senderAddr)
		if err != nil {
			log.Printf("API调用失败，使用模拟数据: %v", err)
			// 如果API调用失败，返回基于地址的模拟数据
			resp.SenderLocation = s.getSimulatedLocation(senderAddr)
		} else {
			log.Printf("API调用成功，使用真实数据")
			resp.SenderLocation = senderLocation
		}
	}

	// 解析收件地址
	if receiverAddr != "" {
		log.Printf("开始解析收件地址: %s", receiverAddr)
		receiverLocation, err := s.geocodeAddress(receiverAddr)
		if err != nil {
			log.Printf("API调用失败，使用模拟数据: %v", err)
			// 如果API调用失败，返回基于地址的模拟数据
			resp.ReceiverLocation = s.getSimulatedLocation(receiverAddr)
		} else {
			log.Printf("API调用成功，使用真实数据")
			resp.ReceiverLocation = receiverLocation
		}
	}

	return resp, nil
}

// geocodeAddress 调用高德地图API获取地址的经纬度信息
func (s *GreeterService) geocodeAddress(address string) (*v1.LocationInfo, error) {
	// 构建请求URL
	params := url.Values{}
	params.Set("key", AmapAPIKey)
	params.Set("address", address)

	requestURL := fmt.Sprintf("%s?%s", AmapGeoURL, params.Encode())

	// 发送HTTP请求
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("请求高德地图API失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var amapResp AmapGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&amapResp); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v", err)
	}

	// 检查API响应状态
	if amapResp.Status != "1" {
		return nil, fmt.Errorf("高德地图API返回错误: %s", amapResp.Info)
	}

	// 检查是否有地理编码结果
	if len(amapResp.Geocodes) == 0 {
		return nil, fmt.Errorf("未找到地址对应的地理位置")
	}

	geocode := amapResp.Geocodes[0]

	// 解析经纬度
	var lng, lat float64
	if _, err := fmt.Sscanf(geocode.Location, "%f,%f", &lng, &lat); err != nil {
		return nil, fmt.Errorf("解析经纬度失败: %v", err)
	}

	return &v1.LocationInfo{
		Lng: lng,
		Lat: lat,
	}, nil
}

// SayHello 方法已被移除，因为 HelloRequest 和 HelloReply 类型不存在
