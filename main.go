package main

import (
	"context"
	"flag"
	"my_test/cache"
	"my_test/db"
	"my_test/kafka"
	"my_test/server"
)

func main() {

	http_port := flag.String("http_port", "8081", "http port")
	external := flag.Bool("external", false, "use external kafka host:port")
	group_id := flag.String("group_id", "0", "kafka consumer group id")
	flag.Parse()

	db.InitConfig(*external)
	db.InitLog()
	cache.InitGCache()
	kafka.InitKafkaCfg()
	go server.StartConsumerServer(context.Background(), *group_id)
	server.StartHttpServer(*http_port)
}
