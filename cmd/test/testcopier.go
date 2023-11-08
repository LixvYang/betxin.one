package main

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/sqlmodel"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

func main() {
	betxinUser := schema.User{
		IdentityNumber: "123123",
		UID:            "123123",
		IsMvmUser:      true,
	}
	sqlU := sqlmodel.User{}
	copier.Copy(&sqlU, &betxinUser)
	fmt.Printf("%#v", sqlU)
}
