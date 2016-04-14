package main

type config struct {
	Debug    bool
	Database string
	Password string

	Add struct {
		AccessKey       string
		SecretAccessKey string
		Region          string
	}

	Console struct {
		Duration int
	}

	Shell struct {
		Shell  string
		Export bool
	}

	LockAgent struct {
		Timeout string
	}
}
