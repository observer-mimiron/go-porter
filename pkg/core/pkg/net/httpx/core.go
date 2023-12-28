package httpx

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cors "github.com/rs/cors/wrapper/gin"
	"go-porter/configs"
	"go-porter/internal/errCode"
	"go-porter/pkg/core/pkg/conf/env"
	"go-porter/pkg/core/pkg/proposal"
	"go-porter/pkg/core/pkg/trace"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
)

func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcs[i] = func(c *gin.Context) {
			ctx := newContext(c)
			defer releaseContext(ctx)

			handler(ctx)
		}
	}

	return funcs
}

var _ Mux = (*mux)(nil)

// Mux http mux
type Mux interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type mux struct {
	engine *gin.Engine
}

func (m *mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.engine.ServeHTTP(w, req)
}

func (m *mux) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{
		group: m.engine.Group(relativePath, wrapHandlers(handlers...)...),
	}
}

//panicHandler panic处理
func (m *mux) panicHandler(logger *zap.Logger) gin.IRoutes {
	return m.engine.Use(func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", string(debug.Stack())))
			}
		}()

		ctx.Next()
	})
}

//errorHandler 错误处理
func (m *mux) errorHandler(logger *zap.Logger, opt *option) gin.IRoutes {
	return m.engine.Use(func(ctx *gin.Context) {
		if ctx.Writer.Status() == http.StatusNotFound {
			return
		}

		ts := time.Now()

		context := newContext(ctx)
		defer releaseContext(context)

		context.init()
		context.setLogger(logger)
		context.ableRecordMetrics()

		if !opt.outTracePaths[ctx.Request.URL.Path] {
			if traceId := context.GetHeader(trace.Header); traceId != "" {
				context.setTrace(trace.New(traceId))
			} else {
				context.setTrace(trace.New(""))
			}
		}

		defer func() {
			var (
				response        interface{}
				businessCode    int
				businessCodeMsg string
				abortErr        error
				traceId         string
			)

			if ct := context.Trace(); ct != nil {
				context.SetHeader(trace.Header, ct.ID())
				traceId = ct.ID()
			}

			// region 发生 Panic 异常发送告警提醒
			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				//转义
				stackInfo = strings.Replace(stackInfo, "\n", "\\n", -1)
				logger.Error("got panic",
					zap.String("panic", fmt.Sprintf("%+v", err)),
					zap.String("stack", stackInfo),
				)

				context.AbortWithError(errCode.ErrServer)

				if notifyHandler := opt.alertNotify; notifyHandler != nil {
					notifyHandler(&proposal.AlertMessage{
						ProjectName:  configs.ProjectName,
						Env:          env.Active().Value(),
						TraceID:      traceId,
						HOST:         context.Host(),
						URI:          context.URI(),
						Method:       context.Method(),
						ErrorMessage: err,
						ErrorStack:   stackInfo,
						Timestamp:    time.Now(),
					})
				}
			}

			// endregion
			// region 发生错误，进行返回 错误处理
			if ctx.IsAborted() {
				for i := range ctx.Errors {
					multierr.AppendInto(&abortErr, ctx.Errors[i])
				}

				if err := context.abortError(); err != nil { // customer err
					// 判断是否需要发送告警通知
					//if err.IsAlert() {
					//	if notifyHandler := opt.alertNotify; notifyHandler != nil {
					//		notifyHandler(&proposal.AlertMessage{
					//			ProjectName:  configs.ProjectName,
					//			Env:          env.Active().Value(),
					//			TraceID:      traceId,
					//			HOST:         context.Host(),
					//			URI:          context.URI(),
					//			Method:       context.Method(),
					//			ErrorMessage: err.Message(),
					//			ErrorStack:   fmt.Sprintf("%+v", err.StackError()),
					//			Timestamp:    time.Now(),
					//		})
					//	}
					//}

					errcode := errCode.ErrServer.ErrCode()
					errmsg := errCode.ErrServer.Message()
					causeErr := errors.Cause(err)

					if e, ok := causeErr.(*errCode.ErrCode); ok { //自定义错误类型
						errcode = e.ErrCode()
						errmsg = e.Message()
					}

					type Code struct {
						Code    int32  `json:"code"`    // 业务码
						Message string `json:"message"` // 描述信息
					}

					ctx.JSON(http.StatusUnauthorized, &Code{
						Code:    errcode,
						Message: errmsg,
					})
				}

			}
			// endregion

			// region 正确返回
			response = context.getPayload()
			if response != nil {
				ctx.JSON(http.StatusOK, response)
			}
			// endregion

			// region 记录指标
			if opt.recordHandler != nil && context.isRecordMetrics() {
				path := context.Path()
				if alias := context.Alias(); alias != "" {
					path = alias
				}

				opt.recordHandler(&proposal.MetricsMessage{
					ProjectName:  configs.ProjectName,
					Env:          env.Active().Value(),
					TraceID:      traceId,
					HOST:         context.Host(),
					Path:         path,
					Method:       context.Method(),
					HTTPCode:     ctx.Writer.Status(),
					BusinessCode: businessCode,
					CostSeconds:  time.Since(ts).Seconds(),
					IsSuccess:    !ctx.IsAborted() && (ctx.Writer.Status() == http.StatusOK),
				})
			}
			// endregion

			// region 记录日志
			var t *trace.Trace
			if x := context.Trace(); x != nil {
				t = x.(*trace.Trace)
			} else {
				return
			}

			decodedURL, _ := url.QueryUnescape(ctx.Request.URL.RequestURI())

			// ctx.Request.Header，精简 Header 参数
			traceHeader := map[string]string{
				"Content-Type":              ctx.GetHeader("Content-Type"),
				configs.HeaderLoginToken:    ctx.GetHeader(configs.HeaderLoginToken),
				configs.HeaderSignToken:     ctx.GetHeader(configs.HeaderSignToken),
				configs.HeaderSignTokenDate: ctx.GetHeader(configs.HeaderSignTokenDate),
			}

			t.WithRequest(&trace.Request{
				TTL:        "un-limit",
				Method:     ctx.Request.Method,
				DecodedURL: decodedURL,
				Header:     traceHeader,
				Body:       string(context.RawData()),
			})

			var responseBody interface{}

			if response != nil {
				responseBody = response
			}

			t.WithResponse(&trace.Response{
				Header:          ctx.Writer.Header(),
				HttpCode:        ctx.Writer.Status(),
				HttpCodeMsg:     http.StatusText(ctx.Writer.Status()),
				BusinessCode:    businessCode,
				BusinessCodeMsg: businessCodeMsg,
				Body:            responseBody,
				CostSeconds:     time.Since(ts).Seconds(),
			})

			t.Success = !ctx.IsAborted() && (ctx.Writer.Status() == http.StatusOK)
			t.CostSeconds = time.Since(ts).Seconds()

			logger.Info("trace-log",
				zap.Any("method", ctx.Request.Method),
				zap.Any("path", decodedURL),
				zap.Any("http_code", ctx.Writer.Status()),
				zap.Any("business_code", businessCode),
				zap.Any("success", t.Success),
				zap.Any("cost_seconds", t.CostSeconds),
				zap.Any("trace_id", t.Identifier),
				zap.Any("trace_info", t),
				zap.Error(abortErr),
			)
			// endregion
		}()

		ctx.Next()
	})
}

func New(logger *zap.Logger, options ...Option) (Mux, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	gin.SetMode(gin.ReleaseMode)
	mux := &mux{
		engine: gin.New(),
	}

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	if !opt.disablePProf {
		if !env.Active().IsPro() {
			pprof.Register(mux.engine) // register pprof to gin
		}
	}

	if !opt.disablePrometheus {
		mux.engine.GET("/metrics", gin.WrapH(promhttp.Handler())) // register prometheus
	}

	if opt.enableCors {
		mux.engine.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:     []string{"*"},
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}

	mux.panicHandler(logger)

	mux.errorHandler(logger, opt)

	//限流器
	if opt.enableRate {
		limiter := rate.NewLimiter(rate.Every(time.Second*1), configs.MaxRequestsPerSecond)
		mux.engine.Use(func(ctx *gin.Context) {
			context := newContext(ctx)
			defer releaseContext(context)

			if !limiter.Allow() {
				context.AbortWithError(errCode.ErrTooManyRequests)
				//context.AbortWithError(errors.Error(
				//	http.StatusTooManyRequests,
				//	ecode.TooManyRequests,
				//	ecode.Text(ecode.TooManyRequests)),
				//)
				return
			}

			ctx.Next()
		})
	}
	mux.engine.NoMethod(wrapHandlers(DisableTraceLog)...)
	mux.engine.NoRoute(wrapHandlers(DisableTraceLog)...)

	system := mux.Group("/system")
	{
		// 健康检查
		system.GET("/health", func(ctx Context) {
			resp := &struct {
				Timestamp   time.Time `json:"timestamp"`
				Environment string    `json:"environment"`
				Host        string    `json:"host"`
				Status      string    `json:"status"`
			}{
				Timestamp:   time.Now(),
				Environment: env.Active().Value(),
				Host:        ctx.Host(),
				Status:      "ok",
			}
			ctx.Payload(resp)
		})
	}

	return mux, nil
}
