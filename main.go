package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func main() {
	// practice1()

	// practice2()

	// practice3()

	// practice4()

	// practice5()

	// practice6()

	// practice7()

	// practice8()

	// practice9()

	// practice10()

	// practice11()
}

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
 * Shared-memory :メモリを共有して通信をやり取りする
 * sync.Mutex.Lock()によってメモリへのアクセスを取得
 * sync.Mutex.Unlock()で解放
 * ちなみにこのコードでも実行される順序は非決定的
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

/***************************************************************
 * Shared-memory :メモリを共有して通信をやり取りする
 * デッドロック発生
 * 各goroutineで a, b をロックしてその後 b, aをロックしようと無限に
 * 待ち続けている状態 * ちなみに time.Slee()を外せばデッドロックは発生しない
 ***************************************************************/

func practice5() {
	type value struct {
		mu    sync.Mutex
		value int
	}
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()
		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()
		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}
	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

/***************************************************************
 * ライブロック
 * 廊下のすれ違いを避ける為両方のgoroutineが左右に移動している状態
 ***************************************************************/

func practice6() {
	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", dirName)
		atomic.AddInt32(dir, 1)
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprintf(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}

/***************************************************************
 * リソース枯渇
 * greedyWorkerはワークループ全体で共有ロックを保持している
 * politeWorkerは必要な時だけロックしている
 ***************************************************************/

func practice7() {
	var wg sync.WaitGroup
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}
		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}

/***************************************************************
 * 合流ポイント
 * time.Sleepを訂正した正しい例
 ***************************************************************/

func practice8() {
	var wg sync.WaitGroup
	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello")
	}

	wg.Add(1)
	go sayHello()
	wg.Wait()
}

func practice9() {
	var wg sync.WaitGroup
	salutation := "hello"

	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()

	fmt.Println(salutation)
}

func practice10() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		// goroutineが開始する前にforによるループが終了する
		// コピーする事で回避可能
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

/***************************************************************
 * コンテキストスイッチベンチマーク
 ***************************************************************/

func practice11(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}

/***************************************************************
 * syncパッケージ
 ***************************************************************/

/***************************************************************
 * WaitGroup
 * 結果を気にしない、もしくは他に結果を収集する手段がある場合
 * ひとまとまりの並行処理の完了を待つのに向いている
 ***************************************************************/

func practice12() {
	var wg sync.WaitGroup

	wg.Add(1) // Addの引数に1を渡して、１つのゴルーチンが起動した事を表している
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done() // WaitGroupに終了する事を伝える
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	wg.Wait() // Waitは全てのゴルーチンが終了したと伝えるまでメインゴルーチンをブロックする
	fmt.Println("All goroutine complete.")
}
