package aop

import (
	"reflect"
	"regexp"
	"bou.ke/monkey"
)

//连接点
type JoinPoint struct {
	Receiver 			interface{}
	Method   			reflect.Method//方法
	Params   			[]reflect.Value
	Result   			[]reflect.Value//结果
}
//这里是挂在一个新函数
func NewJoinPoint(receiver interface{}, params []reflect.Value, method reflect.Method) *JoinPoint {
	point := &JoinPoint{
		Receiver: receiver,
		Params: params,
		Method: method,
	}
	fn := method.Func
	fnType := fn.Type()//type为interface{}的具体类型，
	nout := fnType.NumOut()
	point.Result = make([]reflect.Value, nout)
	for i := 0; i < nout; i++ {
		//默认返回空值
		point.Result[i] = reflect.Zero(fnType.Out(i))
	}
	return point
}

//切面接口
type AspectInterface interface {
	Before(point *JoinPoint) bool
	After(point *JoinPoint)
	Finally(point *JoinPoint)
	GetAspectExpress() string
}

//切面列表,生成一个切片
var aspectList = make([]AspectInterface, 0)

//注册切点
func RegisterPoint(pointType reflect.Type) {
	pkgPth := pointType.PkgPath()//获取包的名称
	receiverName := pointType.Name()//获取pointtype在包中的名称
	/*
		kind获取类别，ptr是一个const常量表示指针类型，
		判断pointtype是不是指针，如果是指针的话，重新获取下包名
		Elem()通过反射回去指针指向的元素类型
		这里上下两段代码的目的都是为了获取传入过来的reflect.Type的包名和类名
	*/
	if pointType.Kind() == reflect.Ptr {
		/* Elem可以获取元素的类型，指针所指对象的类型，获取接口的动态类型 */
		pkgPth = pointType.Elem().PkgPath()//pointtype所在的包名
		receiverName = pointType.Elem().Name()//类型名
	}


	//枚举方法
	for i := 0; i < pointType.NumMethod(); i++ {//nummethod可以获取到方法的数量
		method := pointType.Method(i)//根据索引获取对应的方法
		/*
			方法位置字符串"包名.接收者.方法名"，用于匹配代理
			methodLocation是一个地址字符串，是包名.结构体名.方法名
		*/
		methodLocation := fmt.Sprintf("%s.%s.%s", pkgPth, receiverName, method.Name)
		//这里用到了monkey框架来做动态代理，这部分还没太看懂monkey包怎么用
		var guard *monkey.PatchGuard//patchguard英文翻译叫：补丁 守卫

		//proxy是一个匿名函数
		var proxy = func(in []reflect.Value) []reflect.Value {
			//patch是打桩，unpatch是不打桩，这里没有用到打桩的知识点
			guard.Unpatch()
			defer guard.Restore()
			receiver := in[0]
			point := NewJoinPoint(receiver, in[1:], method)
			defer finallyProcessed(point, methodLocation)
			if !beforeProcessed(point, methodLocation) {
				return point.Result
			}
			point.Result = receiver.MethodByName(method.Name).Call(in[1:])
			afterProcessed(point, methodLocation)
			return point.Result
		}
		//动态创建代理函数
		proxyFn := reflect.MakeFunc(method.Func.Type(), proxy)
		//将pointtype方法跳转到补丁方法proxyFn
		guard = monkey.PatchInstanceMethod(pointType, method.Name, proxyFn.Interface())
	}
}

//注册切面
func RegisterAspect(aspect AspectInterface) {
	aspectList = append(aspectList, aspect)
}

//前置处理
func beforeProcessed(point *JoinPoint, methodLocation string) bool {
	for _, aspect := range aspectList {
		if !isAspectMatch(aspect.GetAspectExpress(), methodLocation) {
			continue
		}
		if !aspect.Before(point) {
			return false
		}
	}
	return true
}

//后置处理
func afterProcessed(point *JoinPoint, methodLocation string) {
	for i := len(aspectList) - 1; i >= 0; i-- {
		aspect := aspectList[i]
		//如果匹配的话
		if !isAspectMatch(aspect.GetAspectExpress(), methodLocation) {
			continue
		}
		aspect.After(point)
	}
}

//最终处理
func finallyProcessed(point *JoinPoint, methodLocation string) {
	for i := len(aspectList) - 1; i >= 0; i-- {
		aspect := aspectList[i]
		if !isAspectMatch(aspect.GetAspectExpress(), methodLocation) {
			continue
		}
		aspect.Finally(point)
	}
}


//传入的参数一个是切面描述，一个是方法地址
func isAspectMatch(aspectExpress, methodLocation string) bool {
	//aspectExpress采用正则表达式
	pattern, err := regexp.Compile(aspectExpress)
	if err != nil {
		return false
	}
	return pattern.MatchString(methodLocation)
}

