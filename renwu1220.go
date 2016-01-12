package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"runtime"
	"time"
)

//var NutrientSolutionConfigurationCh chan int
//var CoolingCh chan interface{}

/*
var Nutrient_supply_ch  := make (chan   int)

var Nutrient_solution_configuration_ch:= make (chan interface { } )
var Cooling_ch := make (chan interface { } )
var Warming _ch:= make (chan interface { } )
var Ventilation_ch := make (chan interface { } )
var Dehumidification_ch := make (chan interface { } )
var Shading_ch := make (chan interface { } )
var CO2_ch := make (chan interface { } )
*/
//var data :=[9]string  {"_","_","_","_","_","_","_","_","_"}

//var WaitCh chan interface{}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Task_management() {
	NutrientSupplyCh := make(chan interface{}, 1)
	NutrientSupplyCh1 := make(chan interface{}, 1)
	NutrientSupplyCh2 := make(chan interface{}, 1) //
	NutrientSolutionConfigurationCh := make(chan int, 1)
	CoolingCh := make(chan interface{}, 1)
	WarmingCh := make(chan interface{}, 1)
	VentilationCh := make(chan interface{}, 1)
	DehumidificationCh := make(chan interface{}, 1)
	ShadingCh := make(chan interface{}, 1)
	CO2Ch := make(chan interface{}, 1)
	//taskname := make(chan interface{}, 1)
	TkNameCh := make(chan string, 1) //查询结果给GO
	WaitCh := make(chan interface{}) //use for wait

	db, err := sql.Open("sqlite3", "./farmwin2000shchgl.db")
	checkErr(err)
	defer db.Close()
	err = db.Ping()
	checkErr(err)

	db.Exec("update 生产过程控制 set '任务状态'=  '待执行' ")
	db.Exec("update 设备状态 set 设备状态='禁用' ")
	time.Sleep(time.Second)

	//chaxun
	rows, err := db.Query(" select 任务名称 from 生产过程控制 where 任务状态 = '待执行' limit 8")
	checkErr(err)
	defer rows.Close()

	go func() {

		waittimer := 60
		for {
			<-NutrientSupplyCh1
			rows2, err := db.Query(" select  参数上限 from 环境模式1  where 参数名称 = '营养液供给间隔'  ")
			checkErr(err)
			for rows2.Next() {
				err = rows2.Scan(&waittimer)
				checkErr(err)
				defer rows2.Close()
				fmt.Printf("营养液供给间隔:%d \n", waittimer)
			}
			NutrientSupplyCh2 <- 1
			fmt.Printf("营养液供给间隔1:%d \n", waittimer)
			time.Sleep(time.Second * time.Duration(waittimer))

			fmt.Printf("营养液供给间隔2:%d \n", waittimer)
		}
	}()

	//db.Exec(" update 生产过程控制 set '任务状态'= '?' where rowid= ? ", " 待执行", 1)
	go func() {
		var channelnum int
		channelnum = 8
		sptimer := 10
		for {
			select {
			case <-NutrientSupplyCh2:

				rows1, err := db.Query(" select 参数上限 from 环境模式1  where 参数名称 = '管路' ")
				checkErr(err)
				for rows1.Next() {
					err = rows1.Scan(&channelnum)
					checkErr(err)
					defer rows1.Close()
					fmt.Printf("管道:%d \n", channelnum)
				}

				rows, err := db.Query(" select 参数上限 from 环境模式1  where 参数名称 = '雾化时长' ")
				checkErr(err)
				for rows.Next() {
					err = rows.Scan(&sptimer)
					checkErr(err)
					defer rows.Close()
					fmt.Printf("雾化时长:%d \n", sptimer)
				}
				//rows.Close()

				//	rows1.Close()

				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=1")
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=2")
				for i := 2; i < channelnum; i++ {
					db.Exec("  update 设备状态 set 设备状态=  '关闭'  where rowid=?", i)
					db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid= ? ", (i + 1))
					time.Sleep(time.Second * time.Duration(sptimer))
					fmt.Printf("管道%d 开启\n", i)
				}

				db.Exec("  update 设备状态 set 设备状态=  '关闭'  where rowid=1")
				db.Exec("  update 设备状态 set 设备状态=  '关闭'  where rowid= ? ", channelnum)
			//	db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=13")

			default:
				runtime.Gosched()
			}
		}
	}()

	go func() {

		for {
			select {

			case <-NutrientSupplyCh:
				db.Exec(" update 生产过程控制 set '任务状态'= ? where rowid= ? ", " 执行中", 1)
				NutrientSupplyCh1 <- 1
				fmt.Printf("营养液供给taskdone\n ")

			case <-NutrientSolutionConfigurationCh:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 2)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=13")
				fmt.Printf("营养液配置taskdone\n")

			case <-CoolingCh:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 3)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=14")
				fmt.Printf("降温taskdone\n")

			case <-WarmingCh:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 4)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=15")
				fmt.Printf("增温taskdone\n")

			case <-DehumidificationCh:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 5)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=16")
				fmt.Printf("除湿taskdone \n")

			case <-ShadingCh:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 6)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=19")
				fmt.Printf("遮荫taskdone \n")

			case <-CO2Ch:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 7)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=18")
				fmt.Printf("CO2taskdone\n")

			case <-VentilationCh:
				db.Exec(" update 生产过程控制 set '任务状态'=  ?  where rowid=?", "执行中", 8)
				db.Exec("  update 设备状态 set 设备状态=  '开启'  where rowid=12")
				fmt.Printf("通风taskdone \n")

			default:

				runtime.Gosched()
			}
		}
	}()

	go func(Chname chan string) {
		//var Taskname string
		//use slice,nexttime
		for {
			Taskname := <-Chname

			//   select Taskname    {
			switch Taskname {
			//select  <-Ch {
			case "营养液供给":
				//以后代码重构改为关联事物
				fmt.Printf("营养液供给taskmake\n")
				NutrientSupplyCh <- 1

			case "营养液配置":
				fmt.Printf("营养液配置taskmake1\n")

				NutrientSolutionConfigurationCh <- 1

			case "降温":
				fmt.Printf("降温taskmake\n")
				CoolingCh <- 1

			case "增温":
				fmt.Printf("增温taskmake\n")
				WarmingCh <- 1

			case "除湿":
				fmt.Printf("除湿taskmake\n")
				DehumidificationCh <- 1
			case "遮荫":
				fmt.Printf("遮荫taskmake\n")
				ShadingCh <- 1
			case "co2":
				fmt.Printf("CO2taskmake\n")
				CO2Ch <- 1
			case "通风":
				fmt.Printf("通风taskmake\n")
				VentilationCh <- 1

			default:
				//fmt.Printf("Default\n")
				runtime.Gosched()
			}
		}

	}(TkNameCh)

	//Scan需要的容器
	cols, _ := rows.Columns() //字段名

	buff := make([]interface{}, len(cols)) //Scan时的容器
	data := make([]string, len(cols))      //字段值
	//	fmt.Println(len(cols))
	//data is slice
	for i, _ := range buff {
		buff[i] = &data[i] //对容器的处理，以便能在Scan之后收集到每个字段值
	}

	for rows.Next() {
		err = rows.Scan(buff...)
		checkErr(err)
		defer rows.Close()
		fmt.Println(data)
		TkNameCh <- data[0]
	}
	<-WaitCh
	//	fmt.Println(datas)
}

func main() {
	//WaitCh := make(chan interface{})
	Task_management()
	//	time.Sleep(time.Second)

}
