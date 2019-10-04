package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	dlog "github.com/envoker/dailylog"
	dperiod "github.com/envoker/dailylog/period"
	"github.com/envoker/golang/logging/logl"
)

func main() {
	fn := logFiller
	//fn := loglFiller
	//fn := configJson
	//fn := configYaml

	if err := fn(); err != nil {
		log.Println(err)
	}
}

func configJson() error {

	a := dlog.Config{
		Dirname:        "logs",
		FilePrefix:     "test",
		FileExt:        ".log",
		KeepMaxDays:    28,
		RotateInterval: 60 * dperiod.Second,
	}

	data, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	var b dlog.Config
	err = json.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", b)

	return nil
}

func configYaml() error {

	a := dlog.Config{
		Dirname:        "logs",
		FilePrefix:     "test",
		FileExt:        ".log",
		KeepMaxDays:    28,
		RotateInterval: 24 * 60,
	}

	data, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	var b dlog.Config
	err = yaml.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", b)

	return nil
}

func logFiller() error {

	config := dlog.Config{
		Dirname:        "logs",
		FilePrefix:     "test_",
		FileExt:        ".log",
		KeepMaxDays:    28,
		RotateInterval: 24 * 60,
	}

	w, err := dlog.New(config)
	if err != nil {
		return err
	}
	defer w.Close()

	logger := log.New(w, "", log.Ldate|log.Lmicroseconds)

	const n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go logLoop(wg, logger, i)
	}
	wg.Wait()

	return nil
}

func logLoop(wg *sync.WaitGroup, logger *log.Logger, index int) {

	defer wg.Done()

	r := rand.New(rand.NewSource(int64(index)))
	for i := 0; i < 100; i++ {
		logger.Printf("routine(%d): %s", index, randString(r, 15))
		runtime.Gosched()
	}
}

func loglFiller() error {

	config := dlog.Config{
		Dirname:        "logs",
		FilePrefix:     "test_",
		FileExt:        ".log",
		KeepMaxDays:    28,
		RotateInterval: 24 * 60,
	}

	w, err := dlog.New(config)
	if err != nil {
		return err
	}
	defer w.Close()

	logger := logl.New(w, "ios ", logl.LEVEL_WARNING, logl.Lmicroseconds)

	const n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go loglLoop(wg, logger, i)
	}
	wg.Wait()

	return nil
}

func loglLoop(wg *sync.WaitGroup, logger *logl.Logger, index int) {

	defer wg.Done()

	r := rand.New(rand.NewSource(int64(index)))
	for i := 0; i < 100; i++ {

		m := fmt.Sprintf("routine(%d): %s", index, randString(r, 15))

		switch k := r.Intn(6); k {
		case 0:
			logger.Fatal(m)
		case 1:
			logger.Error(m)
		case 2:
			logger.Warning(m)
		case 3:
			logger.Info(m)
		case 4:
			logger.Debug(m)
		case 5:
			logger.Trace(m)
		}

		runtime.Gosched()
	}
}

func test1() error {

	config := dlog.Config{
		Dirname:        "logs",
		FilePrefix:     "test_",
		FileExt:        ".log",
		KeepMaxDays:    28,
		RotateInterval: 5,
	}

	w, err := dlog.New(config)
	if err != nil {
		return err
	}
	defer w.Close()

	logger := log.New(w, "record ", log.Lmicroseconds)

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	begin := time.Now()
	for {
		if time.Since(begin) > 30*time.Minute {
			break
		}
		logger.Println(randString(r, 15))
		time.Sleep(27 * time.Second)
	}

	return nil
}
