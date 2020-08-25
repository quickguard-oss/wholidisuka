package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	tmpDir string
	cache  *cache
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	s.tmpDir = filepath.Join(os.TempDir(), "wholidisuka")

	c, err := New(s.tmpDir, true)

	if err != nil {
		panic(err)
	}

	s.cache = c
}

func (s *Suite) TearDownSuite() {
	if err := os.RemoveAll(s.tmpDir); err != nil {
		panic(err)
	}
}

func (s *Suite) TestSetGet_NotExist() {
	if err := s.cache.Clean(); err != nil {
		s.Fail("Failed to clean cache.")
	}

	got, err := s.cache.Get("test.txt", time.Second*0)

	s.Nil(got)
	s.Nil(err)
}

func (s *Suite) TestSetGet_Alive() {
	s.cache.Set("test.txt", []byte("xxx"))

	got, _ := s.cache.Get("test.txt", time.Hour*100)

	s.Equal("xxx", string(got))
}

func (s *Suite) TestSetGet_Expired() {
	s.cache.Set("test.txt", []byte("xxx"))

	got, err := s.cache.Get("test.txt", time.Second*0)

	s.Nil(got)
	s.Nil(err)
}
