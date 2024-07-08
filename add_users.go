package vmtools

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	//"gopkg.in/yaml.v3"
	"github.com/go-faster/yaml"
)

type Config struct {
	Input io.Reader
	Output io.Writer
	indent int
	encoder yaml.Encoder
	Users []User
	Header string
	YamlString string
}
type Opt func(*Config)

type User struct {
	Username string `yaml:"username"`
	Ip string `yaml:"vm_ip"`
}

type AdditionalUsers struct {
	Users []User `yaml:"additional_users"`
}


func NewConfig(opts ...Opt) *Config{
	c := &Config{
		Input: os.Stdin,
		Output: os.Stdout,
		indent: 2,
	}
	for _, opt := range opts{
		opt(c)
	}
	return c
}

func SetIndent(indent int) func(*Config) {
	return func(c *Config){
		c.indent = indent
	}
}

func WithInput(buf io.Reader) func(*Config) {
	return func(c *Config){
		c.Input = buf
	}
}
func WithOutput(buf io.Writer) func(*Config) {
	return func(c *Config){
		c.Output = buf
	}
}
func WithHeader(header string) func(*Config){
	return func(c *Config){
		c.Header = header
	}
}
func (c *Config) GetIndent() int {
	return c.indent
}

func CreateUser(username string, ip string) User {
	u := User{Username:username,Ip:ip}
	return u
}

func (c *Config) GetUsers(ips []string) {
	var usernames []string
	var users []User
	scanner := bufio.NewScanner(c.Input)
	
	for scanner.Scan(){
		usernames = append(usernames,scanner.Text())
	}
	for _, ip := range ips {
		for _, user := range usernames {
			users = append(users, CreateUser(user, ip))
		}
	}
	c.Users = users
}
func (c *Config) GenerateYaml() (string,error) {
	var b bytes.Buffer
	additionalusers := AdditionalUsers{Users:[]User{}}
	for _, user := range c.Users{
		additionalusers.Users = append(additionalusers.Users,user)
	}
	encoder := yaml.NewEncoder(&b)
	defer encoder.Close()
	encoder.SetIndent(c.indent)
	err := encoder.Encode(&additionalusers)
	if err != nil {
		return "",err
	}
	c.YamlString = b.String()
	
	return c.YamlString,nil
}

func (c *Config) WriteYaml() error{
	if len(c.YamlString) == 0 {
		return errors.New("Uninitialized yaml string")
	}
	if len(c.Header) > 0 {
		c.Output.Write([]byte(c.Header))
		c.Output.Write([]byte{'\n'})
	}
	c.Output.Write([]byte(c.YamlString))
	return nil
}