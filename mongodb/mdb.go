package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

const url = "mongodb://root:root@10.0.0.60:27017/"

type Operater struct {
	mogSession *mgo.Session
	dbname     string
	document   string
}

type person struct {
	AGE    int    `bson:"age"`
	NAME   string `bson:"name"`
	HEIGHT int    `bson:"height"`
	CITY   string `bson:"city"`
}

func main() {
	session, err := mgo.Dial(url)

	defer session.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	//
	c := session.DB("ich_yaoqi").C("Persons")
	p := &person{
		AGE:    18,
		NAME:   "yaoqi",
		HEIGHT: 100,
		CITY:   "武汉",
	}
	err1 := c.Insert(p)
	if err1 != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("添加", p.NAME, "成功！！")

}

//插入
// func (operater *Operater) insert(p person) error {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	err := collcetion.Insert(p)
// 	return err
// }

// //查询所有
// func (operater *Operater) queryAll() ([]person, error) {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	p := new(person)
// 	p.AGE = 33
// 	query := collcetion.Find(nil)
// 	ps := []person{}
// 	query.All(&ps)
// 	iter := collcetion.Find(nil).Iter()
// 	//
// 	result := new(person)
// 	for iter.Next(&result) {
// 		fmt.Println("一个一个输出：", result)
// 	}
// 	return ps, nil
// }

// //条件查询
// func (operater *Operater) query() ([]person, error) {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	p := new(person)
// 	p.AGE = 33
// 	query := collcetion.Find(bson.M{"age": bson.M{"$eq": 21}})
// 	ps := []person{}
// 	query.All(&ps)
// 	fmt.Println(ps)
// 	return ps, nil
// }

// //更新一行
// func (operater *Operater) update() error {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	update := person{
// 		33,
// 		"詹姆斯",
// 		201,
// 	}
// 	err := collcetion.Update(bson.M{"name": "周杰伦"}, update)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return err
// }

// //更新所有数据
// func (operater *Operater) updateAll() error {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	update := person{
// 		33,
// 		"詹姆斯",
// 		201,
// 	}
// 	changeinfo, err := collcetion.UpdateAll(bson.M{"name": "周杰伦"}, update)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("共有多少行", changeinfo.Matched, "影响")
// 	return nil
// }

// //单行删除
// func (operater *Operater) delete(seletor interface{}) error {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	return collcetion.Remove(seletor)
// }

// //统计文档中数据的个数
// func (operater *Operater) count() (int, error) {
// 	collcetion := operater.mogSession.DB(operater.dbname).C(operater.document)
// 	i, err := collcetion.Count()
// 	return i, err
// }
