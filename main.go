package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"text/tabwriter"
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

	// practice12()

	// practice13()

	// practice14()

	// practice15()

	// practice16()

	// practice17()

	// practice18()

	practice19()
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
 * Addの呼び出しは監視対象のgoroutineの外で行う事
 *
 * 出来る限り監視対象のgoroutineの直前でAddを呼ぶ事
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

func practice13() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}

/***************************************************************
 * Mutex
 * メモリへのアクセスを同期する
 * クリティカルセクションを保護する

 * Unlockの呼び出しはdeferで行う事がイディオム
 ***************************************************************/

func practice14() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()         //Mutexインスタンスで保護されたクリティカルセクション
		defer lock.Unlock() // 保護したクリティカルセクションの処理が終了した
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	var arithmethis sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmethis.Add(1)
		go func() {
			defer arithmethis.Done()
			increment()
		}()
	}

	for i := 0; i <= 5; i++ {
		arithmethis.Add(1)
		go func() {
			defer arithmethis.Done()
			decrement()
		}()
	}

	arithmethis.Wait()

	fmt.Println("Arithmetic complete.")
}

/***************************************************************
 * RWMutex
 * メモリへのアクセスを同期する
 * クリティカルセクションを保護する
 * Mutexより多くメモリの管理を提供してくれる。
 * 書き込みのロックを要求した場合読み込みのロックが取れる
 *
 * 論理的に意味があると思うときはMutexでは無く、基本RWMutex
 ***************************************************************/

func practice15() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}

/***************************************************************
 * Cond
 * 2つ以上のごルーチン間で、「それ」が発生したというシグナル
 * ゴルーチン上で処理を続ける前にシグナルを受け取りたい時に使う
 * Waitの呼び出しはブロックするだけでなく現在のゴルーチンを一時停止する
 ***************************************************************/

func practice16() {
	conditionTrue := func() bool {
		time.Sleep(3)
		return false
	}
	c := sync.NewCond(&sync.Mutex{})
	c.L.Lock() // Waitへの呼び出しがループに入るときに自動的にUnlockを呼び出すため必要
	for conditionTrue() == false {
		c.Wait()
	}
	c.L.Unlock() // Waitの呼び出しが終わるとLockを呼び出す
}

/***************************************************************
 * Cond - Broadcast
 * シグナルを待っているすべてのゴルーチンにシグナルを伝える。
 * Signalはシグナルを一番長く待っているゴルーチンを見つけ、そのゴルーチンに伝える
 ***************************************************************/

func practice17() {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutieRunning sync.WaitGroup
		goroutieRunning.Add(1)
		go func() {
			goroutieRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutieRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}

/***************************************************************
 * Once
 * Do()に渡された関数がたとえ異なるゴルーチンで呼ばれたとしても、
 * 一度しか実行されないようにする型
 ***************************************************************/

func practice18() {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

/***************************************************************
 * OnceはDo()が呼び出された回数だけを数えていて、Doに渡された一意な関数が
 * 呼び出された回数を数えているわけではない
 *
 * syncパッケージ内の型を使うときは狭い範囲が最適
 ***************************************************************/

func practice19() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }

	var once sync.Once
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count)
}

/***************************************************************
 * Pool
 ***************************************************************/
