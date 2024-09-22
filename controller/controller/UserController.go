/*
包括用户注册，登录，注销，查看用户信息等操作
*/
package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/controller/model"
	"ginEssential/controller/util"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/* POST请求注册时，处理的函数 */
func Register(c *gin.Context) {
	// 获取参数
	var err error
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	fmt.Println("[Mid] Register:", name, telephone, password)

	// 各项数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}
	if len(name) == 0 { // 如果用户没有给昵称，则随机一个
		name = util.RandomString(10)
	}

	// 检查数据库中是否存在该手机号
	if common.IsTelephoneExist(telephone) {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 422, "msg": "手机号已被注册"})
		return
	}

	/* 创建用户 */
	// 给密码加密后再保存
	password, err = util.EncryptedString(password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 500, "msg": "加密出错"})
		return
	}

	// 定义用户结构并插入数据库
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}

	// 向数据库中添加用户
	newID, err := common.AddUser(newUser)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 500, "msg": "添加用户操作出错:" + err.Error()})
		return
	}

	// 设置cookie
	fmt.Println(newID)
	// 返回结果
	c.String(200, name+"注册成功!")
}

/* 用户登录处理函数 */
func Login(c *gin.Context) {
	// 获取各项参数
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	fmt.Println("Login:", name, telephone, password)

	// 数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	// 根据手机号获取用户信息，判断用户是否存在
	user, _ := common.GetUserByTel(telephone)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity,
			gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	// 判断密码是否正确
	err := util.DecryptString(user.Password, password)
	if err != nil { //密码错误
		c.JSON(400, gin.H{"msg": "密码错误"})
		return
	}

	// 设置session，保证id唯一即可
	session := sessions.Default(c)
	session.Set("login", user.ID) //
	// session.Set("login"+strconv.Itoa(int(user.ID)), user.ID) //
	session.Save()

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  strconv.Itoa(int(user.ID)) + "用户登录成功！",
	})
}

// 设置并发送用户的cookie
func AddCookie(c *gin.Context) {
	id := c.MustGet("user_id").(string)
	fmt.Println(id)
	c.SetCookie("ID_cookie", id, 3600, "/", "localhost", false, true)
	c.Next()
	// c.String(200, "We had send you a coookie.")
}

// 获取用户的头部中的cookie，并设置user_ID变量
func GetCookie(c *gin.Context) {
	id, err := c.Cookie("ID_cookie")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "No cookie!"})
		return
	}

	fmt.Println("Get ID from cookie", id)
	c.Set("user_ID", id)
	c.Next()
	// c.String(200, "You have cookie!")
}

// 根据user_ID，到缓存和数据库中判断是否是有效用户,如果有效则将其信息保存到content中
func CheckSeccison(c *gin.Context) {
	s, exist := c.Get("user_ID")
	fmt.Println(s, exist)
	fmt.Printf("%t\n", s)
	if !exist {
		c.AbortWithStatusJSON(403, gin.H{"msg": "无用户记录！"})
	}

	// 类型转换
	id, err := strconv.Atoi(s.(string))
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"msg": "字符串ID转换整数失败"})
	}

	var user model.User
	// TODO：查找是否在缓存中存在

	// 未查找到相应用户数据
	user, _ = common.GetUser(id)
	if user.ID == 0 {
		c.AbortWithStatusJSON(403, gin.H{"msg": "无该用户的记录！"})
		return
	}

	c.Set("user", user)
	c.Next()
	c.String(200, "You cookie is valid in session!")
}

// 从之前的中间件中接受用户数据user,即一个User结构对象
func GetUser(c *gin.Context) {
	var user model.User
	u, exist := c.Get("user")

	// 如果之前的中间件未设置信息
	if !exist {
		c.AbortWithStatusJSON(200, gin.H{"msg": "无该用户的记录！"})
		return
	}

	// 类型断言，将any转变为user
	if ur, ok := u.(model.User); !ok {
		c.String(400, "用户信息转换失败。")
		return
	} else {
		user = ur
	}

	fmt.Println(user)
	c.Set("Get user:", user)
	c.Next()
	c.String(200, "You cookie is valid!")
}
