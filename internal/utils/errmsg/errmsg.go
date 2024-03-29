package errmsg

const (
	SUCCES = 200
	ERROR  = 500

	// 绑定参数错误
	ERROR_BIND          = 0001
	ERROR_AUTH          = 0002
	ERROR_INVAILD_TOKEN = 0003
	ERROR_INVAILD_ARGV  = 0004

	// ERROR_USERNAME_USED code= 1000... 用户模块的错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_TOKEN_EXIST      = 1004
	ERROR_TOKEN_RUNTIME    = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT    = 1008
	ERROR_UPDATE_USER      = 1009
	ERROR_LIST_USER        = 1010
	ERROR_DELETE_USER      = 1011
	ERROR_GET_USER         = 1012
	ERROR_OAUTH            = 1013

	// code= 2000... 文章模块的错误
	ERROR_ART_NOT_EXIST = 2001

	// code= 3000... 分类模块的错误
	ERROR_CATENAME_USED   = 3001
	ERROR_CATE_NOT_EXIST  = 3002
	ERROR_UPDATE_CATENAME = 3003
	ERROR_DELETE_CATENAME = 3004
	ERROR_LIST_CATEGORY   = 3005

	// code=4000... 奖金错误
	ERROR_BONUSE_EXIST     = 4001
	ERROR_BONUSE_NOT_EXIST = 4002
	ERROR_LIST_BONUSE      = 4003
	ERROR_GET_BONUSE       = 4004
	ERROR_DELETE_BONUSE    = 4005
	ERROR_CREATE_BONUSE    = 4006
	ERROR_UPDATE_BONUSE    = 4007

	// code=5000... 话题错误
	ERROR_UPDATE_TOPIC       = 5001
	ERROR_DELETE_TOPIC       = 5002
	ERROR_LIST_TOPIC         = 5003
	ERROR_GET_TOPIC          = 5003
	ERROR_TOPIC_INVAILD_NAME = 5005
	ERROR_TOPICS_NOT_FOUND   = 5006

	// 收藏错误
	ERROR_CREATE_COLLECT         = 6001
	ERROR_CREATE_ALREADY_COLLECT = 6002
	ERROR_NOT_COLLECT            = 6003

	// 购买记录为空
	ERROR_BUY_RECORD_EMPTY = 7001
)

var codeMsg = map[int]string{
	SUCCES:                 "SUCCESS",
	ERROR:                  "FAIL",
	ERROR_BIND:             "输入参数错误",
	ERROR_AUTH:             "认证错误",
	ERROR_INVAILD_TOKEN:    "token错误",
	ERROR_INVAILD_ARGV:     "参数校验错误",
	ERROR_OAUTH:            "oauth错误",
	ERROR_USERNAME_USED:    "用户名已存在！",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_TOKEN_EXIST:      "TOKEN不存在,请重新登陆",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期,请重新登陆",
	ERROR_TOKEN_WRONG:      "TOKEN不正确,请重新登陆",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误,请重新登陆",
	ERROR_USER_NO_RIGHT:    "该用户无权限",
	ERROR_UPDATE_USER:      "更新用户失败",
	ERROR_LIST_USER:        "查询用户列表失败",
	ERROR_DELETE_USER:      "删除用户失败",
	ERROR_GET_USER:         "获取用户失败",

	ERROR_ART_NOT_EXIST: "文章不存在",

	ERROR_CATENAME_USED:   "该分类已存在",
	ERROR_CATE_NOT_EXIST:  "该分类不存在",
	ERROR_UPDATE_CATENAME: "更新分类错误",
	ERROR_DELETE_CATENAME: "删除分类错误",
	ERROR_LIST_CATEGORY:   "查询分类列表错误",

	ERROR_BONUSE_EXIST:     "奖金已存在",
	ERROR_BONUSE_NOT_EXIST: "奖金不存在",
	ERROR_LIST_BONUSE:      "查询奖金列表失败",
	ERROR_UPDATE_BONUSE:    "更新奖金失败",
	ERROR_GET_BONUSE:       "查询奖金失败",
	ERROR_DELETE_BONUSE:    "删除奖金失败",
	ERROR_CREATE_BONUSE:    "创建奖金失败",

	ERROR_UPDATE_TOPIC:       "更新话题错误",
	ERROR_DELETE_TOPIC:       "删除话题错误",
	ERROR_LIST_TOPIC:         "查询话题列表错误",
	ERROR_TOPIC_INVAILD_NAME: "invaild tid",
	ERROR_TOPICS_NOT_FOUND:   "话题未找到",

	ERROR_CREATE_COLLECT:         "创建收藏失败",
	ERROR_CREATE_ALREADY_COLLECT: "收藏已存在",
	ERROR_NOT_COLLECT:            "收藏不存在",

	ERROR_BUY_RECORD_EMPTY: "购买记录为空",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
