package gfmiddleware

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
)

// MapClaims 类型
//  使用map[string]interface{}进行JSON解码
//  如果不提供声明（claims）类型，则这是默认的声明（claims）类型
type MapClaims map[string]interface{}

// GfJWTMiddleware 提供Json Web令牌身份验证实现。
//  失败时，将返回401http响应。
//  成功时，将调用包装好的中间件，并将userID作为c.Get("userID").(string)提供。
//  用户可以通过向LoginHandler发送json请求来获取令牌。
//  然后需要在身份验证头（Authentication header）中传递令牌。
//  示例：Authorization:Bearer XXX_TOKEN_XXX
type GfJWTMiddleware struct {
	// Realm 要显示给用户的名称。
	//  必须的。
	Realm string

	// SigningAlgorithm 签名算法
	//  可能的值为 HS256，HS384，HS512
	//  可选，默认值为HS256。
	SigningAlgorithm string

	// Key 用于签名的密钥。
	//  必需的。
	Key []byte

	// Timeout jwt令牌有效的持续时间。
	//  可选，默认为1小时。
	Timeout time.Duration

	// MaxRefresh 此字段允许客户端刷新其令牌，直到MaxRefresh过期。
	//  请注意，客户端可以在MaxRefresh的最后一刻刷新其令牌。
	//  这意味着令牌的最大有效时间跨度是 TokenTime + MaxRefresh。
	//  可选，默认为0表示不可刷新。
	MaxRefresh time.Duration

	// Authenticator 根据登录信息执行用户身份验证的回调函数。
	//  必须返回用户数据作为用户标识符，它将存储在声明（Claim）数组中。
	//  必需的。
	//  检查错误 error（e）以确定适当的错误信息。
	Authenticator func(r *ghttp.Request) (interface{}, error)

	// Authorizator 应该执行已验证用户授权的回调函数。
	//  仅在身份验证成功后调用。
	//  成功时必须返回true，失败时必须返回false。
	//  可选，默认为"成功"。
	Authorizator func(data interface{}, r *ghttp.Request) bool

	// PayloadFunc 将在登录期间调用的回调函数。
	//  使用此函数可以向web token添加额外的有效载荷数据。
	//  然后在请求期间通过c.Get("JWT_PAYLOAD")获取数据。
	//  注意，有效负载没有加密。jwt.io上提到的属性不能用作map的键。
	//  可选，默认情况下不会设置其他数据。
	PayloadFunc func(data interface{}) MapClaims

	// Unauthorized 用户可以定义自己的Unauthorized函数。
	Unauthorized func(*ghttp.Request, int, string)

	// LoginResponse 用户可以定义自己的LoginResponse函数。
	LoginResponse func(*ghttp.Request, int, string, time.Time)

	// RefreshResponse 用户可以定义自己的RefreshResponse函数。
	RefreshResponse func(*ghttp.Request, int, string, time.Time)

	// LogoutResponse 用户可以定义自己的LogoutResponse函数。
	LogoutResponse func(*ghttp.Request, int)

	// IdentityHandler 设置身份处理函数
	IdentityHandler func(*ghttp.Request) interface{}

	// IdentityKey 设置身份密钥
	IdentityKey string

	// TokenLookup 是"<source>:<name>"形式的字符串，用于从请求中提取token。
	//  可选。默认值 "header:Authorization"。
	//  可能值：
	//  - "header:<name>"
	//  - "query:<name>"
	//  - "cookie:<name>"
	TokenLookup string

	// TokenHeadName 是请求头中的字符串。
	//  默认值是"Bearer"
	TokenHeadName string

	// TimeFunc 提供当前时间。
	//  可以重写它以使用另一个时间值。
	//  如果您的服务器使用的时区与您的令牌不同，这对于测试非常有用。
	TimeFunc func() time.Time

	//  HTTPStatusMessageFunc 当JWT中间件中的某些东西发生故障时，使用HTTP状态消息。
	//   检查错误 error（e）以确定适当的错误信息。
	HTTPStatusMessageFunc func(e error, r *ghttp.Request) string

	// PrivKeyFile 非对称算法的私钥文件
	PrivKeyFile string

	// PubKeyFile 非对称算法的公钥文件
	PubKeyFile string

	// privKey 私钥
	privKey *rsa.PrivateKey

	// pubKey 公钥
	pubKey *rsa.PublicKey

	// SendCookie 可以选择将令牌作为cookie返回
	SendCookie bool

	// SecureCookie 允许不安全的Cookie通过HTTP进行开发
	SecureCookie bool

	// CookieHTTPOnly 允许访问客户端的Cookie以进行开发
	CookieHTTPOnly bool

	// CookieDomain 允许更改Cookie域以进行开发
	CookieDomain string

	// SendAuthorization 允许为每个请求返回授权头
	SendAuthorization bool

	// DisabledAbort 禁用上下文的abort()。
	DisabledAbort bool

	// CookieName 允许更改Cookie名称以进行开发
	CookieName string
}

var (
	// IdentityKey 默认身份密钥
	IdentityKey = "identity"
	// blacklist 黑名单，存储尚未过期但已被停用的令牌。
	blacklist = gcache.New()
)

// New GfJWTMiddleware检查错误的方法
func New(m *GfJWTMiddleware) (*GfJWTMiddleware, error) {
	if err := m.MiddlewareInit(); err != nil {
		return nil, err
	}

	return m, nil
}

func (mw *GfJWTMiddleware) readKeys() error {
	err := mw.privateKey()
	if err != nil {
		return err
	}
	err = mw.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (mw *GfJWTMiddleware) privateKey() error {
	keyData, err := ioutil.ReadFile(mw.PrivKeyFile)
	if err != nil {
		return ErrNoPrivKeyFile
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPrivKey
	}
	mw.privKey = key
	return nil
}

func (mw *GfJWTMiddleware) publicKey() error {
	keyData, err := ioutil.ReadFile(mw.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	mw.pubKey = key
	return nil
}

func (mw *GfJWTMiddleware) usingPublicKeyAlgo() bool {
	switch mw.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// MiddlewareInit 初始化jwt配置。
func (mw *GfJWTMiddleware) MiddlewareInit() error {

	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(data interface{}, r *ghttp.Request) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(r *ghttp.Request, code int, message string) {
			r.Response.WriteJson(g.Map{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.LoginResponse == nil {
		mw.LoginResponse = func(r *ghttp.Request, code int, token string, expire time.Time) {
			r.Response.WriteJson(g.Map{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.RefreshResponse == nil {
		mw.RefreshResponse = func(r *ghttp.Request, code int, token string, expire time.Time) {
			r.Response.WriteJson(g.Map{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.LogoutResponse == nil {
		mw.LogoutResponse = func(r *ghttp.Request, code int) {
			r.Response.WriteJson(g.Map{
				"code":    http.StatusOK,
				"message": "success",
			})
		}
	}

	if mw.IdentityKey == "" {
		mw.IdentityKey = IdentityKey
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(r *ghttp.Request) interface{} {
			claims := ExtractClaims(r)
			return claims[mw.IdentityKey]
		}
	}

	if mw.HTTPStatusMessageFunc == nil {
		mw.HTTPStatusMessageFunc = func(e error, r *ghttp.Request) string {
			return e.Error()
		}
	}

	if mw.Realm == "" {
		mw.Realm = "gf jwt"
	}

	if mw.CookieName == "" {
		mw.CookieName = "jwt"
	}

	if mw.usingPublicKeyAlgo() {
		return mw.readKeys()
	}

	if mw.Key == nil {
		return ErrMissingSecretKey
	}
	return nil
}

// MiddlewareFunc 使GfJWTMiddleware实现中间件接口。
func (mw *GfJWTMiddleware) MiddlewareFunc() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		mw.middlewareImpl(r)
	}
}

func (mw *GfJWTMiddleware) middlewareImpl(r *ghttp.Request) {
	claims, token, err := mw.GetClaimsFromJWT(r)
	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, r))
		return
	}

	if claims["exp"] == nil {
		mw.unauthorized(r, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrMissingExpField, r))
		return
	}

	if _, ok := claims["exp"].(float64); !ok {
		mw.unauthorized(r, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp, r))
		return
	}

	if int64(claims["exp"].(float64)) < mw.TimeFunc().Unix() {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, r))
		return
	}

	in, err := mw.inBlacklist(token)
	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, r))
		return
	}

	if in {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrInvalidToken, r))
		return
	}

	r.SetParam("JWT_PAYLOAD", claims)
	identity := mw.IdentityHandler(r)

	if identity != nil {
		r.SetParam(mw.IdentityKey, identity)
	}

	if !mw.Authorizator(identity, r) {
		mw.unauthorized(r, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, r))
		return
	}

	//c.Next() todo
}

// GetClaimsFromJWT 从JWT令牌获取声明（claims）
func (mw *GfJWTMiddleware) GetClaimsFromJWT(r *ghttp.Request) (MapClaims, string, error) {
	token, err := mw.ParseToken(r)

	if err != nil {
		return nil, "", err
	}

	if mw.SendAuthorization {
		token := r.GetString("JWT_TOKEN")
		if len(token) > 0 {
			r.Header.Set("Authorization", mw.TokenHeadName+" "+token)
		}
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, token.Raw, nil
}

// LoginHandler 客户端可以使用它来获取jwt令牌。
//  有效载荷（Payload）必须为json格式，{"username": "USERNAME", "password": "PASSWORD"}。
//  响应（Reply）的格式为 {"token": "TOKEN"}。
func (mw *GfJWTMiddleware) LoginHandler(r *ghttp.Request) {
	if mw.Authenticator == nil {
		mw.unauthorized(r, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(ErrMissingAuthenticatorFunc, r))
		return
	}

	data, err := mw.Authenticator(r)

	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, r))
		return
	}

	// 创建令牌
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)

	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, r))
		return
	}

	// 设置cookie
	if mw.SendCookie {
		maxage := int64(expire.Unix() - time.Now().Unix())
		r.Cookie.SetCookie(mw.CookieName, tokenString, mw.CookieDomain, "/", time.Duration(maxage)*time.Second)
	}

	mw.LoginResponse(r, http.StatusOK, tokenString, expire)
}

func (mw *GfJWTMiddleware) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	if mw.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(mw.privKey)
	} else {
		tokenString, err = token.SignedString(mw.Key)
	}
	return tokenString, err
}

// LogoutHandler 可用于注销令牌。
//  令牌在注销时仍然需要校验有效性。
//  注销令牌将未过期的令牌放入黑名单。
func (mw *GfJWTMiddleware) LogoutHandler(r *ghttp.Request) {
	claims, token, err := mw.CheckIfTokenExpire(r)
	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, r))
		return
	}

	err = mw.setBlacklist(token, claims)

	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, r))
		return
	}

	mw.LogoutResponse(r, http.StatusOK)
}

// RefreshHandler 可用于刷新令牌。
//  令牌在刷新时仍然需要校验有效性。
//  应放在使用GfJWTMiddleware的端点下。
//  响应（Reply）的格式为{"token": "TOKEN"}。
func (mw *GfJWTMiddleware) RefreshHandler(r *ghttp.Request) {
	tokenString, expire, err := mw.RefreshToken(r)
	if err != nil {
		mw.unauthorized(r, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, r))
		return
	}

	mw.RefreshResponse(r, http.StatusOK, tokenString, expire)
}

// RefreshToken 刷新令牌并检查令牌是否过期
func (mw *GfJWTMiddleware) RefreshToken(r *ghttp.Request) (string, time.Time, error) {
	claims, token, err := mw.CheckIfTokenExpire(r)
	if err != nil {
		return "", time.Now(), err
	}

	// 创建令牌
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(newToken)

	if err != nil {
		return "", time.Now(), err
	}

	// 设置cookie
	if mw.SendCookie {
		maxage := int64(expire.Unix() - time.Now().Unix())
		r.Cookie.SetCookie(mw.CookieName, tokenString, mw.CookieDomain, "/", time.Duration(maxage)*time.Second)
	}

	// 将旧令牌设置为黑名单
	err = mw.setBlacklist(token, claims)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expire, nil
}

// CheckIfTokenExpire 检查令牌是否过期
func (mw *GfJWTMiddleware) CheckIfTokenExpire(r *ghttp.Request) (jwt.MapClaims, string, error) {
	token, err := mw.ParseToken(r)

	if err != nil {
		// 如果我们收到一个错误，并且错误不是一个ValidationErrorExpired，那么我们希望返回错误。
		// 如果错误只是ValidationErrorExpired，我们想继续，因为如果令牌在MaxRefresh时间内，我们仍然可以刷新它。
		if _, ok := err.(*jwt.MalformedTokenError); ok {
			return nil, "", err
		}
	}

	in, err := mw.inBlacklist(token.Raw)

	if err != nil {
		return nil, "", err
	}

	if in {
		return nil, "", ErrInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		return nil, "", ErrExpiredToken
	}

	return claims, token.Raw, nil
}

// TokenGenerator 客户端可用于获取jwt令牌的方法。
func (mw *GfJWTMiddleware) TokenGenerator(data interface{}) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().UTC().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

func (mw *GfJWTMiddleware) jwtFromHeader(r *ghttp.Request, key string) (string, error) {
	authHeader := r.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (mw *GfJWTMiddleware) jwtFromQuery(r *ghttp.Request, key string) (string, error) {
	token := r.GetString(key)

	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

func (mw *GfJWTMiddleware) jwtFromCookie(r *ghttp.Request, key string) (string, error) {
	cookie := r.Cookie.Get(key)

	if cookie == "" {
		return "", ErrEmptyCookieToken
	}

	return cookie, nil
}

func (mw *GfJWTMiddleware) jwtFromParam(r *ghttp.Request, key string) (string, error) {
	token := r.GetString(key)
	if token == "" {
		return "", ErrEmptyParamToken
	}

	return token, nil
}

// ParseToken 解析jwt令牌
func (mw *GfJWTMiddleware) ParseToken(r *ghttp.Request) (*jwt.Token, error) {
	var token string
	var err error

	methods := strings.Split(mw.TokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = mw.jwtFromHeader(r, v)
		case "query":
			token, err = mw.jwtFromQuery(r, v)
		case "cookie":
			token, err = mw.jwtFromCookie(r, v)
		case "param":
			token, err = mw.jwtFromParam(r, v)
		}
	}

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
		}

		// 如果有效，保存令牌字符串
		r.SetParam("JWT_TOKEN", token)

		return mw.Key, nil
	})
}

func (mw *GfJWTMiddleware) unauthorized(r *ghttp.Request, code int, message string) {
	r.Header.Set("WWW-Authenticate", "JWT realm="+mw.Realm)
	mw.Unauthorized(r, code, message)
	if !mw.DisabledAbort {
		r.ExitAll()
	}

}

func (mw *GfJWTMiddleware) setBlacklist(token string, claims jwt.MapClaims) error {
	// MD5的目的是减少密钥长度。
	token, err := gmd5.EncryptString(token)

	if err != nil {
		return err
	}

	exp := int64(claims["exp"].(float64))

	// 全局缓存（gcache）
	err = blacklist.Set(token, true, time.Unix(exp, 0).Sub(mw.TimeFunc()))
	if err != nil {
		return err
	}

	return nil
}

func (mw *GfJWTMiddleware) inBlacklist(token string) (bool, error) {
	// MD5的目标是减少密钥长度。
	tokenRaw, err := gmd5.EncryptString(token)

	if err != nil {
		return false, nil
	}

	// 全局缓存（gcache）
	if in, err := blacklist.Contains(tokenRaw); err != nil {
		return false, nil
	} else {
		return in, nil
	}
}

// ExtractClaims 帮助提取JWT声明（claims）
func ExtractClaims(r *ghttp.Request) MapClaims {
	claims := r.GetParam("JWT_PAYLOAD")
	return claims.(MapClaims)
}

// GetToken 帮助获取JWT令牌字符串
func GetToken(r *ghttp.Request) string {
	token := r.GetString("JWT_TOKEN")
	if len(token) == 0 {
		return ""
	}

	return token
}
