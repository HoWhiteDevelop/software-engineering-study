package go_ws

import "github.com/robfig/cron/v3"

func CleanOfflineConn() {
	//// 创建一个新的定时任务调度器
	c := cron.New()

	// 每天定时执行的条件
	//每分钟执行一次
	spec := `* * * * *`
	// 添加定时任务
	c.AddFunc(spec, func() {
		// fmt.Println("CleanOfflineConn")
		HandelOfflineCoon()
	})

	go c.Start()
}
