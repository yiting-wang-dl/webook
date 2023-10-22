package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/webook/internal/domain"
	"github.com/webook/internal/service"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) { // when to use *gin.Context and when use context.Context? Why it doesn't return anything?
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "Incorrect Email Format")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "Inconsistent Password")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "Password must contain alphabets, numbers, special characters, and minimum 8 characters.")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "Signup Complete")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "Email Already Exists")
	default:
		ctx.String(http.StatusOK, "System Error")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 600,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "System Error")
			return
		}
		ctx.String(http.StatusOK, "Login Successful")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "Invalid Username or Password")
	default:
		ctx.String(http.StatusOK, "System Error")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		h.setJWTToken(ctx, u.Id)
		ctx.String(http.StatusOK, "Login Successful")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "Incorrect Username or Password")
	default:
		ctx.String(http.StatusOK, "System Error")
	}
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // 15 minutes
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 15)), // testing in postman
		},
	}
	fmt.Println("initiate time: ", time.Now())
	fmt.Println("ExpiresAt: ", time.Now().Add(time.Hour*15))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "System Error")
	}
	ctx.Header("x-jwt-token", tokenStr)
	println("tokenStr", tokenStr)
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		//Id       int64  `json:"uid"`
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"` // YYYY-MM-DD
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 1. session id
	//sess := sessions.Default(ctx)
	//uid, ok := sess.Get("userId").(int64)
	//if !ok {
	//	ctx.String(http.StatusOK, "Cannot convert uid into int64")
	//	return
	//}
	// or 2. JWT
	uc, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "Incorrect Birthday Format")
		return
	}
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		//Id:       uid,  // sessionid
		Id:       uc.Uid, // JWT
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	ctx.String(http.StatusOK, "Your profile is Updated")
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	uid, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "Cannot convert uid into int64")
		return
	}
	user, err := h.svc.FindById(ctx, uid)
	if err != nil {
		ctx.String(http.StatusOK, "Fail to retrieve user information")
	}
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		Email:    user.Email,
		Nickname: user.Nickname,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday.Format(time.DateOnly),
	})
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
