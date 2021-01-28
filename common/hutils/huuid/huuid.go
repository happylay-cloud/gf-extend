package huuid

import (
	"github.com/gogf/gf/text/gstr"
	"github.com/google/uuid"
)

// UUID 获取随机32位uuid
//  格式：83a34b1231974072b76d38130482ad47
func UUID() string {
	return gstr.Replace(uuid.New().String(), "-", "")
}
