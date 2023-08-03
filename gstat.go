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
	"os/signal"
	"strings"
	"time"
	//	"sync"
)

const (
	PROXY_LIST string        = "https://www.proxy-list.download/api/v1/get?type=https"
	TIMEOUT    time.Duration = 5
	MAX_R      int           = 10
)

type C struct {
	Success bool
	Obj     *Company
}

var (
	verbouse        bool
	infile, outfile string
	threads         int
)

func init() {
	flag.BoolVar(&verbouse, "v", false, "Output log to StdOut (shorthand)")
	flag.BoolVar(&verbouse, "verbouse", false, "Output log to StdOut")
	flag.StringVar(&infile, "i", `in.txt`, "The input file")
	flag.StringVar(&outfile, "o", "out.json", "The output file")
	flag.IntVar(&threads, "t", 1, "The number of threads")
	flag.Usage = showHelp

}

func main() {
	t := time.Now()
	defer func() {
		fmt.Println("End gstat. Work time:", time.Since(t))
		if panicVal := recover(); panicVal != nil {
			log.Fatalf("Stop programm because %v:\n", panicVal)
		}
	}()

	flagParse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	fi, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	fo, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	out := bufio.NewWriter(fo)
	defer out.Flush()

	if count, err := skipLines(fi, fo); err == io.EOF {
		//if err = continueGet(in, fo); err == io.EOF {
		log.Printf("End of input file %v until seeking last iin. Skip %d lines\n", infile, count)
	} else if err != nil {
		log.Println("SkipLines:", err)
	} else {
		log.Printf("%v iins skiped ...\n", count)
	}

	cin := make(chan any)
	cout := make(chan any)

	go sentDataForTasks(fi, cin)

	pool := NewClientPool(threads, cin, cout)
	go pool.Start(ctx)

	for val := range cout {
		jd, err := json.Marshal(val.(C).Obj)
		if err != nil {
			log.Println("error marshall.", val.(C).Obj)
		} else {
			out.WriteString(string(jd))
			out.WriteString("\n")
		}
	}
}

func sentDataForTasks(fin io.Reader, cin chan<- any) { //read input data to in chanal
	var astr []string
	var iin string
	in := bufio.NewReader(fin)

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
		log.Printf("iin %v to chan.", iin)
		cin <- iin
		//fmt.Printf("send %s; ", iin)
	}
	close(cin)
	fmt.Println("\ngo End of input file... chIn closed.")
	//app.done()
	fmt.Println("App cancel by the end of input.")
}

func getLastIIN(f io.ReadSeeker) (string, error) {
	buf := make([]byte, 2048)
	if _, err := f.Seek(-2048, 2); err != nil {
		return "", err
	}

	if _, err := f.Read(buf); err != nil {
		return "", err
	}

	lines := bytes.Split(buf, []byte{'\n'})
	l := ""
	for i := len(lines) - 1; ; i-- {
		l = strings.TrimSpace(string(lines[i]))
		if l != "" && l[len(l)-1] == '}' {
			break
		}
	}
	arr_s := strings.Split(l, ",")
	arr_iin := strings.Split(strings.TrimSpace(arr_s[0]), ":")
	last_iin := strings.Trim(arr_iin[1], "\"")
	return last_iin, nil
}

func skipLines(fi io.Reader, fo io.ReadSeeker) (int, error) {
	last_iin, err := getLastIIN(fo)
	if err != nil {
		return 0, fmt.Errorf("error getLastiin:%v", err)
	}

	count := 0
	log.Printf("last iin:%v seek next...\n", last_iin) //, iin)
	r := bufio.NewReader(fi)

	iin := ""
	var astr []string
	for {
		count++
		str, err := r.ReadString('\n')
		if err == io.EOF {
			return count, err
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
	return count, nil
}

func flagParse() {

	flag.Parse()

	if !verbouse {
		output_log, err := os.OpenFile("gstat.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("Cant open file for loging gstat.log. Output log to STDOUT\n")
		} else {
			log.SetOutput(output_log)
		}
	}
	fmt.Printf("Input file: %s\nOutput file: %s\nThreads: %v\n", infile, outfile, threads)
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
