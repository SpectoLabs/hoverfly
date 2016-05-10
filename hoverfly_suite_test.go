package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/phayes/freeport"
	"fmt"
	"github.com/SpectoLabs/hub/application"
	"net/http"
	"go/build"
	"strings"
	"time"
	"github.com/SpectoLabs/hoverfly"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	"github.com/SpectoLabs/hoverfly/cache"
	"github.com/dghubble/sling"
)

var (
	baseUrl string
)


func TestHoverfly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly Suite")
}

var _ = BeforeSuite(func() {
	port := freeport.GetPort()
	baseUrl = fmt.Sprintf("http://localhost:%v", port)
	go func() {
		cfg := hoverfly.InitSettings()
		db := cache.GetDB(cfg.DatabasePath)
		defer db.Close()
		requestCache := cache.NewBoltDBCache(db, []byte("requestsBucket"))
		metadataCache := cache.NewBoltDBCache(db, []byte("metadataBucket"))
		tokenCache := cache.NewBoltDBCache(db, []byte(backends.TokenBucketName))
		userCache := cache.NewBoltDBCache(db, []byte(backends.UserBucketName))
		authBackend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)
		hoverfly := hoverfly.GetNewHoverfly(cfg, requestCache, metadataCache, authBackend)
		err := hoverfly.StartProxy()
		if err != nil {
			panic(err)
		}
		go hoverfly.StartAdminInterface()
	}();
	Eventually(func() int {
		resp, err := http.Get(baseUrl + "/state")
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second * 3).Should(BeNumerically("==", http.StatusOK))
})

func DoRequest(r *sling.Sling) (*http.Response, error) {
	req, err := r.Request()
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}