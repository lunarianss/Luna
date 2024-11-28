package common

type I18nObject struct {
	Zh_Hans string `json:"zh_Hans" yaml:"zh_Hans"`
	En_US   string `json:"en_US"   yaml:"en_US"`
}

func (o *I18nObject) PatchZh() {
	if o == nil {
		return
	}
	if o.Zh_Hans == "" {
		o.Zh_Hans = o.En_US
	}
}
