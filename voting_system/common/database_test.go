/* 测试数据库操作的测试函数 */

package common

import (
	"fmt"
	"ginEssential/controller/model"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestInitDB(t *testing.T) {
// 	InitDB()
// 	// t.fl
// }

func TestGetuser(t *testing.T) {
	user, err := GetUser(10)
	fmt.Println("No.10 :", user, err)
	JsonString(user)
	assert.Equal(t, user.ID, 10)

	fmt.Println("")

	user, err = GetUser(100)
	fmt.Println("No.100 :", user, err)
	JsonString(user)
	assert.Equal(t, user.ID, 0)
}

// 测试通过电话查找用户的函数
func TestGetuserByTel(t *testing.T) {
	user, err := GetUserByTel("11111111111")
	fmt.Println("Tel = 11111111111 :", user, err)
	// JsonString(user)
	assert.Equal(t, user.ID, uint(7))

	fmt.Println("")

	user, err = GetUserByTel("1114514")
	fmt.Println("Tel = 1114514 :", user, err)
	// JsonString(user)
	assert.Equal(t, user.ID, uint(0))
}

func TestAddUser(t *testing.T) {
	user, err := AddUser(model.User{Name: "Marisa"})
	if err != nil {
		t.Fatal("Error!", err)
	}
	fmt.Println(user)
}

func TestAddUserList(t *testing.T) {
	// 注意下面不需要再加&model.User了，已经智能到不需要了，有技术的
	// 而且主键似乎不能指定，会自动递增
	// 批量创建在旧版本不支持，更新后支持了
	users := []*model.User{
		{Name: "Sakuya"},
		{Name: "Youmu"},
	}
	n, err := AddUserList(users)
	if err != nil {
		t.Fatal("Error!", err)
	}
	fmt.Println("Inserted items:", n)
}

func TestOmit(t *testing.T) {
	// 除了指定项以外，其余项与user相同
	user := model.User{Name: "Sanae", Age: 16}
	result := DB.Omit("Age").Create(&user)
	fmt.Println("Error =", result.Error)
}

/* 测试函数的模板 */
func TestDemo(t *testing.T) {
	// 对一个结构连续运行first，last，得到的都是首条记录
	user := model.User{}
	DB.First(&user)
	if DB.Error != nil {
		t.Fatal(DB.Error)
	}
}

// 测试各种orm的方法
func TestDB(t *testing.T) {
	// 对一个结构连续运行first，last，得到的都是首条记录
	// 即使用一个数组去接受用数组查询的结果，也只会修改数组的第一个值
	user := model.User{}

	// mysql不区分大小写，所以似乎属性名不必和表中大小写一致
	// 但是如果是定义结构体来查询，因为属于golang的范围所以还是区分大小写
	// DB.Where(map[string]interface{}{"name": "东风谷早苗"}).First(&user)
	// DB.Where(&model.User{2, "东风谷早苗", 18}).First(&user)
	// DB.Where(&model.User{Name: "东风谷早苗"}).First(&user)

	// Plain SQL
	// DB.Find(&user, "name = ?", "泄矢诹访子")
	DB.Where("name = ?", "泄矢诹访子").First(&user)
	// SELECT * FROM users WHERE name = "jinzhu";
	// DB.Where()

	// SELECT * FROM users;
	// DB.Find(&user)

	// SELECT * FROM users WHERE id = 3;
	// DB.First(&user, 3)

	// 下面的语句将查询编号1，1,4,，然后获取最后的编号3作为结果
	// users := [3]*model.User{}
	// DB.Take(&users, []int{1, 5, 4})
	// DB.Last(&user)
	// DB.Take(&user)
	// fmt.Println(users[0])
	// fmt.Println(user)

	// if DB.Error != nil {
	// 	t.Fatal(DB.Error)
	// }
	fmt.Println(user)
	// fmt.Println(DB.RowsAffected)
	// fmt.Println(DB.Error)
}

func TestSetTable(t *testing.T) {
	// js := json.Marshaler(DB_character)

	// JsonString(DB)
	// JsonString(DB_user)
	// JsonString(DB)

	fmt.Println(DB)
	fmt.Println(DB_user)
	music := model.Music{}
	DB.First(&music, 1)
	JsonString(music)
	music.ID = 0
	DB.First(&music, 4)
	JsonString(music)

	// DB_character.AutoMigrate(model.Character{})
	// c := model.Character{Name: "reimu", Poll: 0}
	// DB_character.Create(&c)

	// DB_character.
	// DB_vote.AutoMigrate(model.Vote{})
	// DB_character.First()

	// 这样就可以创建表了
	DB_vote.AutoMigrate(model.Vote{})
}

// 测试读取文件的每一行并处理打印
func TestReadLinesFromFile(t *testing.T) {
	// 替换为你的文件路径
	filename := "works_list.txt"
	// filename := "res_final.txt"

	lines, err := ReadLinesFromFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Lines from the file:")
	for _, line := range lines {
		// fmt.Println(line)
		sep := strings.Split(line, " ")
		fmt.Println(sep[0], sep)
	}
}

// 插入所有作品的记录
func TestInsertWorks(t *testing.T) {
	lines, err := ReadLinesFromFile("works_list.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println(DB_work)

	DB_work.AutoMigrate(model.Work{})
	for _, line := range lines {
		// fmt.Println(line)
		sep := strings.Split(line, " ")
		w := model.Work{Name: sep[2], Poll: 0}
		fmt.Println(sep[2], sep)
		DB_work.Create(&w)
	}
}

// 插入所有音乐的记录
func TestInsertMusic(t *testing.T) {
	lines, err := ReadLinesFromFile("music_list.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	DB.AutoMigrate(model.Music{})
	for _, line := range lines {
		// fmt.Println(line)
		sep := strings.Split(line, "\t")
		id, err := strconv.Atoi(strings.TrimSpace(sep[0]))
		if err != nil {
			fmt.Println("Error when Atoi:", err)
			return
		}
		// fmt.Println(sep[1], id, sep)

		w := model.Music{Name: sep[1], Poll: 0, Work_id: id}
		DB.Create(&w)
	}
}

// 测试gorm语句的执行
func TestGetPartIndex(t *testing.T) {
	// 获取id为10的作品信息
	work := model.Work{}
	DB_work.First(&work, 10)
	JsonString(work)
}

// 获取用户和音乐，并让用户为音乐投票
func TestVote(t *testing.T) {
	u := model.User{}
	DB_user.First(&u, 2)
	m := model.Music{}
	DB.First(&m, 100)
	JsonString(u)
	JsonString(m)
	DB_vote.Create(&model.Vote{
		User_id: int(u.ID),
		Item_id: int(m.ID),
		Part:    1},
	)
}

func TestIsTelephoneExist(t *testing.T) {
	res := IsTelephoneExist("114514")
	fmt.Println(res)

	res = IsTelephoneExist("114514")
	fmt.Println(res)
}

func TestIsUserExist(t *testing.T) {
	res := IsUserExist(111)
	fmt.Println(res)
}

// 用户多次投票的测试
func TestUserVoteForMusic(t *testing.T) {
	UserVoteForMusic(1, 54)
	// UserVoteForMusic(2, 5)
	UserVoteForMusic(5, 8)
	// CancelUserVoteForMusic(1, 3)
}

// 获取音乐排行榜
func TestGetMusicRank(t *testing.T) {
	ms, _ := GetMusciRank(20)
	for _, v := range ms {
		fmt.Println(v.ID, v.Name, v.Poll)
		// JsonString(v)
	}
}

func TestRedis(t *testing.T) {

}

// 测试获取卡牌信息
// func TestCard(t *testing.T) {
// 	card := model.Card{}
// 	DB.First(&card, 3)
// 	if DB.Error != nil {
// 		t.Fatal(DB.Error)
// 	}
// 	fmt.Println(card)
// }
