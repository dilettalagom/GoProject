package rpcservice

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"sync"
)

type Args struct {
	File            string
	N               int
	M               int
	Offset          int
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
func wordcount(wg *sync.WaitGroup, offset int, block_size int, buf []string, channel_reducer *chan map[string]int) {

	defer wg.Done()

	var word_map = make(map[string]int)

	for i := 0; i < len(buf); i++ {

		fmt.Printf("Chiave: %s\n", buf[i])

		key := buf[i]

		_, ok := word_map[key]
		if !ok {
			word_map[key] = 1
		} else {
			word_map[key] += 1
		}
	}

	/*for key, value := range word_map {
		fmt.Println("Key:", key, "Value:", value)
	}
	fmt.Println("\n")*/

	//TODO: send to channel_reducer[j]
	for key, _ := range word_map {
		var sum uint8 = 0
		for j := 0; j < len(key); j++ {
			sum += byte(key[j])
		}
		if len(*channel_reducer) != 0 {
			fmt.Printf("hashing %d\n", int(sum)%len(*channel_reducer))
		}
	}

}

func (w *Wordcounter) Map(args Args, Result *Result) error {

	var wg sync.WaitGroup

	file_stream, err := os.Open(args.File)
	check(err)

	//init channel
	chan_tmp := make(chan map[string]int)
	args.Chan_to_reducer = &chan_tmp

	//getFileSize
	file_info, err := os.Stat(args.File)
	check(err)
	size := file_info.Size()
	if size == 0 {
		os.Exit(1)
	}

	//buffer for file
	buf := make([]byte, size)

	//read the entire file
	buf, err = ioutil.ReadFile(args.File)
	if err != nil {
		log.Fatal(err)
	}

	//closing file --> only for stream, no need to close is though

	err = file_stream.Close()
	if err != nil {
		log.Fatal(err)
	}

	//TODO: create M thread and put them waiting on a barrier

	//get an array of words and compute the number of words for each thread
	list_of_words := strings.Split(string(buf), " ")
	//number_of_words := math.Ceil(float64(len(list_of_words))/float64(num_worker))
	number_of_words := int(float64(len(list_of_words)) / float64(args.N))

	fmt.Printf("All words and number of words for each thread: %d - %d\n", len(list_of_words), number_of_words)

	for i := 0; i < args.N; i += 1 {

		wg.Add(1)

		//compute the offset from which each thread starts and the number of words to count
		args.Offset = i * int(number_of_words)
		words_to_read := int(math.Min(float64(number_of_words), float64(int(len(list_of_words))-(i*int(number_of_words)))))

		fmt.Printf("Quanto legge e da dove parte: %d --- %d\n", words_to_read, args.Offset)

		if i == args.N-1 {
			//last thread may read more than words_to_read words
			go wordcount(&wg, args.Offset, words_to_read, list_of_words[args.Offset:], args.Chan_to_reducer)
		} else {
			go wordcount(&wg, args.Offset, words_to_read, list_of_words[args.Offset:(i+1)*words_to_read], args.Chan_to_reducer)
		}

		wg.Wait()

	}

	return nil

}
