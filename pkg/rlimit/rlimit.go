package rlimit

import "syscall"

type Options struct {
	CPU          uint64
	CPUHard      uint64
	Data         uint64
	FileSize     uint64
	StackSize    uint64
	AddressSpace uint64
	OpenFile     uint64
	DisableCore  bool
}

type Option func(*Options)

func CPU(cpu uint64) Option {
	return func(opt *Options) {
		opt.CPU = cpu
	}
}

func CPUHard(cpuHard uint64) Option {
	return func(opt *Options) {
		opt.CPUHard = cpuHard
	}
}

func Data(data uint64) Option {
	return func(opt *Options) {
		opt.Data = data
	}
}

func FileSize(fileSize uint64) Option {
	return func(opt *Options) {
		opt.FileSize = fileSize
	}
}

func StackSize(stackSize uint64) Option {
	return func(opt *Options) {
		opt.StackSize = stackSize
	}
}

func AddressSpace(addressSpace uint64) Option {
	return func(opt *Options) {
		opt.AddressSpace = addressSpace
	}
}

func OpenFile(openFile uint64) Option {
	return func(opt *Options) {
		opt.OpenFile = openFile
	}
}

func DisableCore(isDisable bool) Option {
	return func(opt *Options) {
		opt.DisableCore = isDisable
	}
}

func NewRLimitOptions(opts ...Option) *Options {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type ParamHolder struct {
	Res   int
	Param syscall.Rlimit
}

func (opts *Options) PrepareRLimitHandler() []ParamHolder {
	res := make([]ParamHolder, 0)
	if opts.CPU > 0 {
		cpuHard := opts.CPUHard
		if cpuHard < opts.CPU {
			cpuHard = opts.CPU
		}
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_CPU,
			Param: syscall.Rlimit{Cur: opts.CPU, Max: cpuHard},
		})
	}

	if opts.Data > 0 {
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_DATA,
			Param: syscall.Rlimit{Cur: opts.Data, Max: opts.Data},
		})
	}

	if opts.FileSize > 0 {
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_FSIZE,
			Param: syscall.Rlimit{Cur: opts.FileSize, Max: opts.FileSize},
		})
	}

	if opts.StackSize > 0 {
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_STACK,
			Param: syscall.Rlimit{Cur: opts.StackSize, Max: opts.StackSize},
		})
	}

	if opts.AddressSpace > 0 {
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_AS,
			Param: syscall.Rlimit{Cur: opts.AddressSpace, Max: opts.AddressSpace},
		})
	}
	if opts.OpenFile > 0 {
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_NOFILE,
			Param: syscall.Rlimit{Cur: opts.OpenFile, Max: opts.OpenFile},
		})
	}
	if opts.DisableCore {
		res = append(res, ParamHolder{
			Res:   syscall.RLIMIT_CORE,
			Param: syscall.Rlimit{Cur: 0, Max: 0}})
	}

	return res
}

func SetRLimits(params []ParamHolder) error {
	for _, p := range params {
		err := syscall.Setrlimit(p.Res, &p.Param)
		if err != nil {
			return err
		}
	}
	return nil
}
