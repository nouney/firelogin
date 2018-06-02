package firelogin

type Opt func(f *Firelogin) error

func WithPort(p string) Opt {
	return func(f *Firelogin) error {
		f.port = p
		return nil
	}
}

func WithAuthHTML(h string) Opt {
	return func(f *Firelogin) error {
		f.authHTML = h
		return nil
	}
}

func WithSuccessHTML(h string) Opt {
	return func(f *Firelogin) error {
		f.successHTML = h
		return nil
	}
}
