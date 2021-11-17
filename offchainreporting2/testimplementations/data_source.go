package testimplementations

import (
	"context"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/tidwall/gjson"
)

var _ median.DataSource = (*RandomDataSource)(nil)

type RandomDataSource struct {
	mu   sync.Mutex
	rand *rand.Rand
}

func NewRandomDataSource(seed int64) *RandomDataSource {
	return &RandomDataSource{
		rand: rand.New(rand.NewSource(seed)),
		mu:   sync.Mutex{},
	}
}

var uInt192StrictUpperBound = (&big.Int{}).Lsh(big.NewInt(1), 192)
var int192StrictUpperBound = (&big.Int{}).Rsh(uInt192StrictUpperBound, 1)

func (r *RandomDataSource) Observe(context.Context) (*big.Int, error) {
	// Want a uniform sample from range of int192. Rand gives sample from
	// [0,2**192), subtracting 2**191 gives sample from [-2**191,2**191).
	r.mu.Lock()
	defer r.mu.Unlock()
	return (&big.Int{}).Sub(
		(&big.Int{}).Rand(r.rand, uInt192StrictUpperBound),
		int192StrictUpperBound,
	), nil
}

var _ median.DataSource = (*SimpleFetchingDataSource)(nil)

type SimpleFetchingDataSource struct {
	url  string
	path string
}

func NewSimpleFetchingDataSource(url string, path string) SimpleFetchingDataSource {
	return SimpleFetchingDataSource{
		url:  url,
		path: path,
	}
}

func (s SimpleFetchingDataSource) Observe(ctx context.Context) (*big.Int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("status code %d while querying data source %s",
			resp.StatusCode, s.url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	value := gjson.GetBytes(body, s.path)
	obs := big.NewInt(int64(value.Float()))
	return obs, nil
}

// This data source reads the current value of a file on disk.
type FileDataSource struct {
	path string
}

func NewFileDataSource(path string) FileDataSource {
	return FileDataSource{
		path: path,
	}
}

func (s FileDataSource) Observe(ctx context.Context) (*big.Int, error) {
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, err
	}
	obs := big.NewInt(i)
	return obs, nil
}

type ConstantDataSource struct {
	value *big.Int
}

func NewConstantDataSource(value *big.Int) *ConstantDataSource {
	return &ConstantDataSource{value}
}

func (s ConstantDataSource) Observe(ctx context.Context) (*big.Int, error) {
	result := new(big.Int)
	result.Set(s.value)
	return result, nil
}
