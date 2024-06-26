package rlimit

import (
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
)

func TestNewRLimitOptions(t *testing.T) {
	opts := NewRLimitOptions(
		WithCPU(1),
		WithCPUHard(2),
		WithData(3),
		WithFileSize(4),
		WithStackSize(5),
		WithAddressSpace(6),
		WithOpenFile(7),
		WithDisableCore(true),
	)

	assert.Equal(t, uint64(1), opts.CPU)
	assert.Equal(t, uint64(2), opts.CPUHard)
	assert.Equal(t, uint64(3), opts.Data)
	assert.Equal(t, uint64(4), opts.FileSize)
	assert.Equal(t, uint64(5), opts.StackSize)
	assert.Equal(t, uint64(6), opts.AddressSpace)
	assert.Equal(t, uint64(7), opts.OpenFile)
	assert.Equal(t, true, opts.DisableCore)
}

func TestRLimits_PrepareRLimit(t *testing.T) {
	opts := &Options{
		CPU:          1,
		CPUHard:      2,
		Data:         3,
		FileSize:     4,
		StackSize:    5,
		AddressSpace: 6,
		OpenFile:     7,
		DisableCore:  true,
	}

	ret := opts.PrepareRLimitHandler()
	assert.Equal(t, 7, len(ret))
	assert.Equal(t, uint64(1), ret[0].Param.Cur)
	assert.Equal(t, uint64(2), ret[0].Param.Max)
	assert.Equal(t, uint64(3), ret[1].Param.Cur)
	assert.Equal(t, uint64(3), ret[1].Param.Max)
	assert.Equal(t, uint64(4), ret[2].Param.Cur)
	assert.Equal(t, uint64(4), ret[2].Param.Max)
	assert.Equal(t, uint64(5), ret[3].Param.Cur)
	assert.Equal(t, uint64(5), ret[3].Param.Max)
	assert.Equal(t, uint64(6), ret[4].Param.Cur)
	assert.Equal(t, uint64(6), ret[4].Param.Max)
	assert.Equal(t, uint64(7), ret[5].Param.Cur)
	assert.Equal(t, uint64(7), ret[5].Param.Max)
	assert.Equal(t, uint64(0), ret[6].Param.Cur)
	assert.Equal(t, uint64(0), ret[6].Param.Max)
}

func TestSetRLimits(t *testing.T) {
	opts := &Options{
		CPU:          1,
		CPUHard:      2,
		Data:         3,
		FileSize:     4,
		StackSize:    5,
		AddressSpace: 6,
		OpenFile:     7,
		DisableCore:  true,
	}

	handlers := opts.PrepareRLimitHandler()

	err := SetRLimits(handlers)
	assert.Nil(t, err)

	rlimit := &syscall.Rlimit{}

	_ = syscall.Getrlimit(syscall.RLIMIT_CPU, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(1))
	assert.Equal(t, rlimit.Max, uint64(2))

	_ = syscall.Getrlimit(syscall.RLIMIT_DATA, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(3))
	assert.Equal(t, rlimit.Max, uint64(3))

	_ = syscall.Getrlimit(syscall.RLIMIT_FSIZE, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(4))
	assert.Equal(t, rlimit.Max, uint64(4))

	_ = syscall.Getrlimit(syscall.RLIMIT_STACK, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(5))
	assert.Equal(t, rlimit.Max, uint64(5))

	_ = syscall.Getrlimit(syscall.RLIMIT_AS, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(6))
	assert.Equal(t, rlimit.Max, uint64(6))

	_ = syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(7))
	assert.Equal(t, rlimit.Max, uint64(7))

	_ = syscall.Getrlimit(syscall.RLIMIT_CORE, rlimit)
	assert.Equal(t, rlimit.Cur, uint64(0))
	assert.Equal(t, rlimit.Max, uint64(0))
}
