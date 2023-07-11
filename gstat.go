package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	//	"sync"
)

const (
	PROXY_LIST string        = "https://www.proxy-list.download/api/v1/get?type=https"
	TIMEOUT    time.Duration = 5
)

// const F_OUT string = "out.txt"
const MAX_R int = 20

type C struct {
	Success bool
	Obj     *Company
}

var (
	verbouse, help  bool
	infile, outfile string
	threads         int
)

func init() {
	flag.BoolVar(&verbouse, "v", false, "Output log to StdOut (shorthand)")
	flag.BoolVar(&verbouse, "verbouse", false, "Output log to StdOut")
	flag.StringVar(&infile, "i", `in.txt`, "The input file")
	flag.StringVar(&outfile, "o", "out.json", "The output file")
	flag.IntVar(&threads, "t", 1, "The number of threads")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")
	flag.BoolVar(&help, "help", false, "Show help")

}

func main() {
	t := time.Now()
	defer func() {
		log.Println("Time:", time.Since(t))
	}()

	flag.Parse()
	if help {
		showHelp()
	}
	log.Println("HELP:", help)

	if !verbouse {
		output_log, err := os.OpenFile("gstat.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal("Cant open ouput file for loging ", "gstst.log", ".\n", err)
		}
		log.SetOutput(output_log)
	}

	fi, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Println("defer FI close...")
		fi.Close()
	}()
	in := bufio.NewReader(fi)

	fo, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Println("defer FO close...")
		fo.Close()
	}()

	out := bufio.NewWriter(fo)
	defer func() {
		log.Println("defer out flash:")
		out.Flush()
	}()
	fmt.Printf("Input file: %s\nOutput file: %s\nThreads: %v\n", infile, outfile, threads)
	if err = skipLines(in, fo); err == io.EOF {
		//if err = continueGet(in, fo); err == io.EOF {
		log.Printf("End of input file %v untill seeking last iin.\n", infile)
	}

	run2(in, out)
}

func run2(in *bufio.Reader, out *bufio.Writer) {
	var (
		val  any
		astr []string
		iin  string
		//mtx *sync.Mutex
		// i int
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cin := make(chan any)
	cout := make(chan any)
	pool := NewClientPool(ctx, cin, cout)
	for i := 0; i < threads; i++ {
		pool.Add(i, NewTask(i))
	}

	go func() {
		for {
			str, err := in.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			astr = strings.Split(string(str), ";")
			if len(astr) < 5 {
				continue
			}

			iin = strings.TrimSpace(astr[4])
			if iin == "" {
				continue
			}
			cin <- iin
			//fmt.Printf("send %s; ", iin)
		}
		close(cin)
		fmt.Println("\ngo End of input file... chIn closed.")
		cancel()
		fmt.Println("cancel.")

	}()

	pool.Start()

	for val = range pool.chOut {
		jd, _ := json.Marshal(val.(C).Obj)
		out.WriteString(string(jd))
		out.WriteString("\n")
		//fmt.Printf("OK IIN %v to JSON \n", val.(C).Obj.Bin)

	}
	fmt.Println("End gstat.")

}

func skipLines(r *bufio.Reader, f *os.File) error {
	count := 0
	buf := make([]byte, 2048)

	if _, err := f.Seek(-2048, 2); err != nil {
		fmt.Println("err seek", err)
		return err
	}

	if _, err := f.Read(buf); err != nil {
		fmt.Println("err resd", err)
		return err
	}
	lines := bytes.Split(buf, []byte{'\n'})
	l := ""
	for i := 1; ; i++ {
		l = strings.TrimSpace(string(lines[len(lines)-i]))
		if l != "" {
			break
		}
	}
	arr_s := strings.Split(l, ",")
	arr_iin := strings.Split(strings.TrimSpace(arr_s[0]), ":")
	last_iin := strings.Trim(arr_iin[1], "\"")
	log.Printf("last iin:%v seek next...\n", last_iin) //, iin)

	iin := ""
	var astr []string
	for {
		count++
		str, err := r.ReadString('\n')
		if err == io.EOF {
			return err
		}
		if err != nil {
			panic(err)
		}

		astr = strings.Split(string(str), ";")
		if len(astr) < 5 {
			continue
		}

		iin = strings.TrimSpace(astr[4])
		if iin == last_iin {
			break
		}
	}
	log.Printf("%v iins skiped ...\n", count) //, iin)
	return nil
}

func showHelp() {
	fmt.Printf("gstat Get Info by IIN.\n(C)2023 mrmioxin@gmail.com\ngstat get a data about individual business by the Individual Identification Number (IIN) from old.stat.gov.kz/api. If it registred.\n")
	flag.VisitAll(func(f *flag.Flag) {
		if f.DefValue == "" {
			fmt.Printf("\t-%s: %s\n", f.Name, f.Usage)
		} else {
			fmt.Printf("\t-%s: %s (Default: %s)\n", f.Name, f.Usage, f.DefValue)
		}
	})
	os.Exit(0)
}