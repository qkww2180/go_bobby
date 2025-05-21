package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

type UserAccount struct {
	ID             int     `gorm:"column:id;primary_key"`
	UserId         int     `gorm:"user_id"`
	Balance        float64 `gorm:"balance"`
	TradingBalance float64 `gorm:"trading_balance"`
}

func (UserAccount) TableName() string {
	return "user_account"
}

var lock sync.Mutex

// 转入和转出的时候，都要加锁，否则会出现并发问题
func SagaAdjustBalance(db *sql.Tx, uid int, amount float64) error {
	lock.Lock()
	defer lock.Unlock()

	if amount < 0 {
		var balance float64
		db.QueryRow("select balance from dtm.user_account where user_id = ?", uid).Scan(&balance)
		if balance < -amount {
			return fmt.Errorf("余额不足")
		}
	}
	_, err := db.Exec("update dtm.user_account set balance = balance + ? where user_id = ?", amount, uid)
	if err != nil {
		return err
	}
	return nil
}

var db *gorm.DB

func initDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root",
		"127.0.0.1",
		"3306",
		"dtm")
	newLogger := glog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		glog.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  glog.Info,   // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	return nil
}

// MustBarrierFromGin 1
func MustBarrierFromGin(c *gin.Context) *dtmcli.BranchBarrier {
	ti, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	fmt.Println(err)
	return ti
}

// 服务发现， 库存服务有5个
func main() {
	err := initDB()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.POST("/SagaBTransIn", func(c *gin.Context) {
		barrier := MustBarrierFromGin(c)
		tx := db.Begin()
		sourceTx := tx.Statement.ConnPool.(*sql.Tx)
		err := barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
			fmt.Println("开始转入")
			userID := 1
			err := SagaAdjustBalance(sourceTx, userID, 100)
			if err != nil {
				fmt.Printf("转入失败:%s\r\n", err.Error())
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
			return
		}

		return
	})

	r.POST("/SagaBTransInCom", func(c *gin.Context) {
		fmt.Println("转入失败， 开始补偿")
		//userID := 1
		//err := SagaAdjustBalance(db, userID, -100)
		//if err != nil {
		//	fmt.Printf("转入补偿失败:%s\r\n", err.Error())
		//	return
		//}
		fmt.Println("转入补偿成功")
	})

	r.POST("/SagaBTransOut", func(c *gin.Context) {
		barrier := MustBarrierFromGin(c)
		tx := db.Begin()
		sourceTx := tx.Statement.ConnPool.(*sql.Tx)

		err := barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
			fmt.Println("开始转出")
			userID := 3
			err := SagaAdjustBalance(sourceTx, userID, -100)
			if err != nil {
				if err.Error() == "余额不足" {
					c.JSON(http.StatusConflict, gin.H{})
				}
				fmt.Printf("转出失败:%s\r\n", err.Error())
				c.JSON(500, gin.H{"msg": err.Error()})
			}
			fmt.Println("转出成功")
			return nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		return
	})

	r.POST("/SagaBTransOutCom", func(c *gin.Context) {
		fmt.Println("转出失败， 开始补偿")
		//userID := 3
		//err := SagaAdjustBalance(db, userID, 100)
		//if err != nil {
		//	fmt.Printf("转出补偿失败:%s\r\n", err.Error())
		//	return
		//}
		fmt.Println("转出补偿成功")
	})

	r.GET("start", func(c *gin.Context) {
		req := gin.H{}
		dmtServer := "http://127.0.0.1:36789/api/dtmsvr"
		qsBusi := "http://127.0.0.1:8089"
		saga := dtmcli.NewSaga(dmtServer, shortuuid.New()).
			// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransOutCom"
			Add(qsBusi+"/SagaBTransOut", qsBusi+"/SagaBTransOutCom", req).
			// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransInCom"
			Add(qsBusi+"/SagaBTransIn", qsBusi+"/SagaBTransInCom", req)
		// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
		saga.WaitResult = true
		err := saga.Submit()
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Run(":8089")
}
