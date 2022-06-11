package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Setting  配置结构.
type Setting struct {
	vp *viper.Viper
}

// NewSetting 读取配置文件内容创建配置对象.
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()

	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	s := &Setting{vp}
	s.WatchSettingChange()

	return s, nil
}

// WatchSettingChange  监听配置文件变化.
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
