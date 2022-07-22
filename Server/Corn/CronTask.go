package Corn

import (
	"TencentVideoCheck/Server/Config"
	"database/sql"
	"fmt"
	"github.com/robfig/cron"
	"log"
)

func CronTask() {

	//*/5 * * * * ? 每5秒执行一次

	//每天23：30执行一次
	//0 30 23 * * ?

	//每天中午12：00执行一次
	//0 0 12 * * ?

	refreshCron := cron.New()
	checkCron := cron.New()

	//检查Cookie是否失效
	err := refreshCron.AddFunc("0 0 12 * * ?", func() {
		dsn := Config.GetDsn()
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Printf("连接数据库出错：%v\n", err)
			return
		}

		rows, err := db.Query("SELECT * FROM user")
		if err != nil {
			log.Fatalf("查询数据库出错：", err)
			return
		}

		for rows.Next() {
			var cookie string
			var barrage int
			var check int
			var download int
			var giving int
			var obtained string
			var userInfo string
			var notice string
			var noticeToken string
			err = rows.Scan(&cookie, &barrage, &check, &download, &giving, &obtained, &userInfo, &notice, &noticeToken)
			if err != nil {
				log.Fatalf("遍历数据库出错：", err)
				return
			}
			//查询打印结果集
			//fmt.Println(cookie)
			if Refresh(cookie) {
				log.Printf("Cookie未失效")
			} else {
				log.Printf("Cookie已失效,正在删除Cookie")
				deleteDB, err := db.Prepare("DELETE FROM `user` WHERE `Cookie`=?")
				if err != nil {
					fmt.Println(err)
				}
				_, err = deleteDB.Exec(cookie)
				if err != nil {
					fmt.Printf("删除数据出错：%v\n", err)
				}
			}
		}
	})
	if err != nil {
		log.Printf("检查Cookie定时任务出错：%v\n", err)
		return
	}

	//自动领取已完成任务产生的V力值
	checkCron.AddFunc("*/5 * * * * ?", func() {
		dsn := Config.GetDsn()
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Printf("连接数据库出错：%v\n", err)
			return
		}

		rows, err := db.Query("SELECT * FROM user")
		if err != nil {
			log.Fatalf("查询数据库出错：", err)
			return
		}

		for rows.Next() {
			var cookie string
			var barrage int
			var check int
			var download int
			var giving int
			var watch int
			var obtained string
			var userInfo string
			var notice string
			var noticeToken string
			err = rows.Scan(&cookie, &barrage, &check, &download, &giving, &watch, &obtained, &userInfo, &notice, &noticeToken)
			if err != nil {
				log.Fatalf("遍历数据库出错：", err)
				return
			}
			//查询打印结果集
			//fmt.Println(cookie)

			//弹幕签到
			BarrageCheck(cookie)
			//签到
			Check(cookie)
			//下载签到
			DownloadCheck(cookie)
			//赠送签到
			GivingCheck(cookie)
			//观看60分钟签到
			WatchCheck(cookie)
			//本月获得V力值
			Obtained(cookie)
			//用户信息
			UserInfo(cookie)
		}
	})

	refreshCron.Start()
	checkCron.Start()
	select {}
}