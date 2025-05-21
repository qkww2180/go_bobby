package main
/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/
import (
	"OldPackageTest/grpclb_test/proto"
	"context"
	"fmt"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(
		"consul://192.168.1.103:8500/user-srv?wait=14s&tag=srv",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		userSrvClient := proto.NewUserClient(conn)
		rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
			Pn:    1,
			PSize: 2,
		})
		if err != nil {
			panic(err)
		}
		for index, data := range rsp.Data {
			fmt.Println(index, data)
		}
	}

}
