package cmd

import (
	"github/go-project/tour/internal/sql2struct"
	"log"

	"github.com/spf13/cobra"
)

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql",
	Long:  "sql",
	Run:   func(cmd *cobra.Command, Args []string) {},
}

var sql2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql struct",
	Long:  "sql struct",
	Run: func(cmd *cobra.Command, args []string) {
		dbinfo := &sql2struct.DBInfo{
			DBType:   dbType,
			Host:     host,
			Userame:  username,
			Password: password,
			Charset:  charset,
		}
		dbModel := sql2struct.NewDBModel(dbinfo)
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel connect error : %v", err)
		}
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel GetColumns error : %v", err)
		}

		template := sql2struct.NewStructTemplate()
		tplCols := template.AssemblyColumns(columns)
		err = template.Generate(tableName, tplCols)
		if err != nil {
			log.Fatalf("template generate err: %v", err)
		}
	},
}

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

func init() {
	sqlCmd.AddCommand(sql2structCmd)
	sql2structCmd.Flags().StringVarP(&username, "username", "u", "root", "please input username !")
	sql2structCmd.Flags().StringVarP(&password, "password", "p", "123456", "please input password !")
	sql2structCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1:3306", "please input host !")
	sql2structCmd.Flags().StringVarP(&charset, "charset", "", "utf8mb4", "please input charset !")
	sql2structCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "please input dbType !")
	sql2structCmd.Flags().StringVarP(&dbName, "db", "", "mysql", "please input dbName !")
	sql2structCmd.Flags().StringVarP(&tableName, "table", "", "general_log", "please input tableName !")
}
