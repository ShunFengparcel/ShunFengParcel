# ELK 栈部署指南

## 概述

ELK 栈（Elasticsearch + Logstash + Kibana）用于日志收集、存储、分析和可视化。

## 架构图

```
应用服务 (ShunFengParcel)
    ↓ (JSON 日志文件)
Filebeat (日志采集)
    ↓
Logstash (日志处理)
    ↓
Elasticsearch (日志存储)
    ↓
Kibana (日志可视化)
```

## 1. Elasticsearch 部署

### Docker 部署

```bash
# 创建网络
docker network create elk

# 启动 Elasticsearch
docker run -d \
  --name elasticsearch \
  --net elk \
  -p 9200:9200 \
  -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  -v es-data:/usr/share/elasticsearch/data \
  elasticsearch:8.11.0
```

### 验证安装

```bash
curl http://localhost:9200
```

## 2. Kibana 部署

```bash
docker run -d \
  --name kibana \
  --net elk \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" \
  kibana:8.11.0
```

访问：http://localhost:5601

## 3. Logstash 部署

### 创建配置文件

```bash
mkdir -p /etc/logstash/conf.d
```

创建 `logstash.conf`:

```conf
input {
  file {
    path => "/var/log/shunfeng/*.log"
    start_position => "beginning"
    codec => json
    type => "shunfeng-app"
  }
}

filter {
  # 解析 JSON 日志
  if [type] == "shunfeng-app" {
    json {
      source => "message"
    }
    
    # 添加地理位置信息（如果有 IP）
    if [client_ip] {
      geoip {
        source => "client_ip"
      }
    }
    
    # 时间戳处理
    date {
      match => ["timestamp", "ISO8601"]
      target => "@timestamp"
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "shunfeng-logs-%{+YYYY.MM.dd}"
  }
  
  # 调试输出（可选）
  stdout {
    codec => rubydebug
  }
}
```

### 启动 Logstash

```bash
docker run -d \
  --name logstash \
  --net elk \
  -p 5044:5044 \
  -v /etc/logstash/conf.d:/usr/share/logstash/pipeline \
  -v /path/to/logs:/var/log/shunfeng \
  logstash:8.11.0
```

## 4. Filebeat 部署（推荐）

Filebeat 比 Logstash 更轻量，适合日志采集。

### 安装 Filebeat

```bash
# Windows
# 下载：https://www.elastic.co/downloads/beats/filebeat
# 解压到 C:\Program Files\Filebeat

# Linux
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-8.11.0-linux-x86_64.tar.gz
tar xzvf filebeat-8.11.0-linux-x86_64.tar.gz
cd filebeat-8.11.0-linux-x86_64/
```

### 配置 filebeat.yml

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - E:/gowork/src/ShunFengParcel/ShunFengParcel/logs/*.log
  json.keys_under_root: true
  json.add_error_key: true
  fields:
    service: shunfeng-parcel
    env: production

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "shunfeng-logs-%{+yyyy.MM.dd}"

setup.kibana:
  host: "localhost:5601"

# 日志级别
logging.level: info
```

### 启动 Filebeat

```bash
# Windows
.\filebeat.exe -e -c filebeat.yml

# Linux
./filebeat -e -c filebeat.yml
```

## 5. Docker Compose 一键部署

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  elasticsearch:
    image: elasticsearch:8.11.0
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms512m -Xmx512m
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - es-data:/usr/share/elasticsearch/data
    networks:
      - elk

  kibana:
    image: kibana:8.11.0
    container_name: kibana
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    networks:
      - elk

  logstash:
    image: logstash:8.11.0
    container_name: logstash
    volumes:
      - ./logstash/pipeline:/usr/share/logstash/pipeline
      - ./logs:/var/log/shunfeng
    ports:
      - "5044:5044"
    depends_on:
      - elasticsearch
    networks:
      - elk

volumes:
  es-data:
    driver: local

networks:
  elk:
    driver: bridge
```

启动：

```bash
docker-compose up -d
```

## 6. Kibana 配置

### 创建索引模式

1. 访问 Kibana: http://localhost:5601
2. 进入 **Management** → **Stack Management** → **Index Patterns**
3. 点击 **Create index pattern**
4. 输入：`shunfeng-logs-*`
5. 选择时间字段：`@timestamp`
6. 点击 **Create**

### 查看日志

1. 进入 **Discover**
2. 选择索引模式：`shunfeng-logs-*`
3. 查看实时日志

## 7. 日志解析规则

### Grok 模式（如果使用文本格式）

```conf
filter {
  grok {
    match => {
      "message" => "\[%{TIMESTAMP_ISO8601:timestamp}\] %{LOGLEVEL:level} %{GREEDYDATA:message}"
    }
  }
}
```

### JSON 解析（推荐）

```conf
filter {
  json {
    source => "message"
    target => "parsed"
  }
  
  mutate {
    rename => {
      "[parsed][timestamp]" => "@timestamp"
      "[parsed][level]" => "log_level"
      "[parsed][message]" => "log_message"
    }
  }
}
```

## 8. 智能告警规则

### 在 Kibana 中创建告警

1. 进入 **Stack Management** → **Rules and Connectors**
2. 点击 **Create rule**

### 告警规则示例

#### 错误率告警

```
条件：
- 索引：shunfeng-logs-*
- 查询：level: "ERROR"
- 时间窗口：5 分钟
- 阈值：count > 100

动作：
- 发送邮件
- 发送钉钉/企业微信通知
```

#### 慢请求告警

```
条件：
- 索引：shunfeng-logs-*
- 查询：duration_ms > 3000
- 时间窗口：5 分钟
- 阈值：count > 10

动作：
- 发送告警通知
```

#### 服务异常告警

```
条件：
- 索引：shunfeng-logs-*
- 查询：level: "FATAL"
- 时间窗口：1 分钟
- 阈值：count > 0

动作：
- 立即发送紧急告警
- 触发 PagerDuty
```

## 9. 实时索引构建

### 索引模板

```json
PUT _index_template/shunfeng-logs-template
{
  "index_patterns": ["shunfeng-logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 1,
      "index.refresh_interval": "5s"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "level": { "type": "keyword" },
        "message": { "type": "text" },
        "service_id": { "type": "keyword" },
        "trace_id": { "type": "keyword" },
        "duration_ms": { "type": "long" },
        "error": { "type": "text" }
      }
    }
  }
}
```

### 索引生命周期管理（ILM）

```json
PUT _ilm/policy/shunfeng-logs-policy
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_size": "50GB",
            "max_age": "7d"
          }
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

## 10. 性能优化

### Elasticsearch 优化

```yaml
# elasticsearch.yml
indices.memory.index_buffer_size: 30%
thread_pool.write.queue_size: 1000
```

### Logstash 优化

```yaml
# logstash.yml
pipeline.workers: 4
pipeline.batch.size: 125
pipeline.batch.delay: 50
```

## 11. 监控和维护

### 检查集群健康

```bash
curl http://localhost:9200/_cluster/health?pretty
```

### 查看索引状态

```bash
curl http://localhost:9200/_cat/indices?v
```

### 清理旧索引

```bash
# 删除 30 天前的索引
curator_cli --host localhost delete_indices --filter_list \
  '[{"filtertype":"age","source":"name","direction":"older","timestring":"%Y.%m.%d","unit":"days","unit_count":30}]'
```

## 12. 安全配置

### 启用认证

```yaml
# elasticsearch.yml
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
```

### 创建用户

```bash
bin/elasticsearch-users useradd admin -p password -r superuser
```

## 故障排查

### 常见问题

1. **Elasticsearch 无法启动**
   - 检查内存设置
   - 检查磁盘空间
   - 查看日志：`docker logs elasticsearch`

2. **Filebeat 无法连接**
   - 检查网络连接
   - 验证 Elasticsearch 地址
   - 查看 Filebeat 日志

3. **日志未显示**
   - 检查索引模式
   - 验证时间范围
   - 检查 Filebeat 配置

## 总结

完成以上步骤后，你将拥有：
- ✅ 完整的日志收集系统
- ✅ 实时日志分析能力
- ✅ 智能告警机制
- ✅ 可视化日志查询界面

下一步可以根据业务需求创建自定义仪表板和更复杂的告警规则。
