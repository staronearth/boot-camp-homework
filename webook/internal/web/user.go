package web

import (
	"boot-camp-homework/webook/internal/domain"
	"boot-camp-homework/webook/internal/service"
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

const biz = "login"

// 确保 UserHandler 上实现了 handler 接口
var _ handler = &UserHandler{}

// 这个给更优雅
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	jwth        Token
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, jwth Token) *UserHandler {
	const (
		emailRegexPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		//设置密码长度不超过72字节
		passwordPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,72}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		jwth:        jwth,
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
	// PUT "/login/sms/code" 发验证码
	// POST "/login/sms/code" 校验验证码
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
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
		return
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
	if errors.Is(err, service.ErrUserDuplicate) {
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
	domainuser, err := u.svc.Login(context.Request.Context(), domain.User{
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

	if err = u.jwth.SetJWTToken(context, domainuser.Id); err != nil {
		context.String(http.StatusInternalServerError, "系统错误")
		return
	}
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

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 是不是一个合法的手机号
	// 考虑正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	err := u.codeSvc.Send(ctx.Request.Context(), biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "发送太频繁,请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//这边，可以加上各种校验
	ok, err := u.codeSvc.Verify(ctx.Request.Context(), biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(ctx.Request.Context(), req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//这里要写登录成功的逻辑
	err = u.jwth.SetJWTToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证码校验通过",
	})
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
