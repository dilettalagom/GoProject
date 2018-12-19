package rpcservice

import (
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

type Args struct {
	File            string
	N               int
	M               int
	Offset          int
	Num_rows        int
	Chan_to_reducer *chan map[string]int
}

type Result map[string]int

type Wordcounter int

//utils
func check(err error) {
	if err != nil {
		panic(err)
	}
}

//goroutines
func mapper(wg *sync.WaitGroup, file *os.File, offset int, block_size int, channel_reducer *chan map[string]int) {

	defer wg.Done()

	//change seek
	_, err := file.Seek(int64(offset), io.SeekStart)
	check(err)

	//fmt.Printf("seek %d\n", seek)
	//creating word_map
	var word_map = make(map[string]int)

	//Split of file
	/*scanner := bufio.NewScanner(file)

	for scanner.Scan()  {
		for i:=0; i < block_size; i++{
			line := scanner.Text()
			//fmt.Printf("riga %s\n", line)
			parts := strings.Split(line, " ")
			for _, newKey := range parts {
				//creating word_map
				_, ok := word_map[newKey]
				if (!ok) {
					word_map[newKey] = 1
				}else{
					word_map[newKey] += 1
				}
			}
		}
		break
	}*/

	//TODO: read byte per byte

	//TEST
	for key, value := range word_map {
		fmt.Println("Key:", key, "Value:", value)
	}
	fmt.Println("\n")

	//TODO:per ogni chiave trasforma in byte e fa hash
	//invio sul canale indicizzato dall'hashing

}

func (w *Wordcounter) Map(args Args, Result *Result) error {

	var wg sync.WaitGroup

	file_stream, err := os.Open(args.File)
	check(err)

	//init channel
	chan_tmp := make(chan map[string]int)
	args.Chan_to_reducer = &chan_tmp

	//getFileSize
	file_info, err := os.Stat(file_stream.Name())
	check(err)
	size := file_info.Size()
	if size == 0 {
		os.Exit(1)
	}

	//TODO: goroutines REDUCER-> wait barrier +wg.add(M)

	//compute partitions
	num_rows := int(math.Ceil(float64(int(size) / args.N)))

	for i := 0; i < num_rows; i += 1 {

		size_tmp := int(math.Min(float64(num_rows), float64(int(size)-(i*num_rows)))) //amount to read

		remainder := int(float64(int(size) - (i * num_rows)))

		args.Offset = i * num_rows

		if i == num_rows-1 {
			args.Num_rows = remainder
			//TODO: goroutines MAPPER
			wg.Add(1)
			fmt.Printf("ULTIMOOO offset - size %d\n", args.Offset, remainder)
			go mapper(&wg, file_stream, args.Offset, remainder, args.Chan_to_reducer)
		} else {
			args.Num_rows = size_tmp
			//TODO: goroutines MAPPER
			wg.Add(1)
			fmt.Printf("offset - size %d\n", args.Offset, size_tmp)
			go mapper(&wg, file_stream, args.Offset, size_tmp, args.Chan_to_reducer)
		}
		//TODO: wait of REDUCERS

		wg.Wait()

	}

	return nil

}
