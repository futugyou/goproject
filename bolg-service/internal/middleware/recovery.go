package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goproject/blog-service/global"
	"github.com/goproject/blog-service/pkg/app"
	"github.com/goproject/blog-service/pkg/errcode"
)

func Recovery() gin.HandlerFunc {
	// defailtMailer := email.NewEmail(&email.SMTPInfo{
	// 	Host:     global.EmailSetting.Host,
	// 	Port:     global.EmailSetting.Port,
	// 	IsSSL:    global.EmailSetting.IsSSL,
	// 	UserName: global.EmailSetting.UserName,
	// 	Password: global.EmailSetting.Password,
	// 	From:     global.EmailSetting.From,
	// })
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallersFrames().Errorf(c, "panic recover err : %v", err)

				// err := defailtMailer.SendEmail(
				// 	global.EmailSetting.To,
				// 	fmt.Sprintf("error happen time : %d", time.Now().Unix()),
				// 	fmt.Sprintf("error message : %v", err),
				// )
				// if err != nil {
				// 	global.Logger.Panicf("mail SendEmail err: %v", err)
				// }

				app.NewResponse(c).ToErrorResponse(errcode.ServiceError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
