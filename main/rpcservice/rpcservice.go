package rpcservice

import (
	"./barrier"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

type Args struct {
	File   string
	N      int
	M      int
	Offset int
}

type Result map[string]int

type Wordcounter int

//goroutines
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
			//fmt.Printf("hashing %d\n", c_index)
			m := map[string]int{key: word_map[key]}
			(*channel_reducer)[c_index] <- m
		}
	}

}

func reducer(br *barrier.Barrier, wg *sync.WaitGroup, channel *chan map[string]int, channel_daddy *chan map[string]int) {

	defer wg.Done()

	//starts its life waiting at the barrier
	br.Wait_on_barrier()

	map_recieved := make(map[string]int)

readChannel:
	for {
		select {
		//has been received a message_map on this channel
		case m_tmp := <-*channel:
			for key, value := range m_tmp {
				map_recieved[key] += value
			}

		//timeout expired
		case <-time.After(10 * time.Second): //TODO: decidere il valore del TIMEOUT
			close(*channel)
			break readChannel
		}
	}

	//return values to father
	*channel_daddy <- map_recieved

	/*TEST:
	fmt.Printf("\n\n REDUCER: IN TOTO HO RICEVUTO\n")
	for key, value := range map_recieved {
		fmt.Println("Key:", key, "Value:", value)
	}
	fmt.Printf("REDUCER: muoio\n\n")*/

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

	//read the entire file
	buf, err := ioutil.ReadFile(args.File)
	if err != nil {
		log.Fatal(err)
	}

	//create M thread_reducer and put them waiting on a barrier with specific channel
	for i := 0; i < args.M; i++ {
		wg.Add(1)
		go reducer(br, &wg, &chan_to_reducers[i], &chan_backTo_daddy[i]) //listening only on a single channel
	}

	//get an array of words and compute the number of words for each thread
	list_of_words := strings.Split(string(buf), " ") //TODO: bisogna controllare il formato del file: se ci sono \n lui non li elimina e crea una chiave < parola \n parola >

	number_of_words := int(float64(len(list_of_words)) / float64(args.N))
	fmt.Println("All words and number:", len(list_of_words), "and len of words for each thread:", number_of_words)

	for i := 0; i < args.N; i += 1 {

		wg.Add(1)

		//compute the offset from which each thread starts and the number of words to count
		args.Offset = i * int(number_of_words)
		words_to_read := int(math.Min(float64(number_of_words), float64(int(len(list_of_words))-(i*int(number_of_words)))))

		//TEST: fmt.Println("Quanto legge e da dove parte: ", words_to_read, " --- ", args.Offset)

		if i == args.N-1 {
			//last thread may read more than words_to_read words
			go mapper(&wg, br, list_of_words[args.Offset:], &chan_to_reducers)
		} else {
			go mapper(&wg, br, list_of_words[args.Offset:(i+1)*words_to_read], &chan_to_reducers)
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
