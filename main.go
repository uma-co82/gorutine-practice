package main

import (
	"fmt"
	"time"
)

func main() {
	// practice1
	practice1()

	// practice2
	practice2()
}

/***************************************************************
 *                     First Practice  Start                   *
 ***************************************************************/

/***************************************************************
 * 1つ目のfor分で go 使用
 * goroutineで実行しない場合は6秒に1行ずつ出力される
 *
 * goroutineで実行するとほぼ同時に3行(bufferが3の為)出力される
 * その後6秒待ち,bufferを解放するので次の3行が出力される
 ***************************************************************/
func practice1() {
	const totalExecuteNum = 6
	const maxConcurrencyNum = 3

	sig := make(chan string, maxConcurrencyNum)
	res := make(chan string, totalExecuteNum)

	defer close(sig)
	defer close(res)

	fmt.Printf("start concurrency execute %s\n", time.Now())

	for i := 0; i < totalExecuteNum; i++ {
		go wait6Sec(sig, res, fmt.Sprintf("No. %d", i))
	}
	for {
		if len(res) >= totalExecuteNum {
			break
		}
	}

	fmt.Printf("end concurrency execute %s \n", time.Now())

}

func wait6Sec(sig chan string, res chan string, name string) {
	sig <- fmt.Sprintf("sig %s", name)

	time.Sleep(6 * time.Second)

	fmt.Printf("%s:end wait 6 sec \n", name)

	res <- fmt.Sprintf("sig %s\n", name)

	v := <-sig

	fmt.Printf("buffer 解放 %s\n", v)
}

/***************************************************************
 *                     First Practice Fin                      *
 ***************************************************************/

/***************************************************************
 *                  Secound Practice Start                     *
 ***************************************************************/

/***************************************************************
 * main goroutine に time.Sleep(time.Second)の記述が
 * 無ければpractice2は呼ばれない
 * 理由は各goroutineより先に main goroutineが終了する為
 * main goroutineが終了した時に全体のプロセスが終了する
 * ここで言うmain goroutineはpractice2自体
 ***************************************************************/

func practice2() {
	fmt.Printf("main start\n")

	go fmt.Printf("hello No.%d", 1)
	go fmt.Printf("hello No.%d", 2)
	go fmt.Printf("hello No.%d", 3)
	// time.Sleep(time.Second) // <- ここが無ければ practice2(n) は実行されない

	fmt.Printf("main end\n")
}

/***************************************************************
 *                  Secound Practice Fin                       *
 ***************************************************************/
