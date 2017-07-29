package bacon

type Options struct {
	allMeat bool
}

type OptionFn func(o *Options) *Options

// SetAllMeat returns an OptionFn that sets all meat or meat and filler
func SetAllMeat(allMeat bool) OptionFn {
	return func(o *Options) *Options {
		o.allMeat = allMeat
		return o
	}
}

func defaultOptions() *Options {
	return &Options{
		allMeat: false,
	}
}
