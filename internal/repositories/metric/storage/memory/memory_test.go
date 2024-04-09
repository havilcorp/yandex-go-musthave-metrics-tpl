package memory

import (
	"context"
	"math/rand"
	"testing"

	_ "net/http/pprof"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/stretchr/testify/require"
)

type WantGauge struct {
	key     string
	val     float64
	wantKey string
	wantVal float64
	isError bool
}

type WantCounter struct {
	key     string
	val     int64
	wantKey string
	wantVal int64
	isError bool
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func BenchmarkAddGauge(b *testing.B) {
	b.ReportAllocs()
	store := NewMemStorage()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		key := randSeq(5)
		val := rand.Float64()
		store.AddGauge(ctx, key, val)
	}
}

func BenchmarkAddCounter(b *testing.B) {
	b.ReportAllocs()
	store := NewMemStorage()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		key := randSeq(5)
		val := rand.Int63()
		store.AddCounter(ctx, key, val)
	}
}

func ExampleMemStorage_AddGauge() {
	store := NewMemStorage()
	store.AddGauge(context.Background(), "KEY", 1.1)
}

func ExampleMemStorage_AddCounter() {
	store := NewMemStorage()
	store.AddCounter(context.Background(), "KEY", 1)
}

func ExampleMemStorage_AddGaugeBulk() {
	store := NewMemStorage()
	list := make([]domain.Gauge, 1)
	list = append(list, domain.Gauge{
		Key:   "Alloc",
		Value: 100.123,
	})
	store.AddGaugeBulk(context.Background(), list)
}

func ExampleMemStorage_AddCounterBulk() {
	store := NewMemStorage()
	list := make([]domain.Counter, 1)
	list = append(list, domain.Counter{
		Key:   "Counter",
		Value: 100,
	})
	store.AddCounterBulk(context.Background(), list)
}

func TestAddGauge(t *testing.T) {
	tests := []struct {
		name string
		want WantGauge
	}{
		{
			name: "test1",
			want: WantGauge{
				key:     "Alloc",
				val:     100.654,
				wantKey: "Alloc",
				wantVal: 100.654,
				isError: false,
			},
		},
		{
			name: "test2",
			want: WantGauge{
				key:     "Alloc",
				val:     100.654,
				wantKey: "Undefined",
				wantVal: 100.654,
				isError: true,
			},
		},
	}
	store := NewMemStorage()
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.AddGauge(ctx, test.want.key, test.want.val)
			if err != nil {
				panic(err)
			}
			val, err := store.GetGauge(ctx, test.want.wantKey)
			require.Equal(t, err != nil, test.want.isError)
			if test.want.isError == false {
				require.Equal(t, test.want.wantVal, val)
			}
		})
	}
}

func TestAddGaugeBulk(t *testing.T) {
	tests := []struct {
		name string
		want [2]WantGauge
	}{
		{
			name: "test1",
			want: [2]WantGauge{
				{
					key:     "Alloc",
					val:     100.123,
					wantKey: "Alloc",
					wantVal: 100.123,
				},
				{
					key:     "Alloc2",
					val:     100.2,
					wantKey: "Alloc2",
					wantVal: 100.2,
				},
			},
		},
	}
	store := NewMemStorage()
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gauge := make([]domain.Gauge, 0)
			gauge = append(gauge, domain.Gauge{
				Key:   test.want[0].key,
				Value: test.want[0].val,
			})
			gauge = append(gauge, domain.Gauge{
				Key:   test.want[1].key,
				Value: test.want[1].val,
			})
			err := store.AddGaugeBulk(ctx, gauge)
			if err != nil {
				panic(err)
			}
			var val0 float64
			var val1 float64
			list, err := store.GetAllGauge(ctx)
			if err != nil {
				t.Error(err)
			}
			for k, v := range list {
				if test.want[0].key == k {
					val0 = v
				}
				if test.want[1].key == k {
					val1 = v
				}
			}
			if test.want[0].isError == false {
				require.Equal(t, test.want[0].wantVal, val0)
			}
			if test.want[1].isError == false {
				require.Equal(t, test.want[1].wantVal, val1)
			}
		})
	}
}

func TestAddCounter(t *testing.T) {
	tests := []struct {
		name string
		want WantCounter
	}{
		{
			name: "test1",
			want: WantCounter{
				key:     "Alloc",
				val:     100,
				wantKey: "Alloc",
				wantVal: 100,
				isError: false,
			},
		},
		{
			name: "test2",
			want: WantCounter{
				key:     "Alloc",
				val:     100,
				wantKey: "Undefined",
				wantVal: 100,
				isError: true,
			},
		},
	}
	store := NewMemStorage()
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.AddCounter(ctx, test.want.key, test.want.val)
			if err != nil {
				panic(err)
			}
			val, err := store.GetCounter(ctx, test.want.wantKey)
			require.Equal(t, err != nil, test.want.isError)
			if test.want.isError == false {
				require.Equal(t, test.want.wantVal, val)
			}
		})
	}
}

func TestAddCounerBulk(t *testing.T) {
	tests := []struct {
		name string
		want [2]WantCounter
	}{
		{
			name: "test1",
			want: [2]WantCounter{
				{
					key:     "Alloc",
					val:     100,
					wantKey: "Alloc",
					wantVal: 100,
				},
				{
					key:     "Alloc2",
					val:     12,
					wantKey: "Alloc2",
					wantVal: 12,
				},
			},
		},
	}
	store := NewMemStorage()
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counter := make([]domain.Counter, 0)
			counter = append(counter, domain.Counter{
				Key:   test.want[0].key,
				Value: test.want[0].val,
			})
			counter = append(counter, domain.Counter{
				Key:   test.want[1].key,
				Value: test.want[1].val,
			})
			err := store.AddCounterBulk(ctx, counter)
			if err != nil {
				panic(err)
			}
			var val0 int64
			var val1 int64
			list, err := store.GetAllCounters(ctx)
			if err != nil {
				t.Error(err)
			}
			for k, v := range list {
				if test.want[0].key == k {
					val0 = v
				}
				if test.want[1].key == k {
					val1 = v
				}
			}
			if test.want[0].isError == false {
				require.Equal(t, test.want[0].wantVal, val0)
			}
			if test.want[1].isError == false {
				require.Equal(t, test.want[1].wantVal, val1)
			}
		})
	}
}
