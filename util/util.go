package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"
	"sort"
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
)

type DecodeAndValidator struct {
}

func (this *DecodeAndValidator) Validate(r *http.Request, param interface{}) error {
	if err := JsonDecode(r, param); err != nil {
		return err
	}

	if ok, err := validator.ValidateStruct(param); !ok {
		return err
	}
	return nil
}

func JsonDecode(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("[JsonDecode] body: %s", string(body))

	if string(body) != "" {
		if err := json.Unmarshal(body, v); err != nil && err != io.EOF {
			logrus.Errorf("[JsonDecode] body:%s, #error:%s", string(body), err.Error())
			return err
		}
	}

	return nil
}

// CopyStrings copies the contents of the specified string slice
// into a new slice.
func CopyStrings(s []string) []string {
	c := make([]string, len(s))
	copy(c, s)
	return c
}

// SortStrings sorts the specified string slice in place. It returns the same
// slice that was provided in order to facilitate method chaining.
func SortStrings(s []string) []string {
	sort.Strings(s)
	return s
}

// ShuffleStrings copies strings from the specified slice into a copy in random
// order. It returns a new slice.
func ShuffleStrings(s []string) []string {
	shuffled := make([]string, len(s))
	perm := rand.Perm(len(s))
	for i, j := range perm {
		shuffled[j] = s[i]
	}
	return shuffled
}

// HandleCrash simply catches a crash and logs an error. Meant to be called via defer.
func HandleCrash() {
	r := recover()
	if r != nil {
		callers := ""
		for i := 0; true; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			callers = callers + fmt.Sprintf("%v:%v\n", file, line)
		}
		logrus.Infof("Recovered from panic: %#v (%v)\n%v", r, r, callers)
	}
}

// Forever loops forever running f every d.  Catches any panics, and keeps going.
func Forever(f func(), period time.Duration) {
	for {
		func() {
			defer HandleCrash()
			f()
		}()
		time.Sleep(period)
	}
}

// MakeJSONString returns obj marshalled as a JSON string, ignoring any errors.
func MakeJSONString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
