/*
需要用到go-odbc库，下载地址：https://github.com/weigj/go-odbc
64位的WIN不能使用ODBC来调用access数据库了，
所以必须要安装32位的GCC和32位的GO来安装go-odbc。
然后要用32位的GO来编译本程序才行。。。。
*/
package main

import (
	"fmt"
	"os"
	"io"
	"flag"
	"database/sql" 
	_"github.com/weigj/odbc/driver"
)

func d(err error,s string){
	if (err != nil){
		fmt.Println(s," . err :",err.Error())
		os.Exit(1)	
	}
}

func main() {
	var(
		//db *sql.DB
		db2 *sql.DB
		//stmt *sql.Stmt
		//rows *sql.Rows
		result sql.Result
		err error
		csql string
	)
	flag.Parse()
	v3 := flag.Arg(0)
	println("输入参数:",v3)
	println("请先确认目录内的cfg.mdb为3.37版本数据库文件！")
	println("检查cfg.mdb文件版本是否为3.37！");
	/*3.37版有tuser表，没有userinfo表*/
	/*db,err = sql.Open("odbc","driver={Microsoft Access Driver (*.mdb)};dbq=cfg.mdb");
	d(err,"sql.Open")
	csql = "select * from tuser";
	_,err = db.Query(csql);
	if(err != nil){
		//说明不存在tuser表
		d(err,"cfg.mdb文件不是3.37版本")
	}else{
		//说明存在tuser表

	}
	db.Close();*/
	err = os.Rename("cfg.mdb", "cfg-3.37.mdb")
	println("备份原cfg.mdb文件！")
	d(err,"备份文件时出错")
	/*****************************************/
	//db,err := sql.Open("odbc", "DSN=jyd");	//linux
	//db,err := sql.Open("odbc","Driver=MDBTools;DBQ=/www/jyd/1.mdb");	//linux
	//db, err := sql.Open("odbc", "DSN=21;")

	db2,err = sql.Open("odbc","driver={Microsoft Access Driver (*.mdb)};dbq=v5");
	d(err,"致命错误，是否V5.5版的cfg.mdb文件没有了？：")
	defer db2.Close()
	if e := db2.Ping(); e != nil {
		println("V5.5数据库连接错误！",e)
		os.Exit(2)
	}
	csql = "DELETE * FROM treeinfolocal"     
	_,err = db2.Exec(csql);
	println("执行，清空V5.5版本的表数据:",err)
	
	csql = `select id,
		Caption,parentid,loginacc,
		loginpsd,dvsip,dvsport,dvschannel,dvslinkmode 
		into treeinfolocal in 'd:\\cfg.mdb' from treeinfo `;
		/*复制到新的表中*/
	csql = `insert into treeinfolocal(id,Caption,parentid,loginacc,loginpsd,
			dvsip,dvsport,dvschannel,dvslinkmode)  
			select id,Caption,parentid,loginacc,loginpsd,
			dvsip,dvsport,dvschannel,dvslinkmode from treeinfo in 'd:\\1.mdb' `;
		/*复制到已有的表中dvslinkmode=0,selfdvskind='DVS',sipnodekind='EU'
		insert into tableB IN '" & toDBFile & "' select * FROM tableA where userName='abc*/
	/*db,err = sql.Open("odbc","driver={Microsoft Access Driver (*.mdb)};dbq=d:\\1.mdb");
	d(err,"sql.Open")
	defer db.Close()*/
	result,err = db2.Exec(csql)
	if err != nil {
    	println("拷贝V3.37版本数据到V5.5版本数据库...",err)
	}else{
		println("计算拷贝数据总数...")
		if c,err := result.RowsAffected(); err != nil{
			println("拷贝数据总数：",c)
		}else{
			println("计算拷贝数据总数时出错：",err)
		}
	}
	
	//csql = "update treeinfolocal set selfdvskind=? WHERE parentid<>-1 ";
	//_,err = db2.Exec(csql,"DVS");
	//println("更新V5.5数据库...",err)
	csql = `update treeinfolocal set selfdvskind='',belone=10 WHERE parentid=-1 `;
	_,err = db2.Exec(csql);
	println("更新V5.5数据库...",err)
	csql = "update treeinfolocal set sipnodekind='DO',dvslinkmode=1 ";
	_,err = db2.Exec(csql);
	println("更新V5.5数据库...",err)
	csql = "update treeinfolocal set sipnodekind='EU',dvslinkmode=0 WHERE dvsip<>'' ";
	_,err = db2.Exec(csql);
	println("更新V5.5数据库...",err)
	csql = "update treeinfolocal set sipnodekind='EU',dvslinkmode=0 WHERE dvsip<>'' ";
	_,err = db2.Exec(csql);
	println("更新V5.5数据库...",err)

	println("开始生成Cfg.mdb文件...")
	w, err := CopyFile("v5", "cfg.mdb")  
    if err != nil {  
        println("生成Cfg.mdb文件出错：",err.Error())  
    }  
    println("生成Cfg.mdb文件完成！字节数：",w)

	/*stmt, err = db.Prepare("select id,Caption from Treeinfo where id > ?")
	d(err,"db.Prepare")
	defer stmt.Close()

	rows, err = stmt.Query("1")
	d(err,"stmt.Query")
	defer rows.Close()
	
	for rows.Next() {
		var (
		Caption string
		id int
		)
		if err = rows.Scan(&id,&Caption); err == nil{
			//fmt.Println("Caption:\t",Caption,"\t id:",id)
		}else{
			fmt.Println("err:",err)
		}
   }*/
}
func CopyFile(src, des string) (w int64, err error) {  
    srcFile, err := os.Open(src)  
    if err != nil {  
        println(err)  
    }  
    defer srcFile.Close()  
  
    desFile, err := os.Create(des)  
    if err != nil {  
        println(err)  
    }  
    defer desFile.Close()  
  
    return io.Copy(desFile, srcFile)  
}  