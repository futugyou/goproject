### 前言： 最近把这本书从头到尾敲了一遍，记录下其中一些组件的使用方式。
##### 这本书的github原地址为 [go-programming-tour-book](https://github.com/go-programming-tour-book)

### 第一章： 使用flag和cobra实现简单命令行工具
1. #### flag 基本命令行
    <details>
    <summary> 点击展开 </summary>

    ```golang
    //go run main.go   --name=09   go --name=7655
    var nameFlag Name
    flag.Var(&nameFlag, "name", "help info")//声明一个参数 09
    flag.Parse()
    	
    goCmd := flag.NewFlagSet("go", flag.ExitOnError)//一个新的子命令 go
    goCmd.StringVar(&name, "name", "go project", "help info")//子命令的参数 7655
    phpCmd := flag.NewFlagSet("php", flag.ExitOnError)//另一个新的子命令 php
    phpCmd.StringVar(&name, "n", "php project", "help info")
    
    args := flag.Args()
    switch args[0] {
    case "go":
        _ = goCmd.Parse(args[1:])
    case "php":
        _ = phpCmd.Parse(args[1:])
    }
    ```
    </details>

2. #### cobra 命令行
    ```golang
    go get -u github.com/spf13/cobra //从来不加version的我在grpc那个项目给自己挖了个坑
    ```
    <details>
    <summary>1. 建立一个root空命令</summary>
    
    ```golang
    var rootCmd = &cobra.Command{
    	Use:   "",
    	Short: "",
    	Long:  "",
    	Run: func(cmd *cobra.Command, args []string) {
    	},
    }
    
    // Execute Execute  在main中调用此函数
    func Execute() error {
    	return rootCmd.Execute()
    }
    
    func init() {
    	rootCmd.AddCommand(wordCmd) // 这三个都是root的子命令,只贴一个word
    	rootCmd.AddCommand(timeCmd) 
    	rootCmd.AddCommand(sqlCmd) // 这个涉及了template和sql的基本使用
    }
    ```
    </details>
    
    <details>
    <summary>2. 建立一个word子命令</summary>
    
    ```golang
    var str string //俩参数
    var mode int8
    var wordCmd = &cobra.Command{
    	Use:   "word",  // 关键字
    	Short: "change word", // short和long都是说明
    	Long:  desc,
    	Run: func(cmd *cobra.Command, args []string) {
            var content string
            ... //具体内容就不贴了
        }
    }
       
    func init() {
        //两个参数  go run main.go word --str=hello  --mode=0
    	wordCmd.Flags().StringVarP(&str, "str", "s", "", "please input word !")
    	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "please intout change mode !")
    }

    ```
    </details>