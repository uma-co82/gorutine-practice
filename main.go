package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// practice1
	// practice1()

	// practice2
	// practice2()

	// practice3
	practice3()
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
	// time.Sleep(time.Second) // <- ここが無ければ 各goroutine は実行されない

	fmt.Printf("main end\n")
}

/***************************************************************
 *                  Secound Practice Fin                       *
 ***************************************************************/

/***************************************************************
 *                  Third Practice Start                       *
 ***************************************************************/

/***************************************************************
 * Shared-memory :メモリを共有して通信をやり取りする
 * 結果は `hello5` が5回表示される
 * i が 5になった段階でtime.Sleep(time.Second)が呼ばれmain goroutine
 * が止まり,その間に各goroutineが実行される為
 *
 * 以下で回避可能
 * go func(i int) {
 * 	fmt.Println("hello", i)
 * }(i)
 ***************************************************************/

func practice3() {
	fmt.Printf("main start\n")

	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println("hello", i)
		}()
	}
	time.Sleep(time.Second)

	fmt.Printf("main end\n")
}

/***************************************************************
 *                  Third Practice  Fin 		                   *
 ***************************************************************/

/***************************************************************
 *                  Fourth Practice Start		                   *
 ***************************************************************/

/***************************************************************
 * Shared-memory :メモリを共有して通信をやり取りする
 * sync.Mutex.Lock()によってメモリへのアクセスを取得
 * sync.Mutex.Unlock()で解放
 * *ちなみにこのコードでも実行される順序は非決定的
 ***************************************************************/

func practice4() {
	var memoryAccess sync.Mutex
	var data int

	go func() {
		memoryAccess.Lock()
		data++
		memoryAccess.Unlock()
	}()

	memoryAccess.Lock()
	if data == 0 {
		fmt.Printf("the value is 0\n")
	} else {
		fmt.Printf("the value is %v\n", data)
	}
	memoryAccess.Unlock()
}
