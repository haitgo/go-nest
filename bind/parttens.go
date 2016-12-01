package bind

import (
	"go-nest/cache"
	"regexp"
)

//过滤回调函数
type ValidCall func(str *string) bool

var (
	customValids = cache.NewMemory(0)             //用户添加的正则判断规则
	compiles     = cache.NewMemory(len(Parttens)) //正则编译map，保证只实现一次regexp.compile
)

//注册自定义过滤规则的函数，比如可以做一些关键词过滤的处理等
func RegValid(key string, call ValidCall) {
	customValids.Set(key, call)
}

//正则判断
//str字符串，parteen预设正则或者自定义正则，验证成功返回true，否则返回false
//partten 参见parttens.go
func Regexp(str, partten string) bool {
	if partten == "-" { //如果为不开启验证
		return true
	}
	if call := customValids.Get(partten); call != nil { //优先判断是否是自定义过滤函数
		cbk := call.(ValidCall)
		return cbk(&str)
	}
	var pt *regexp.Regexp
	p := compiles.Get(partten)
	if p != nil {
		pt = p.(*regexp.Regexp)
	} else if partten, ok := Parttens[partten]; ok {
		cp, err := regexp.Compile(partten) //编译正则
		if err != nil {
			return false
		}
		compiles.Set(partten, cp)
		pt = cp
	} else {
		return false
	}
	return pt.Match([]byte(str))
}

//正则式集合，可以自行添加
var Parttens = map[string]string{
	"Text":           `^[\S\s]+$`,             //任意类型
	"HexCmd":         `^(\d{1,3}\t)+\d{1,3}$`, //16进制modbus命令
	"String":         `^[\S]+$`,
	"Date":           `^\d{4}(\-|\/|\.)\d{1,2}(\-|\/|\.)\d{1,2}$`, //这里只是一个简单的日期判断
	"Phone":          `^[1]\d{10}$`,                               //简单判断手机号码
	"Md5":            `^[a-zA-Z0-9]{32}$`,                         //md5加密
	"Email":          "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$",
	"CreditCard":     "^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$",
	"ISBN10":         "^(?:[0-9]{9}X|[0-9]{10})$",
	"ISBN13":         "^(?:[0-9]{13})$",
	"UUID3":          "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$",
	"UUID4":          "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
	"UUID5":          "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
	"UUID":           "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
	"Alpha":          "^[a-zA-Z]+$",
	"Alphanumeric":   "^[a-zA-Z0-9]+$",
	"Numeric":        "^[-+]?[0-9]+$",
	"Int":            "^(?:[-+]?(?:0|[1-9][0-9]*))$",
	"Float":          "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$",
	"Hexadecimal":    "^[0-9a-fA-F]+$",
	"Hexcolor":       "^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$",
	"RGBcolor":       "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$",
	"ASCII":          "^[\x00-\x7F]+$",
	"Multibyte":      "[^\x00-\x7F]",
	"FullWidth":      "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]",
	"HalfWidth":      "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]",
	"Base64":         "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$",
	"PrintableASCII": "^[\x20-\x7E]+$",
	"DataURI":        "^data:.+\\/(.+);base64$",
	"Latitude":       "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$",
	"Longitude":      "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$",
	"DNSName":        `^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{1,62}){1}(.[a-zA-Z0-9]{1}[a-zA-Z0-9_-]{1,62})*$`,
	"URL":            `^((ftp|https?):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`,
	"SSN":            `^\d{3}[- ]?\d{2}[- ]?\d{4}$`,
	"WinPath":        `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`,
	"UnixPath":       `^((?:\/[a-zA-Z0-9\.\:]+(?:_[a-zA-Z0-9\:\.]+)*(?:\-[\:a-zA-Z0-9\.]+)*)+\/?)$`,
	"Semver":         "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$",
}
