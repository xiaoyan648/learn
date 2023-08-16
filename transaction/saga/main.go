package busi

import (
	"fmt"
	"log"
	"time"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

// busi address
const (
	qsBusiAPI  = "/api/busi_start"
	qsBusiPort = 8082
)

var qsBusi = fmt.Sprintf("http://localhost:%d%s", qsBusiPort, qsBusiAPI)
var (
	user1Balance = 1000
	user2Balance = 500
)

// QsMain will be call from dtm/qs
func QsMain() {
	QsStartSvr()
	QsFireRequest()
	select {}
}

// QsStartSvr quick start: start server
func QsStartSvr() {
	app := gin.New()
	qsAddRoute(app)
	log.Printf("quick start examples listening at %d", qsBusiPort)
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", qsBusiPort))
	}()
	time.Sleep(100 * time.Millisecond)
}

func qsAddRoute(app *gin.Engine) {
	app.POST(qsBusiAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn")
		user1Balance = 500
		log.Printf("user1 %d, user2 %d", user1Balance, user2Balance)
		// c.JSON(200, "")
		// c.JSON(500, "") // retry
		c.JSON(409, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(qsBusiAPI+"/TransInCompensate", func(c *gin.Context) {
		log.Printf("TransInCompensate")
		user1Balance = 1000
		log.Printf("user1 %d, user2 %d", user1Balance, user2Balance)
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOut", func(c *gin.Context) {
		log.Printf("TransOut")
		user2Balance = 1000
		log.Printf("user1 %d, user2 %d", user1Balance, user2Balance)
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOutCompensate", func(c *gin.Context) {
		log.Printf("TransOutCompensate")
		user2Balance = 500
		log.Printf("user1 %d, user2 %d", user1Balance, user2Balance)
		c.JSON(200, "")
	})
}

const dtmServer = "http://localhost:36789/api/dtmsvr"

// QsFireRequest quick start: fire request
func QsFireRequest() string {
	req := &gin.H{"amount": 30} // load of micro-service
	// DtmServer is the url of dtm
	saga := dtmcli.NewSaga(dtmServer, shortuuid.New()).
		// add a TransOut sub-transaction，forward operation with url: qsBusi+"/TransOut", reverse compensation operation with url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// add a TransIn sub-transaction, forward operation with url: qsBusi+"/TransIn", reverse compensation operation with url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// saga.WithRetryLimit(5)
	// submit the created saga transaction，dtm ensures all sub-transactions either complete or get revoked
	err := saga.Submit()
	if err != nil {
		panic(err)
	}
	log.Printf("user1 %d, user2 %d", user1Balance, user2Balance)
	return saga.Gid
}
