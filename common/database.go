/* 数据库的相关操作，即初始化，以及增删改查 */

package common

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	// "ginEssential/model"
	"errors"
	"ginEssential/controller/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB // 数据库
	DB_character *gorm.DB // 角色数据表
	DB_user      *gorm.DB // 用户数据表
	DB_work      *gorm.DB // 作品数据表
	DB_vote      *gorm.DB // 投票数据表，关联用户和角色
)

// const DataBase = "javabook"
// const Table = "mg"
const DataBase = "gin"
const Table = "mooc_user"

/* 连接数据库 */
func InitDB() *gorm.DB {
	// 注意我的mysql的端口是3306，gin的端口是8080，
	// 连接sql用的是mysql的端口即3306
	// driverName := "mysql"
	host := "localhost"
	port := "3306"
	username := "root"
	database := DataBase
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		username,
		password,
		host,
		port,
		database,
		charset,
	)

	db, err := gorm.Open(mysql.Open(args))
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	// db.AutoMigrate(&model.User{})
	return db
}

/* 初始化数据库对象 */
func init() {
	DB = InitDB()
	fmt.Println("~~~Value DB was inited.")
}

/* 获取数据库对象，单例模式 */
func GetDB() *gorm.DB {
	if DB == nil {
		DB = InitDB()
		fmt.Println("~~~Value DB was inited.")
	}
	return DB
}

/* 设置当前数据表 */
func SetTable(name string) *gorm.DB {
	// fmt.Println("Current Table is", name)
	return DB.Table(name)
}

/* 查询数据库，检查某个电话号码是否已经存在，返回该布尔值 */
func IsTelephoneExist(tel string) bool {
	var user model.User
	res := DB.Where("telephone = ?", tel).Limit(1).Find(&user)
	// res := DB.Where("telephone = ?", tel).First(&user)
	// err := res.Error
	// fmt.Println(err)
	// return !errors.Is(err, gorm.ErrRecordNotFound)
	return res.RowsAffected != 0
}

/* 添加一个用户，返回该用户的id以及错误 */
func AddUser(user model.User) (uint, error) {
	res := DB.Create(&user)
	err := res.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return user.ID, err
}

/* 查询数据库，检查某个用户是否已经存在，返回相应布尔值 */
func IsUserExist(id int) bool {
	var user model.User
	res := DB.Limit(1).Find(&user, id)
	return res.RowsAffected != 0
}

/*
查询某一ID的用户信息，返回相应结果用户和错误
如果数据库中没有，则ID=0，但是err为nil
*/
func GetUser(id int) (model.User, error) {
	var user model.User
	res := DB.Limit(1).Find(&user, id)
	return user, res.Error
}

/*
查询某一电话号码的用户信息，返回相应结果用户和错误
如果数据库中没有，则ID=0，但是err为nil
*/
func GetUserByTel(tel string) (model.User, error) {
	var user model.User
	res := DB.Limit(1).Where("telephone = ?", tel).Find(&user)
	return user, res.Error
}

/* 添加一系列用户 */
func AddUserList(users []*model.User) (int, error) {
	db := DB.Create(users)
	err := db.Error
	return len(users), err
}

// 某用户为某音乐投票
func UserVoteForMusic(uid int, mid int) (err error) {
	// 判断用户是否存在
	fmt.Println(uid, "vote for", mid)
	var user model.User
	resUser := DB.First(&user, uid)
	if resUser.Error != nil {
		return resUser.Error
	}
	if user.ID == 0 {
		return errors.New("no this user")
	}
	if user.Poll == 0 { // 没有票了
		return errors.New("no poll can be used")
	}

	// 判断音乐是否存在
	var music model.Music
	resMusic := DB.First(&music, mid)
	if resMusic.Error != nil {
		return resMusic.Error
	}
	if music.ID == 0 {
		return errors.New("no this music")
	}

	// 还要判断是否是重复投票（

	// 先搜索音乐所在的部分编号，然后插入相应记录
	res := DB.Create(&model.Vote{User_id: uid, Item_id: mid, Part: 1}) // music属于第一个部门
	err = res.Error
	if err != nil {
		return err
	}

	// 然后将对应音乐的票数加一,将对应用户的可票数减一
	DB.Model(&music).Update("poll", music.Poll+1)
	DB.Model(&user).Update("poll", user.Poll-1)

	// 更新缓存的排名
	target := strconv.Itoa(int(music.ID))
	index, err := REDIS_DB.ZRank(CTX, "rank", target).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	fmt.Println(index)
	REDIS_DB.ZIncrBy(CTX, "rank", 1, target)
	return err
}

// 某用户取消为某音乐的投票
func CancelUserVoteForMusic(uid int, mid int) error {
	// 判断用户是否存在
	var user model.User
	res := DB.First(&user, uid)
	if res.Error != nil {
		return res.Error
	}
	if user.ID == 0 {
		return errors.New("no this user")
	}

	// 判断音乐是否存在
	var music model.Music
	res = DB.First(&music, mid)
	if res.Error != nil {
		return res.Error
	}
	if music.ID == 0 {
		return errors.New("no this musci")
	}

	// 判断该投票记录是否存在

	// 先搜索相应记录的编号，然后删除相应记录
	var vote model.Vote
	// vote = model.Vote{
	// 	User_id: uid,
	// 	Item_id: mid,
	// 	Part:    1,
	// }
	// res = DB.Delete(&vote)
	res = DB.Limit(1).Where("user_id = ? and item_id = ? and part =?", uid, mid, 1).Delete(&vote)
	err := res.Error
	fmt.Println(err)
	if err != nil {
		return err
	}

	// 将对应用户的可票数加一,将对应音乐的票数减一
	DB.Model(&music).Update("poll", music.Poll-1)
	DB.Model(&user).Update("poll", user.Poll+1)

	// 更新缓存的排名
	target := strconv.Itoa(int(music.ID))
	index, err := REDIS_DB.ZRank(CTX, "rank", target).Result()
	if err == redis.Nil { // 不存在热榜缓存中
		return nil
	} else if err != nil {
		return err
	}

	fmt.Println(index)
	REDIS_DB.ZIncrBy(CTX, "rank", -1, target)
	return err
}

// 获取音乐排行榜的前num项，并缓存
func GetMusciRank(num int) ([]model.Music, error) {
	// 先在缓存中查找
	rres, err := REDIS_DB.ZRevRange(CTX, "rank", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var musics []model.Music
	number := len(rres)
	if number > 0 { // 存在缓存
		fmt.Println("缓存命中！")
		// 获取榜上的所有id
		ids := make([]int, number)
		for i, v := range rres {
			ids[i], _ = strconv.Atoi(v)
		}

		err = DB.Find(&musics, ids).Error
		if err != nil {
			return nil, err
		}
		sort.Sort(model.Musics(musics))
		return musics, nil
	}

	// 再去SQL中查找
	res := DB.Limit(num).Order("poll desc").Find(&musics)
	if res.Error != nil {
		return nil, err
	}

	// 存入缓存，并设置过期时间
	for _, v := range musics {
		fmt.Println(v)
		REDIS_DB.ZAdd(CTX, "rank", MusicScore(&v))
	}
	REDIS_DB.Expire(CTX, "rank", time.Hour)
	return musics, res.Error
}
