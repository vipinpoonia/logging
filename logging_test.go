package logging

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

func func1() error {
	err := errors.New("TestError")
	return err
}

func func2() error {
	err := func1()
	log.Info(err)
	return err
}
func test() {
	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	err := func2()
	log.WithError(err).Error()
	log.WithError(nil).Error()
	// Calls os.Exit(1) after logging
	//log.Fatal("Bye.")
	// Calls panic() after logging
	//log.Panic("I'm bailing.")

	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
	time.Sleep(5 * time.Second)
	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())

}

func TestLogging(t *testing.T) {
	config := LogConfig{
		SentryDsn: os.Getenv("SentryDsn"),
		Env:       "local",
	}
	log.Info("Initializing logging..")

	err := Init(config)
	test()
	if err != nil {
		t.Errorf(err.Error())
	}

}
