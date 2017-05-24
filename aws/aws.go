package aws

import (
	"sync"
)

var (
	af     *Funcs
	afInit sync.Once
)

// NS - the aws namespace
func NS() *Funcs {
	afInit.Do(func() { af = &Funcs{} })
	return af
}

// AddFuncs -
func AddFuncs(f map[string]interface{}) {
	f["aws"] = NS

	// global aliases - for backwards compatibility
	f["ec2meta"] = NS().EC2Meta
	f["ec2dynamic"] = NS().EC2Dynamic
	f["ec2tag"] = NS().EC2Tag
	f["ec2region"] = NS().EC2Region
}

// Funcs -
type Funcs struct {
	meta     *Ec2Meta
	metaInit sync.Once
	info     *Ec2Info
	infoInit sync.Once
}

// EC2Region -
func (a *Funcs) EC2Region(def ...string) string {
	a.metaInit.Do(a.initMeta)
	return a.meta.Region(def...)
}

// EC2Meta -
func (a *Funcs) EC2Meta(key string, def ...string) string {
	a.metaInit.Do(a.initMeta)
	return a.meta.Meta(key, def...)
}

// EC2Dynamic -
func (a *Funcs) EC2Dynamic(key string, def ...string) string {
	a.metaInit.Do(a.initMeta)
	return a.meta.Dynamic(key, def...)
}

// EC2Tag -
func (a *Funcs) EC2Tag(tag string, def ...string) string {
	a.infoInit.Do(a.initInfo)
	return a.info.Tag(tag, def...)
}

func (a *Funcs) initMeta() {
	if a.meta == nil {
		a.meta = NewEc2Meta()
	}
}

func (a *Funcs) initInfo() {
	if a.info == nil {
		a.info = NewEc2Info()
	}
}
