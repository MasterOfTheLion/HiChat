package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/middleware"
	"HiChat/models"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func List(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "获取用户列表失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func LoginByNameAndPassWord(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	data, err := dao.FindUserByName(name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "登录失败",
		})
		return
	}

	if data.Name == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名不存在",
		})
		return
	}

	ok := common.CheckPassWord(password, data.Salt, data.PassWord)
	if !ok {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "密码错误",
		})
		return
	}

	Rsp, err := dao.FindUserByNameAndPwd(name, data.PassWord)
	if err != nil {
		zap.S().Info("登陆失败", err)
	}

	token, err := middleware.GenerateToken(Rsp.ID, "yk")
	if err != nil {
		zap.S().Info("生成token失败", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"tokens":  token,
		"userId":  Rsp.ID,
	})
}

func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	repassword := ctx.Request.FormValue("Identity")

	if user.Name == "" || password == "" || repassword == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或密码不能为空",
			"data":    user,
		})
		return
	}

	_, err := dao.FindUser(user.Name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户已注册",
			"data":    user,
		})
		return
	}

	if password != repassword {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}

	salt := fmt.Sprintf("%d", rand.Int31())

	user.PassWord = common.SaltPassWord(password, salt)
	user.Salt = salt
	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t
	dao.CreateUser(user)
	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "新增用户成功",
		"data":    user,
	})
}

func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("类型转换失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注销账号失败",
		})
		return
	}
	user.ID = uint(id)
	Name := ctx.Request.FormValue("name")
	PassWord := ctx.Request.FormValue("password")
	Email := ctx.Request.FormValue("email")
	Phone := ctx.Request.FormValue("phone")
	avatar := ctx.Request.FormValue("icon")
	gender := ctx.Request.FormValue("gender")
	if Name != "" {
		user.Name = Name
	}
	if PassWord != "" {
		salt := fmt.Sprintf("%d", rand.Int31())
		user.Salt = salt
		user.PassWord = common.SaltPassWord(PassWord, salt)
	}
	if Email != "" {
		user.Email = Email
	}
	if Phone != "" {
		user.Phone = Phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		zap.S().Info("参数不匹配", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数不匹配",
		})
		return
	}

	Rsp, err := dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("更新用户失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "修改信息失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "修改成功",
		"data":    Rsp.Name,
	})
}

func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("类型转换失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注销账号失败",
		})
		return
	}

	user.ID = uint(id)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("注销用户失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注销账号失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注销账号失败",
	})
}
