package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"math/rand"
	"sync/atomic"
	"time"
)

type Counter struct {
	pass  *int64
	block *int64
	total *int64
}

func main() {
	//流量计数器,为了流控打印日志更直观,和集成nacos数据源无关。
	counter := Counter{
		pass:  new(int64),
		block: new(int64),
		total: new(int64),
	}

	//nacos server地址
	sc := []constant.ServerConfig{
		{
			ContextPath: "/nacos",
			Port:        8848,
			IpAddr:      "39.107.30.137",
		},
	}

	//nacos client 相关参数配置,具体配置可参考github.com/nacos-group/nacos-sdk-go
	cc := constant.ClientConfig{
		NamespaceId: "public",
		TimeoutMs:   5000,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	//注册流控规则Handler
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	//创建NacosDataSource数据源
	nds, err := nacos.NewNacosDataSource(client, "sentinel-go", "flow", h)
	if err != nil {
		panic(err)
	}

	err = nds.Initialize()
	if err != nil {
		panic(err)
	}

	go timerTask(&counter)

	//模拟流量
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				atomic.AddInt64(counter.total, 1)
				a, b := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					atomic.AddInt64(counter.block, 1)
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					atomic.AddInt64(counter.pass, 1)
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					a.Exit()
				}
			}
		}()
	}
	<-ch
}

func timerTask(counter *Counter) {
	var (
		oldTotal, oldPass, oldBlock int64
	)

	for {
		time.Sleep(time.Second)
		globalTotal := atomic.LoadInt64(counter.total)

		oneSecondTotal := globalTotal - oldTotal
		oldTotal = globalTotal

		globalPass := atomic.LoadInt64(counter.pass)
		oneSecondPass := globalPass - oldPass
		oldPass = globalPass

		globalBlock := atomic.LoadInt64(counter.block)
		oneSecondBlock := globalBlock - oldBlock
		oldBlock = globalBlock

		fmt.Println("total:", oneSecondTotal, "pass:", oneSecondPass, "block:", oneSecondBlock)
	}
}
