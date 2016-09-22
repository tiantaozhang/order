package orderalipay

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"os"
	"sync"
)

func DoPerf(tc int, logf string, call func(int64)) (int64, error) {
	stdout := os.Stdout
	stderr := os.Stderr
	if len(logf) > 0 {
		f, err := os.OpenFile(logf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return 0, err
		}
		os.Stdout = f
		os.Stderr = f
		log.SetWriter(f)
	}
	ws := sync.WaitGroup{}
	ws.Add(tc)
	beg := util.Now()
	for i := 0; i < tc; i++ {
		go func(v int64) {
			call(v)
			ws.Done()
		}(int64(i))
	}
	ws.Wait()
	end := util.Now()
	if len(logf) > 0 {
		os.Stdout.Close()
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetWriter(os.Stdout)
	}
	return end - beg, nil
}
