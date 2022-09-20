package yaml

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var path = "/etc/esh.conf"

func init() {
	rootPath, _ := os.UserHomeDir()
	path = filepath.Join(rootPath, "esh.conf")
}

const (
	ConfigOK = iota
	ConfigNotExist
	DefaultUser = "default_user"
	DefaultPwd  = "default_password"
	DefaultPort = "default_port"
)

type Yaml struct {
	Global  map[string]string `yaml:"global"`
	Conn    map[string]string `yaml:"conn"`
	Cluster map[string]string `yaml:"cluster"`
}

func (y *Yaml) checkAndInit() int {
	file, err := os.Open(path)
	defer func() {
		y.init()
	}()
	if err != nil {
		if !os.IsExist(err) {
			file, err = os.OpenFile(path, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, 0777)
			os.Chmod(path, 0777)
			return ConfigNotExist
		} else {
			log.Println(err)
		}
	}
	defer func() { file.Close() }()
	return ConfigOK
}

func (y *Yaml) init() {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temGlobal := y.Global
	if temGlobal == nil {
		temGlobal = make(map[string]string)
	}
	if temGlobal[DefaultUser] == "" {
		temGlobal[DefaultUser] = "root"
	}
	if temGlobal[DefaultPwd] == "" {
		temGlobal[DefaultPwd] = "password"
	}
	if temGlobal[DefaultPort] == "" {
		temGlobal[DefaultPort] = "22"
	}
	y.Global = temGlobal
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
}

func (y Yaml) GetConn(k string) ([]string, error) {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, &y)
	return AesDecrypt(y.Conn[k])
}

func (y *Yaml) SetConn(conns map[string][]string) bool {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temConn := y.Conn
	if temConn == nil {
		temConn = make(map[string]string)
	}
	for k, v := range conns {
		if len(v) != 4 {
			log.Println("setConn error, len(conns) != 3,len(conns)=", len(v), k)
			return false
		}
		data := fmt.Sprintf(strings.Join(v, "\n"))
		temConn[k] = AesEncrypt(data)
	}
	y.Conn = temConn
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
	return true
}

func (y *Yaml) DelConn(conns []string) {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temConn := y.Conn
	if temConn == nil {
		return
	}
	for _, k := range conns {
		delete(temConn, k)
	}
	y.Conn = temConn
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
	return
}

func (y Yaml) ListConn() (list []interface{}) {

	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, &y)
	for k, _ := range y.Conn {
		list = append(list, k)
	}
	return
}

func (y Yaml) GetGlobal(k string) string {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, &y)
	return y.Global[k]
}

func (y *Yaml) SetGlobal(k map[string]string) bool {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temGlobal := y.Global
	if temGlobal == nil {
		temGlobal = make(map[string]string)
	}
	for k, v := range k {
		temGlobal[k] = v
	}
	y.Global = temGlobal
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
	return true
}

func (y Yaml) GetCluster(c string) string {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, &y)
	return y.Cluster[c]
}

func (y *Yaml) DelCluster(c []string) {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temCluster := y.Cluster
	if temCluster == nil {
		return
	}
	for _, k := range c {
		delete(temCluster, k)
	}
	y.Cluster = temCluster
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
	return
}

func (y Yaml) ListCluster() (list []string) {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, &y)
	for k, _ := range y.Cluster {
		list = append(list, k)
	}
	return
}

func (y *Yaml) AddCluster(c string, l string) {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temCluster := y.Cluster
	if temCluster == nil {
		temCluster = make(map[string]string)
	}
	if temCluster[c] == "" {
		temCluster[c] = l
	} else {
		temCluster[c] = strings.Join(union(strings.Split(temCluster[c], ","), strings.Split(l, ",")), ",")
	}
	y.Cluster = temCluster
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
	return
}

func (y *Yaml) DelClusterElement(c string, l []string) {
	yamlFile, _ := ioutil.ReadFile(path)
	yaml.Unmarshal(yamlFile, y)
	temCluster := y.Cluster
	if temCluster == nil {
		return
	}
	if temCluster[c] == "" {
		return
	} else {
		temCluster[c] = strings.Join(difference(strings.Split(temCluster[c], ","), l), ",")
	}
	y.Cluster = temCluster
	amendFile, _ := yaml.Marshal(y)
	ioutil.WriteFile(path, amendFile, 077)
	return
}

func NewYaml() *Yaml {
	_yaml := Yaml{}
	status := _yaml.checkAndInit()
	if status == ConfigNotExist {
		_, err := os.Open(path)
		if os.IsExist(err) {
			log.Println("No yaml ,Success create in ", path)
		} else {
			log.Println("No yaml！ ,Create in path fail! Please try again with administrator privileges！")
		}
	}
	return &_yaml
}

//返回交集
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//返回并集
func union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		m[v]++
	}
	for k, _ := range m {
		nn = append(nn, k)
	}
	return nn
}

//返回差集 slice1-slice2
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
