package profile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestProfile_Load(t *testing.T) {
	RegisterTestingT(t)

	for _, tcase := range []struct {
		file  string
		error string
	}{
		{"", "impossible to read the profile file: open :"},
		{"./none", "impossible to read the profile file: open ./none:"},
		{"../.resources/invalid.json", "impossible to unmarshal the profile file: invalid character 'u' looking for beginning of value"},
		{"../.resources/default.json", ""},
	} {
		profile := Profile{}

		if tcase.error == "" {
			Expect(profile.Load(tcase.file)).Should(Succeed())
			Expect(profile.isLoaded).Should(BeTrue())
		} else {
			Expect(profile.Load(tcase.file)).Should(MatchError(MatchRegexp(tcase.error)))
			Expect(profile.isLoaded).Should(BeFalse())
		}
		Expect(profile.file).Should(Equal(tcase.file))
	}
}

func TestProfile_SubscribeAlteration(t *testing.T) {
	RegisterTestingT(t)

	// create unique filename
	uid := make([]byte, 16)
	rand.New(rand.NewSource(int64(time.Now().Nanosecond()))).Read(uid)
	filename := fmt.Sprintf("%x%x%x%x%x.json", uid[:4], uid[4:6], uid[6:8], uid[8:10], uid[10:])

	// create the file
	if _, err := os.Create(filename); err != nil {
		t.Fatalf("failed to create %s: %s", filename, err)
	}
	defer func() { os.Remove(filename) }()

	var prf = Profile{file: filename}
	var isAltered bool
	var mutex sync.Mutex

	// invalid subscription
	watcher, err := prf.SubscribeAlteration(nil)
	Expect(err).Should(MatchError("handler can't be nil"))

	// valid subscription
	watcher, err = prf.SubscribeAlteration(func(_ *Profile, _ error) {
		mutex.Lock()
		defer mutex.Unlock()
		isAltered = true
	})
	Expect(err).Should(Succeed())
	defer func() { watcher.Close() }()

	// goroutine to edit the file
	go func() {
		time.Sleep(time.Millisecond)
		ioutil.WriteFile(filename, uid, 0644)
	}()

	// check if the file was altered
	Eventually(func() bool {
		mutex.Lock()
		defer mutex.Unlock()
		return isAltered
	}, 25*time.Millisecond).Should(BeTrue())
}

func TestPlugin_AccessTo(t *testing.T) {
	RegisterTestingT(t)

	profile := Plugin{}
	jsonConfig := `{
  "0": {
    "0.0": "string",
    "0.1": 0,
    "0.2": 1.2,
    "0.3": true
  },
  "1": [
    "a",
    "b",
    "c"
  ],
  "2": [
    {
      "2[0].0": null,
      "2[0].1": null
    },
    {
      "2[1].0": null,
      "2[1].1": null
    }
  ]
}`
	Expect(json.Unmarshal([]byte(jsonConfig), &profile.Config)).Should(Succeed())

	for _, tcase := range []struct {
		path  []interface{}
		value interface{}
		error error
	}{
		{[]interface{}{}, nil, ErrEmptyAccessPath},
		{[]interface{}{"X"}, nil, ErrInvalidAccessPath},
		{[]interface{}{true}, nil, ErrInvalidIndexType},
		{[]interface{}{0}, nil, ErrInvalidAccessPath},
		{[]interface{}{"0", "0.4"}, nil, ErrInvalidAccessPath},
		{[]interface{}{"0", "0.0", "0.0.0"}, nil, ErrInvalidAccessPath},
		{[]interface{}{"1", 4}, nil, ErrOutOfBoundIndex},
		{[]interface{}{"2", 0, nil}, nil, ErrInvalidIndexType},

		{[]interface{}{"0", "0.0"}, "string", nil},
		{[]interface{}{"0", "0.1"}, 0.0, nil},
		{[]interface{}{"0", "0.2"}, 1.2, nil},
		{[]interface{}{"0", "0.3"}, true, nil},

		{[]interface{}{"1"}, []interface{}{"a", "b", "c"}, nil},
		{[]interface{}{"1", 0}, "a", nil},
		{[]interface{}{"1", 2}, "c", nil},

		{[]interface{}{"2", 0}, map[string]interface{}{"2[0].0": nil, "2[0].1": nil}, nil},
		{[]interface{}{"2", 0, "2[0].0"}, nil, nil},
		{[]interface{}{"2", 1, "2[1].0"}, nil, nil},
	} {
		v, err := profile.AccessTo(tcase.path...)
		if tcase.error != nil {
			Expect(err).Should(MatchError(tcase.error))
		} else {
			Expect(err).Should(Succeed())
			if tcase.value == nil {
				Expect(v).Should(BeNil())
			} else {
				Expect(v).Should(Equal(tcase.value))
			}
		}
	}
}
