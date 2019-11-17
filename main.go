package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/http"
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

	// practice19()

	// practice20()

	// practice21()

	// practice22()

	// practice23()

	// practice24()

	// practice25()

	// practice26()

	// practice27()

	// practice28()

	// practice29()

	// practice30()

	// practice31()

	// practice32()

	// practice33()

	// practice34()

	// practice35()

	// practice36()

	// practice37()

	// practice38()

	// practice39()

	// practice40()

	// practice41()

	// practice42()

	// practice43()

	// practice44()

	primeNum()

	// practice45()
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
 * オブジェクトプールパターンは使うものを決まった数だけ、作る方法(EX.)DB接続)
 * プール内に使用可能なインスタンスがあるか確認
 * あれば呼び出し元に返す
 * なければNewメンバー変数を呼び出し、新しいインスタンスを作成
 * 作業が終われば、呼び出し元はPutを呼んで、使っていたインスタンスをプールに戻す
 * そうすると他のプロセスが使える
 ***************************************************************/

func practice20() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

func practice21() {
	var numClassCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numClassCreated++
			mem := make([]byte, 1024)
			return &mem
		},
	}

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numClassCreated)
}

/***************************************************************
 * チャネル
 * ゴルーチン間の通信に使うのが最適
 * 変数名の末尾を"Stream","Ch","c"にするのが良い
 * チャネルを扱う時は常に初期化すること！
 ***************************************************************/

func practice22() {
	// 送受信
	var dataStream chan interface{}
	dataStream = make(chan interface{})

	// 受信のみ
	var reserveStream <-chan interface{}
	reserveStream = make(<-chan interface{})

	// 送信のみ
	var sendStream chan<- interface{}
	sendStream = make(chan<- interface{})

	// 正しい記述
	reserveStream = dataStream
	sendStream = dataStream

	fmt.Printf("%v", reserveStream)
	fmt.Printf("%v", sendStream)
}

func practice23() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels!"
	}()
	// stringStreamに何か入ってくるまでブロックして待っている
	fmt.Println(<-stringStream)

	valueStream := make(chan int)
	// チャネルを閉じる
	close(valueStream)
	integer, ok := <-valueStream
	fmt.Printf("(%v): %v", ok, integer)

	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	// rangeはチャンネルが閉じた時に自動的にループを終了する
	for integer := range intStream {
		fmt.Printf("%v\n", integer)
	}
}

// 複数のゴルーチン一度に解放する
func practice24() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

// バッファ付きチャネル
func practice25() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	// intStreamに4つ入るまでブロックされる
	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v. \n", integer)
	}
}

// チャネルの権限をゴルーチン毎に明確にする
// 1. チャネルを初期化
// 2. 書き込みを行うか、他のゴルーチンに所有権を渡します。
// 3. チャネルを閉じる
// 4. 上の3つの手順をカプセル化して読み込みチャネルを経由して公開
func practice26() {
	// 受信専用のチャネルを返す
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

/***************************************************************
 * select文
 * 上から順番に評価されない
 * 読み込みや書き込みのチャネルは全て同時に取り扱われ
 * どれが用意できたかを確認する
 * どのチャネルも準備できていない場合はselect全体がブロックする
 ***************************************************************/

func practice27() {
	var c1, c2 <-chan interface{}
	var c3 chan<- interface{}

	select {
	case <-c1:
		// 何かする
	case <-c2:
		// 何かする
	case c3 <- struct{}{}:
		// 何かする
	}
}

// selectブロックに入ってからおよそ5秒後にブロックをやめる
// case <-cでclose()を感知
func practice28() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

// 複数のチャネルが同時に読み込める様になった時
// case文全体ではそれぞれの条件が等しく選択される可能性がある
func practice29() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

// 全てのチャネルがブロックされた時
func practice30() {
	var c <-chan int
	select {
	case <-c: // このcase分はnilチャンネルから読み込んでいるので決してブロックが解放されない
	case <-time.After(1 * time.Second): // タイムアウトを実現する
		fmt.Println("Timed out.")
	}
}

// チャネルが１つも読み込めず、その間に何かする必要ある時
func practice31() {
	start := time.Now()
	var c1, c2 <-chan int

	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}

// default節はfo-selectループの中で使われる
// ゴルーチンの結果報告を待つ間に他のゴルーチンで仕事を進められます
func practice32() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second) //5秒間スリープした後closeする
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done: // 5秒間スリープした後close()を感知
			break loop
		default:
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achived %v cycles of work before signalled to stop.\n", workCounter)
}

/***************************************************************
 * プリミティブ紹介
 ***************************************************************/

/***************************************************************
 * 拘束
 * パフォーマンス向上、可読性向上、同期の問題を考えなくて良い
 * - イミュータブルなデータ
 * - 拘束によって保護されたデータ
 ***************************************************************/

/***************************************************************
 * アドホック拘束
 * 拘束を規約によって達成する  <- 実際ムリ！！
 ***************************************************************/

func practice33() {
	data := make([]int, 4)

	// loopDataからしかdataにアクセスしない様な規約があったとする
	// しかし、人が触る物なのでいつミスが起こるか分からない
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

/***************************************************************
 * レキシカル拘束
 * コンパイラによって拘束を強制する
 * - チャネルへの読み書きのうち必要な権限だけ公開する
 ***************************************************************/

func practice34() {
	// resultsチャネルへの書き込みができるスコープを制限
	// 書き込み権限を拘束
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	// 読み込み用に拘束
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	// 読み込み権限
	results := chanOwner()
	consumer(results)
}

func practice35() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	// それぞれのごルーチンがdataの一部にしかアクセスできない様に拘束している
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}

/***************************************************************
 * for-selectループ
 * 外部から停止の命令が来るまで無限に繰り返すゴルーチンを作る事が多い
 ***************************************************************/

// doneチャネルが閉じられていなければ、select文を抜けてforループの本体の残りの処理を続ける。
// (doneチャネルが閉じられるまで永遠に続ける処理)
// defaultの中に書いてもおk
func practice36() {
	done := make(chan int)
	for {
		select {
		case <-done:
			return
		default:
		}
		// (doneチャネルが閉じられるまで永遠に続ける処理)
	}
}

/***************************************************************
 * ゴルーチンリークを避ける
 * ゴルーチンはランタイムによってガベージコレクションされないので片付けたい
 *
 * もしあるゴルーチンがゴルーチンの生成の責任を持っているのであれば
 * そのゴルーチンを停止できるようにする責任もある
 ***************************************************************/

// doWorkを含むごルーチンはこのプロセスが生きている限りずっとメモリ内に残る
func practice37() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return completed
	}

	doWork(nil)

	fmt.Println("Done.")
}

// ゴルーチンの親子間で親から子にキャンセルのシグナルを送れるようにする
// このシグナルは慣習として、doneという名前の読み込み専用チャネルにする
// 読み込み版
func practice38() {
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done: // close(done)を感知する
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}

// 書き込み版
func practice39() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done: // close(done)を感知
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 0; i <= 3; i++ {
		fmt.Printf("%v: %v\n", i, <-randStream)
	}
	close(done)
}

/***************************************************************
 * orチャネル
 * 1つ以上のdoneチャネルを1つのdoneチャネルにまとめ、まとめてるチャネルのうちの
 * どれか１つのチャネルが閉じられたら、まとめたチャネルも閉じる
 ***************************************************************/

// 任意の数のチャネルを１つのチャネルにまとめる事ができる
// まとめているチャネルのどれか１つでも閉じたり書き込まれたら、すぐに合成
// されたチャネルが閉じる
func practice40() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	// 異なる時間待機してチャネルが閉じられるが、、、同時に閉じる
	<-or(sig(2*time.Hour), sig(5*time.Minute), sig(1*time.Second), sig(1*time.Hour), sig(1*time.Minute))
	fmt.Printf("done after %v", time.Since(start))
}

/***************************************************************
 * エラーハンドリング
 ***************************************************************/

func practice41() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

/***************************************************************
 * パイプライン
 * システムの抽象化に使う
 * データを受け取って、何らかの処理を行って、どこかに渡すという一連の作業に過ぎない
 * こららの操作をパイプラインのステージと呼ぶ
 * パイプラインを使うことで、各ステージでの懸念事項を切り分けられる
 *
 * - ステージは受け取るものと返すものが同じ型
 * - ステージは引き回せるように具体化されてなければならない(golangで意識しない)
 ***************************************************************/

func practice42() {
	// ステージの役割
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	// ステージの役割
	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	// パイプライン
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

/***************************************************************
 * ！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
 * ！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
 * !！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！重要！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
 * ！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
 * ！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
 ***************************************************************/
/***************************************************************
 * パイプライン構築のベストプラクティス
 ***************************************************************/

// practice42をチャネルを使って書き換える
func practice43() {
	// パイプラインの始めには常にチャネルへの変換を必要とするデータの塊がある
	// ジェネレーターと呼ばれる
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int, len(integers))
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done: // ゴルーチンリークを避ける
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done) // 各ゴルーチンをキルする

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

/***************************************************************
 * 便利ジェネレーター
 * パイプラインのジェネレーターは一塊の値をチャネル上のストリームに変換する関数
 ***************************************************************/

func practice44() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}

	rand := func() interface{} {
		return rand.Int()
	}

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}

	fmt.Printf("message: %s...", message)
}

/***************************************************************
 * ファンアウト、ファンイン
 * パイプライン内のあるステージで計算量が大きく、処理に時間が掛かる場合
 * パイプライン全体の実行に長い時間がかかってしまいます。
 ***************************************************************/

// 並行処理で素数出してみた
func primeNum() {
	generator := func(done <-chan interface{}, num int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()

		return intStream
	}

	hoge := func(num int) bool {
		if num == 1 || num == 2 {
			return true
		}

		for i := 2; i < num; i++ {
			if num%i == 0 {
				return false
			}
		}

		return true
	}

	primeFinder := func(done <-chan interface{}, valueStream <-chan int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				if hoge(v) {
					intStream <- v
				}
			}
		}()

		return intStream
	}

	done := make(chan interface{})
	defer close(done)

	for v := range primeFinder(done, generator(done, 100)) {
		fmt.Println(v)
	}
}
