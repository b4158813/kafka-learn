# kafka-learn

一些基本的kafka服务端和客户端操作

## 制作consumer/producer服务的镜像
```bash
docker build -f Dockerfile -t kafka-learn .
```

## 部署kafka集群 + 部署基于此kafka的本地服务
```bash
docker-compose -f docker-compose.yml up -d
```

## 一些疑难解答

- 本go项目主要是学习kafka的原理和部署，主要内容：
  -  开启http服务，通过向`/users/{id}`发送Post请求的方式将id写入到kafka的producer中
  - 写入后，consumer会消费，直接打印对应消息

- `docker-compose.yml`文件中定义了两套服务，一个是kafka-server，另一个是本go项目的服务，两者通过docker network定义的app-tier内网进行通信
- 两个容器在通过docker network通信时，kafka客户端连接地址应该使用对内的localhost:9092的形式。如果需要通过宿主机网络中的服务访问，则客户端连接地址应该使用kafka:9094（即 docker-compose.yml中定义的服务名 + EXTERNAL的端口号）
- 本机的go服务启动时，可以指定以下参数：
  ```shell
  -http_port 8080 # http服务端口号
  -external true # 是否外部网络访问kafka服务
  -group_id 1 # 消费者组id
  ```

- 同一个消费者组id的消费者们就像是一个分布式消费者集群，共同互斥地消费同一个topic下的所有partition的消息，kafka会通过某种负载均衡策略（range或者round-robin）指定消费者消费哪个分区的数据（如果改消费者组的消费者数量大于了该topic下的partition数量，则根据特定的负载策略会有消费者永远无法消费数据）