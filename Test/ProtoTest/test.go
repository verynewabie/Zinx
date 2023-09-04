package main

import (
	pb "Zinx/Test/ProtoTest/Proto" // 导入特定目录下的特定包,不加pb则全部导入
	"fmt"
	"google.golang.org/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "Aceld",
		Age:    16,
		Emails: []string{"https://legacy.gitbook.com/@aceld", "https://github.com/aceld"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "13113111311",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "14141444144",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "19191919191",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	newData := &pb.Person{}
	err = proto.Unmarshal(data, newData)
	if err != nil {
		fmt.Println("unmarshal err:", err)
	}

	fmt.Println(newData)
}
