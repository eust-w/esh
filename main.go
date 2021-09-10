package main

import (
	"esh/ssh"
	"esh/yaml"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var k = yaml.NewYaml()
var rootCmd = &cobra.Command{
	Use: "root",
}

var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"c", "clu", "C", "CLUSTER", "CL"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `use connect to connect remote ssh or run command`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("in development")
	},
}

var clusterRunCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"c", "clu", "C", "CLUSTER", "CL"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `add a cluster or add a element to the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			if len(args) <= 1 {
				fmt.Println("please enter a cluster name and commands")
				return
			}
			name, args = args[0], args[1:]
		}else if len(args) <=0{
			fmt.Println("please enter commands")
			return
		}
		command := strings.Join(args, " ")
		runFlag := strings.Trim(command, "") == ""
		names := strings.Split(k.GetCluster(name), ",")
		ch := make(chan [2]string, len(names))
		flag := make(chan bool, 1)
		for _, _name := range names {
			name := _name
			go func() {
				data, err := k.GetConn(name)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				ip, user, password, port := data[0], data[1], data[2], data[3]
				flag <- true
				if err := ssh.MultiRun(name, ip, user, password, port, command, runFlag, ch); err != nil {
					fmt.Println("Error: " + err.Error())
				}
				<-flag
			}()
		}

		for i := 0; i < len(names); i++ {
			_tem, err := <-ch
			if !err {
				close(ch)
				close(flag)
				break
			}
			fmt.Println(_tem[0], ":\n", _tem[1])
		}
		close(ch)
		close(flag)
	}}

var clusterAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"c", "clu", "C", "CLUSTER", "CL"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `add a cluster or add a element to the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Println("please enter a cluster name!")
			return
		}
		command := strings.Join(args, ",")
		k.AddCluster(name, command)
	}}

var clusterListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "L", "list", "CLUSTER", "CL"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `show all cluster name`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("|", strings.Join(k.ListCluster(), "|"), "|")
	}}

var clusterShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s", "S", "Show", "SHOW"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `show the cluster's elements'`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" && len(args) <= 0 {
			fmt.Println("please enter a cluster name!")
			return
		}
		if name == "" {
			name = args[0]
		}
		fmt.Println("|", strings.Replace(k.GetCluster(name), ",", "|", -1), "|")
	}}

var clusterDelCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"d", "Del", "DEL", "SHOW"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `show the cluster's elements'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("please enter some cluster our cluster' elements that need to del")
			return
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			k.DelCluster(args)
		}
		k.DelClusterElement(name, args)
	}}

var conCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r", "ru", "R", "RUN", "RU"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `use connect to connect remote ssh or run command`,
	Run: func(cmd *cobra.Command, args []string) {
		var ip, user, password, port string
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			ip, _ = cmd.Flags().GetString("ip")
			if ip == "" {
				name = args[0]
				args = args[1:]
				data, err := k.GetConn(name)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				ip, user, password, port = data[0], data[1], data[2], data[3]
			} else {
				user, _ = cmd.Flags().GetString("user")
				if user == "" {
					user = k.GetGlobal(yaml.DefaultUser)
				}
				password, _ = cmd.Flags().GetString("password")
				if password == "" {
					password = k.GetGlobal(yaml.DefaultPwd)
				}
				port, _ = cmd.Flags().GetString("port")
				if port == "" {
					port = k.GetGlobal(yaml.DefaultPort)
				}
			}
		} else {
			data, err := k.GetConn(name)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ip, user, password, port = data[0], data[1], data[2], data[3]
		}
		command := strings.Join(args, " ")
		runFlag := strings.Trim(command, "") == ""
		if err := ssh.Run(ip, user, password, port, command, runFlag); err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
	},
}

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "add", "ad", "AD", "A", "ADD"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   `add remote ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		ip, _ := cmd.Flags().GetString("ip")
		if name == "" || ip == "" {
			fmt.Println("please enter name and ip!")
			return
		}
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			user = k.GetGlobal(yaml.DefaultUser)
		}
		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			password = k.GetGlobal(yaml.DefaultPwd)
		}
		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			port = k.GetGlobal(yaml.DefaultPort)
		}
		k.SetConn(map[string][]string{name: {ip, user, password, port}})
		//to do
		//utils.DeleteHistory("esh add")
	},
}

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s", "set", "se", "S", "SE", "SET"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   "set global config",

	Run: func(cmd *cobra.Command, args []string) {
		data := make(map[string]string)
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		port, _ := cmd.Flags().GetString("port")
		if user != "" {
			data[yaml.DefaultUser] = user
		}
		if password != "" {
			data[yaml.DefaultPwd] = password
		}
		if port != "" {
			data[yaml.DefaultPort] = port
		}
		k.SetGlobal(data)
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li", "lis", "list", "L", "LI", "LIS", "LIST"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   "list remote ssh",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(k.ListConn()...)
	},
}

var delCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"d", "de", "del", "D", "DE", "DEL"},
	Long:    `if find bugs, please contact me www.longtao.fun`,
	Short:   "del an remote ssh use name",
	Run: func(cmd *cobra.Command, args []string) {
		k.DelConn(args)
	},
}

func init() {
	clusterAddCmd.Flags().StringP("name", "n", "", "name")
	clusterShowCmd.Flags().StringP("name", "n", "", "name")
	clusterDelCmd.Flags().StringP("name", "n", "", "name")
	clusterRunCmd.Flags().StringP("name", "n", "", "name")

	addCmd.Flags().StringP("name", "n", "", "name")
	addCmd.Flags().StringP("ip", "i", "", "ip")
	addCmd.Flags().StringP("user", "u", "", "user")
	addCmd.Flags().StringP("password", "p", "", "password")
	addCmd.Flags().StringP("port", "o", "", "port")

	conCmd.Flags().StringP("ip", "i", "", "ip")
	conCmd.Flags().StringP("user", "u", "", "user")
	conCmd.Flags().StringP("password", "p", "", "password")
	conCmd.Flags().StringP("port", "o", "", "port")

	conCmd.Flags().StringP("name", "n", "", "name")

	setCmd.Flags().StringP("port", "o", "", "default port")
	setCmd.Flags().StringP("user", "u", "", "default user name")
	setCmd.Flags().StringP("password", "p", "", "default password")
}

func main() {
	clusterCmd.AddCommand(clusterAddCmd, clusterShowCmd, clusterListCmd, clusterDelCmd, clusterRunCmd)
	rootCmd.AddCommand(addCmd, conCmd, listCmd, delCmd, setCmd, clusterCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
