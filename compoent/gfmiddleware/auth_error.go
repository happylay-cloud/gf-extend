package gfmiddleware

import "errors"

var (
	// ErrMissingSecretKey 表示需要密钥
	ErrMissingSecretKey = errors.New("需要密钥")

	// ErrForbidden 当给出HTTP状态403时
	ErrForbidden = errors.New("您没有访问此资源的权限")

	// ErrMissingAuthenticatorFunc 表示需要验证器
	ErrMissingAuthenticatorFunc = errors.New("GfJWTMiddleware.Authenticator函数未定义")

	// ErrMissingLoginValues 表示试图在没有用户名或密码的情况下进行身份验证的用户
	ErrMissingLoginValues = errors.New("缺少用户名或密码")

	// ErrFailedAuthentication 表示身份验证失败，可能是错误的用户名或密码
	ErrFailedAuthentication = errors.New("用户名或密码不正确")

	// ErrFailedTokenCreation 表示创建JWT令牌失败，原因未知
	ErrFailedTokenCreation = errors.New("创建JWT令牌失败")

	// ErrExpiredToken 表示JWT令牌已过期。不能刷新。
	ErrExpiredToken = errors.New("令牌已过期")

	// ErrInvalidToken 表示JWT令牌无效。不能刷新。
	ErrInvalidToken = errors.New("令牌无效")

	// ErrEmptyAuthHeader 如果使用HTTP请求头进行身份验证，则需要设置身份验证请求头
	ErrEmptyAuthHeader = errors.New("认证请求头为空")

	// ErrMissingExpField 令牌中缺少exp字段
	ErrMissingExpField = errors.New("缺少exp字段")

	// ErrWrongFormatOfExp 字段必须是float64格式
	ErrWrongFormatOfExp = errors.New("exp必须是float64格式")

	// ErrInvalidAuthHeader 表示认证头是无效的，例如可能有错误的域名
	ErrInvalidAuthHeader = errors.New("认证请求头无效")

	// ErrEmptyQueryToken 如果使用URL查询进行身份验证，则查询令牌变量为空
	ErrEmptyQueryToken = errors.New("查询令牌为空")

	// ErrEmptyCookieToken 如果使用cookie进行身份验证，则可以抛出标记cokie为空
	ErrEmptyCookieToken = errors.New("cookie令牌为空")

	// ErrEmptyParamToken 如果使用路径中的参数进行身份验证，则可以抛出路径中的参数为空
	ErrEmptyParamToken = errors.New("参数令牌为空")

	// ErrInvalidSigningAlgorithm 表示签名算法无效，需要 HS256, HS384, HS512, RS256, RS384 或 RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile 表示给定的私钥不可读
	ErrNoPrivKeyFile = errors.New("私钥文件不可读")

	// ErrNoPubKeyFile 表示给定的公钥不可读
	ErrNoPubKeyFile = errors.New("公钥文件不可读")

	// ErrInvalidPrivKey 表示给定的私钥无效
	ErrInvalidPrivKey = errors.New("私钥无效")

	// ErrInvalidPubKey 表示给定的公钥无效
	ErrInvalidPubKey = errors.New("公钥无效")
)
