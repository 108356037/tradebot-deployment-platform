package pkg

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  int64
	WriteTimeout int64
}

type HelmSettingS struct {
	ChartPath    helmChartPath
	Namespace    string
	InitYamlPath string
}

type DevS struct {
	Mode string
}

type helmChartPath struct {
	AllInOne string
	Jupyter  string
	Grafana  string
	C9       string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
