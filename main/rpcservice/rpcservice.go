package rpcservice

import (
	"./barrier"
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Args struct {
	File string
	N    int
	M    int
}

type Result map[string]int

type Wordcounter int

/* mapper routine*/
func mapper(wg *sync.WaitGroup, br *barrier.Barrier, buf []string, channel_reducer *[]chan map[string]int) {

	defer wg.Done()

	word_map := make(map[string]int)

	for i := 0; i < len(buf); i++ {

		key := buf[i]

		_, ok := word_map[key]
		if !ok {
			word_map[key] = 1
		} else {
			word_map[key] += 1
		}
	}

	//waiting all the other mappers
	br.Wait_on_barrier()

	for key, _ := range word_map {
		var sum uint8 = 0
		for j := 0; j < len(key); j++ {
			sum += byte(key[j])
		}

		if len(*channel_reducer) != 0 {
			c_index := int(sum) % len(*channel_reducer)
			m := map[string]int{key: word_map[key]}
			(*channel_reducer)[c_index] <- m
		}
	}

}

/* reducer routine*/
func reducer(br *barrier.Barrier, wg *sync.WaitGroup, channel *chan map[string]int, channel_daddy *chan map[string]int) {

	defer wg.Done()

	//starts its life waiting at the barrier
	br.Wait_on_barrier()

	map_received := make(map[string]int)

readChannel:
	for {
		select {
		//has been received a message_map on this channel
		case m_tmp := <-*channel:
			for key, value := range m_tmp {
				map_received[key] += value
			}

		//timeout expired
		case <-time.After(5 * time.Second):
			break readChannel
		}
	}

	//return values to father
	*channel_daddy <- map_received

}

func (w *Wordcounter) Map(args Args, Result *Result) error {

	var wg sync.WaitGroup

	//init channells mappers->reducer
	chan_to_reducers := []chan map[string]int{}
	for k := 0; k < args.M; k++ {
		tmp := make(chan map[string]int, 10000)
		chan_to_reducers = append(chan_to_reducers, tmp)
	}

	//init channels reducers->father
	chan_backTo_daddy := []chan map[string]int{}
	for k := 0; k < args.M; k++ {
		tmp := make(chan map[string]int, 10000)
		chan_backTo_daddy = append(chan_backTo_daddy, tmp)
	}

	//creating barrier
	br := barrier.New(args.N + args.M)

	//create M thread_reducer and put them waiting on a barrier with specific channel
	for i := 0; i < args.M; i++ {
		wg.Add(1)
		go reducer(br, &wg, &chan_to_reducers[i], &chan_backTo_daddy[i]) //listening only on a single channel
	}

	file_to_chunk, err := os.Open(args.File)
	if err != nil {
		panic(err)
	}

	var list_of_words []string
	scanner := bufio.NewScanner(file_to_chunk)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		for i := range parts {
			list_of_words = append(list_of_words, strings.TrimFunc(parts[i], func(r rune) bool {
				return !unicode.IsLetter(r)
			}))

		}
	}

	number_of_words := int(float64(len(list_of_words)) / float64(args.N))
	fmt.Println("All words and number:", len(list_of_words), "and len of words for each thread:", number_of_words)

	offset := 0

	for i := 0; i < args.N; i += 1 {

		wg.Add(1)

		//compute the offset from which each thread starts and the number of words to count
		offset = i * int(number_of_words)
		words_to_read := int(math.Min(float64(number_of_words), float64(int(len(list_of_words))-(i*int(number_of_words)))))

		if i == args.N-1 {
			//last thread may read more than words_to_read words
			go mapper(&wg, br, list_of_words[offset:], &chan_to_reducers)
		} else {
			go mapper(&wg, br, list_of_words[offset:(i+1)*words_to_read], &chan_to_reducers)
		}
	}

	wg.Wait()

	//merging all children's results
	for child := 0; child < args.M; child++ {
		m_tmp := <-chan_backTo_daddy[child]

		for key, value := range m_tmp {
			(*Result)[key] += value
		}
	}

	return nil

}
