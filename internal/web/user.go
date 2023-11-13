package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/webook/internal/domain"
	"github.com/webook/internal/service"
	ijwt "github.com/webook/internal/web/jwt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService,
	//hdl ijwt.Handler,
	codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
		//Handler:        hdl,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/logout", h.LogoutJWT)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
	ug.GET("/refresh_token", h.RefreshToken)

	// SMS validation
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
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
		//err = h.setJWTToken(ctx, u.Id)
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			ctx.String(http.StatusOK, "System Error")
			return
		}
		ctx.String(http.StatusOK, "Login Successful")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "Incorrect Username or Password")
	default:
		ctx.String(http.StatusOK, "System Error")
	}
}

//func (h *UserHandler) Logout(ctx *gin.Context) {
//	sess := sessions.Default(ctx)
//	sess.Options(sessions.Options{
//		MaxAge: -1,
//	})
//	sess.Save()
//}

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
	//uc, ok := ctx.MustGet("user").(UserClaims)
	uc, ok := ctx.MustGet("user").(ijwt.UserClaims)
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

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System Error",
			//Msg: err.Error(),
		})
		zap.L().Error("SMS Validation Failed",
			// Do not log this in prod env
			// using it to test in dev env
			//zap.String("phone", req.Phone),
			zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Code is incorrect, please enter again",
		})
		return
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System Error",
		})
		return
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "Login Successful",
	})

}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// validate Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Please Enter Phone Number",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "Code Sent",
		})
	case service.ErrCodeSendTooMany:
		// if lots of this warning, someone is trying to break in
		zap.L().Warn("Validation Code Sent Too Frequently")
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Code Sent Too Frequent, Please Try Again In A Later Time",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System Error",
		})
		// log
		log.Println(err)
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	//sess := sessions.Default(ctx)
	//uid, ok := sess.Get("userId").(int64)
	//if !ok {
	//	ctx.String(http.StatusOK, "Cannot convert uid into int64")
	//	return
	//}
	uc, ok := ctx.MustGet("user").(ijwt.UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := h.svc.FindById(ctx, uc.Uid)
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
		Nickname: user.Nickname,
		Email:    user.Email,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday.Format(time.DateOnly),
	})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// FE will add refresh-token in Authorization
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// invalid token, or redis issue
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "System Error"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "Log Out Successfully"})
}

// before implementing jwt under web
//func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
//	uc := UserClaims{
//		Uid:       uid,
//		UserAgent: ctx.GetHeader("User-Agent"),
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // 15 minutes
//			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 15)), // testing in postman
//		},
//	}
//	fmt.Println("initiate time: ", time.Now())
//	fmt.Println("ExpiresAt: ", time.Now().Add(time.Hour*15))
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
//	tokenStr, err := token.SignedString(JWTKey)
//	if err != nil {
//		ctx.String(http.StatusOK, "System Error")
//	}
//	ctx.Header("x-jwt-token", tokenStr)
//	println("tokenStr", tokenStr)
//}

//var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")
//
//type UserClaims struct {
//	jwt.RegisteredClaims
//	Uid       int64
//	UserAgent string
//}
