package main

// const MySQLDSN = "root:123456@tcp(127.0.0.1:3306)/betxin?charset=utf8mb4&parseTime=True"

// func connectDB(dsn string) *gorm.DB {
// 	db, err := gorm.Open(mysql.Open(dsn))
// 	if err != nil {
// 		panic(fmt.Errorf("connect db fail: %w", err))
// 	}
// 	return db
// }

// func main() {
// 	db := connectDB(MySQLDSN)
// 	query.SetDefault(db)
// 	ctx := context.Background()
// 	u := &sqlmodel.User{
// 		IdentityNumber: "5678",
// 	}

// 	err := query.User.WithContext(ctx).Debug().Create(u)
// 	if err != nil {
// 		panic(err)
// 	}

// 	time.Sleep(10 * time.Second)
// 	ret, err := query.User.WithContext(ctx).Debug().
// 		Where(query.User.IdentityNumber.Eq("5678")).
// 		Update(query.User.UID, "1233333-123123-21-31-23-12-")

// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(ret)
// }
