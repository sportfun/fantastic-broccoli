package profile

import (
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
	watcher, err = prf.SubscribeAlteration(func(_ *Profile) {
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
