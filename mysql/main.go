package main
//mysql数据库操作使用，by：yq 2018-04-03
//基本CURD都操作使用了一遍，使用便捷程度还是比较高，整体操作还需要封装，下一步看ORM的性能和原生的比较
//有时间做做大数据量性能，和dotnet做一下对比
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "root:yaoqi1717@tcp(139.199.177.131:13306)/yq_godb")
	defer db.Close()
	check(err)
	//insert(db)
	selectdb(db)
	selectMax(db)
}

func insert(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO t_users(username, password,age,creat_time,remark) VALUES(?,?,?,?,?)")
	defer stmt.Close()
	check(err)
	stmt.Exec("yaoqi", "yaoqipass", 20, time.Now(), "mark")
	// timer1 := time.NewTimer(time.Second * 2)
	// for index := 0; index < 1000; index++ {
	// 	stmt.Exec("yaoqi", "yaoqi11",20,time.Now(),index)
	//     stmt.Exec("testuser", "123123",23,time.Now(),index)
	// }
}

func selectdb(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM t_users limit 100")
    defer rows.Close()
	check(err)
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
}
func del(db *sql.DB, _id int) {
	stmr, err := db.Prepare("delete from t_users where id = ?")
	defer stmr.Close() //一定注意关闭，编写优质代码的defer
	check(err)
	res, err := stmr.Exec(_id)
	check(err)
	fmt.Println(res)
	num, err := res.RowsAffected() //返回影响行
    check(err)
    fmt.Println(num)
}


func update(db *sql.DB, _id int) {
	stmr, err := db.Prepare("update t_users set password=? where id = ?")
	defer stmr.Close() //一定注意关闭，编写优质代码的defer
	check(err)
	res, err := stmr.Exec("init",_id)
	check(err)
	fmt.Println(res)
	num, err := res.RowsAffected() //返回影响行
    check(err)
    fmt.Println(num)
}


func selectMax(db *sql.DB){
	var _maxid int
	err := db.QueryRow("select id from t_users ORDER BY id DESC LIMIT 1").Scan(&_maxid)
	check(err)
	fmt.Println(_maxid)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
