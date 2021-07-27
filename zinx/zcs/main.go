package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx/zinx/zcs/pb"
)

func main() {
	// 定义一个person结构对象
	person := &pb.Person{
		Name: "lfxpupa",
		Age: 18,
		Emails: []string{"lfxpupa.com", "1325132780@qq.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "1325132780",
				Type: pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "15718815231",
				Type: pb.PhoneType_MOBILE,
			},
		},
	}
	fmt.Println(person)
	// 编码 得到一个二进制文件
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err", err)
	}
	// data为传输的数据 对端需要按照message Person 格式进行解析
	// 解码
	newData := &pb.Person{}
	proto.Unmarshal(data, newData)
	fmt.Println(newData)
}
