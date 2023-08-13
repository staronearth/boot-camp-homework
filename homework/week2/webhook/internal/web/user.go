package web

import (
	"boot-camp-homework/homework/week2/webhook/internal/domain"
	"boot-camp-homework/homework/week2/webhook/internal/service"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"unicode/utf8"
)

var (
	ErrBirthdayFormat    = errors.New("日期格式不对:(1999-09-09)")
	ErrBirthdayIncorrect = errors.New("超出日期限制:(1880-00-00~2023-08-13)")
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		//设置密码长度不超过72字节
		passwordPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,72}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.Signup)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	type ProfileRes struct {
		NickName        string `json:"nick_name"`
		BrithDays       string `json:"brith_days"`
		PersonalProfile string `json:"personal_profile"`
	}
	cliams, _ := ctx.Get("cliams")
	//sess := sessions.Default(ctx)
	//id := sess.Get("UserId")
	usercliams, ok := cliams.(*UserCliams)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
	}

	uinfo, err := u.svc.Profile(ctx.Request.Context(), usercliams.Id)
	if err != nil {
		ctx.String(http.StatusOK, "修改信息失败")
		return
	}
	res := ProfileRes{
		NickName:        uinfo.NickName,
		BrithDays:       uinfo.BrithDays,
		PersonalProfile: uinfo.PersonalProfile,
	}
	ctx.JSON(http.StatusOK, res)
	//context.String(http.StatusOK, "这是profile")
	return
}

func (u *UserHandler) Signup(context *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	if err := context.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		context.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		context.String(http.StatusOK, "两次输入的密码不一致")
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		context.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}
	// 调用一下svc的方法
	err = u.svc.SignUp(context, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		context.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	context.String(http.StatusOK, "注册成功")
	//log.Println(req)
}

func (u *UserHandler) Login(context *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := context.Bind(&req); err != nil {
		return
	}
	domainuser, err := u.svc.Login(context, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		context.String(http.StatusOK, "用户名和密码不对")
		return
	}
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	// 在这里登录成功
	// 设置session
	sess := sessions.Default(context)
	// 我可以随便设置值
	// 你要放在session里面的值
	sess.Set("UserId", domainuser.Id)
	sess.Options(sessions.Options{
		//Secure:   false,
		//HttpOnly: false,
		MaxAge: 180,
	})
	sess.Save()
	context.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) LoginJWT(context *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := context.Bind(&req); err != nil {
		return
	}
	domainuser, err := u.svc.Login(context, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		context.String(http.StatusOK, "用户名和密码不对")
		return
	}
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	//步骤2
	//在这里用JWT设置登录态
	//生成一个JWT token
	//

	key := []byte("4B2B17F68975BA8C806846C5CC898070F3F62423BEC9D87F2DBE844B4C14F137")
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &UserCliams{
		Id: domainuser.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserAgent: context.Request.UserAgent(),
	})
	tokenStr, err := token.SignedString(key)
	if err != nil {
		context.String(http.StatusInternalServerError, "系统错误")
		return
	}
	context.Header("x-jwt-token", tokenStr)
	context.String(http.StatusOK, "登录成功")
	return
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName        string `json:"nick_name"`
		BrithDays       string `json:"brith_days"`
		PersonalProfile string `json:"personal_profile"`
	}
	cliams, _ := ctx.Get("cliams")
	//sess := sessions.Default(ctx)
	//id := sess.Get("UserId")
	usercliams, ok := cliams.(*UserCliams)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if utf8.RuneCountInString(req.NickName) == 0 || utf8.RuneCountInString(req.NickName) > 100 {
		ctx.String(http.StatusOK, "你的昵称不符合长度")
		return
	}
	if utf8.RuneCountInString(req.PersonalProfile) > 1000 {
		ctx.String(http.StatusOK, "你的个人简介超过限制")
		return
	}
	err := boolBirthday(req.BrithDays)
	if err != nil {
		if errors.Is(err, ErrBirthdayFormat) {
			ctx.String(http.StatusOK, "日期格式不对:(1999-09-09)")
			return
		}
		if errors.Is(err, ErrBirthdayIncorrect) {
			ctx.String(http.StatusOK, "超出日期限制:(1880-00-00~2023-08-13)")
			return
		}
	}

	err = u.svc.Edit(ctx.Request.Context(), domain.UserInfo{
		Id:              usercliams.Id,
		NickName:        req.NickName,
		BrithDays:       req.BrithDays,
		PersonalProfile: req.PersonalProfile,
	})
	if err != nil {
		ctx.String(http.StatusOK, "修改信息失败")
		return
	}
	ctx.String(http.StatusOK, "修改信息成功")
}
func boolBirthday(birthday string) error {
	layout := "2006-01-02"
	birthtime, err := time.Parse(layout, birthday)
	if err != nil {
		return ErrBirthdayFormat
	}
	now := time.Now()
	if timesub := now.Sub(birthtime).Seconds(); timesub < 0 || timesub > 150*3600*24*365 {
		return ErrBirthdayIncorrect
	}
	return nil
}

type UserCliams struct {
	Id int64
	jwt.RegisteredClaims
	UserAgent string
}
