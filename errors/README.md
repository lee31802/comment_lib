无论使用什么语言，错误处理都是系统中很关键的一个点。优雅的错误处理能够极大的提高代码的整洁度，而代码整洁度又直接影响可维护性，但是要实现优雅的错误处理却并非易事。那究竟难在哪里呢？
https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
1. 重复的错误处理代码

   > ``` go
   > _, err = fd.Write(p2[e:f])
   > if err != nil {
   >  fmt.Println(err, ...)
   >  return err
   > }
   > ```
   >
   > 在分层系统中（例如：Controller、Service、DAO），每一层都会重复以上代码

2. 原始错误的上下文

   > ``` go
   > func AuthenticateRequest(r *Request) error {
   > 	return authenticate(r.User)
   > }
   > ```
   >
   > 如果 `authenticate` 返回错误，那么 `AuthenticateRequest` 会将错误返回给调用者，调用者也可能会这样做，依此类推。 在程序的顶部，程序的主体将错误打印到屏幕或日志文件，所有打印的都会是： `No such file or directory`

3. 原始错误 vs 错误码错误

   > ``` go
   > errors.New("connection error")
   > 
   > // vs 
   > 
   > var errno uint32 = 10001
   > errors.New(errno, "connection error")
   > ```
   >
   > 标准库或第三方应用库返回的是一般是上一种错误，然而H5或APP不可能根据字符串进行错误判断。因此在业务中需要把所有的错误进行统一封状为错误码错误返回，那么错误码错误中需要保存原始错误么？

4. RPC错误 vs 业务错误

   > - 一方面，与业务层的错误相似，RPC框架也会自己的错误。而RPC框架往往也会集成一些包括过载处理、异常节点剔除的功能，依赖于对两种错误的识别能力。该如何设计错误才能让两者区分开来呢？
   > - 另一方面，无论是RPC框架错误和业务错误，调用端都需要进行统一解码（decode）。该如何设计错误才能让两者融合起来呢？


### 错误模型

简单来看，所有问题是相互独立的，但是透过现象来看本质。以上问题又都有关联，在于进行错误模型设计。从业界各种框架的设计情况来看，可以把错误分为以下三种：

![](./doc/error-model.png)

- **Error codes model**

  > Errors are raised under various circumstances, from network failures to unauthenticated connections, each of which is associated with a particular code. 

  > 从网络故障到未经验证的连接，各种情况下都会引发错误，每种错误都可以都与特定错误码关联。

- **Standard error model**
  > If an error occurs, return error codes instead, with an optional string error message that provides further details about what happened. 

  > 如果发生错误，则返回错误代码以及一条可选的字符串错误信息，该信息提供有关所发生事件的详细信息。

- **Richer error model**
  > Enables servers to return and clients to consume additional error details. It further specifies a standard set of error message types to cover the most common needs (such as invalid parameters, quota violations, and stack traces). 

  > 允许服务器和客户端返回、使用额外的错误详细信息。它进一步指定了一组标准的错误消息类型，以满足最常见的需求（例如无效参数、配额冲突和堆栈跟踪）


三种类型的错误，层层递进，能够囊括的信息也越来越多。当然，信息越是丰富，框架实现难度越高，对使用者也越友好。在微信，svrkit 选择的模型是`Error codes model`；开源框架 grpc 选择的模型则是` Standard error model` ，但是本身支持`Richer error model`，参考：[googleapis](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto)

### 问题解决
想清楚了问题，再看解决问题的方案：

- 难点 1、2：可以按照以下两种思想来解决：[《Errors are values》](https://blog.golang.org/errors-are-values)、[《Don’t just check errors, handle them gracefully》](https://www.cyningsun.com/09-09-2019/dont-just-check-errors-handle-them-gracefully-cn.html)

- 难点 3：错误的主要作用有以下两点：

  > - 根据错误的类型，进行针对性的处理
  > - 错误原因追踪
  >
  > 前者可以使用错误码来代替，后者可以简化为message。即，毋需保留错误本身，只需要将错误转化为错误码和message。

- 难点 4：即根据需要选择合适的错误模型，统一业务错误和框架错误。区分业务错误和框架错误，可以将 code 分段，框架优先占有指定的号段。

  > 更进一步，code 号段可以融合到服务治理中，在服务申请阶段分配对应的号段


### 实践
toC 选择 `Standard error model` 为基础实践：
1. 业务方应当自行在业务层自行定义 error code，并在系统中使用 `errors.New()` 将各种错误类型转换为 `Standard error model` 
    > error code 应当 **> `1000`**，避免与 `grpc` 和 `microkit` 的相关错误冲突
2. 在服务端，业务方应使用 
    > - `errors.Wrap` 在每一层包裹下层错误以及所需携带的参数，并逐层传递
    > - 在RPC入口使用 `errors.As()` 返回 `Standard error model` 
    > - 在RPC入口使用 `logkit` 打印 error ， `logkit` 会自动 `errors.unwrap` 输出所有携带数据
3. 在调用端，业务方应使用 `errors.As()` 返回 `Standard error model` 


### 参考资料

- [**Error Handling**: How gRPC deals with errors, and gRPC error codes](https://grpc.io/docs/guides/error/)
- [Google API 错误模型](https://cloud.google.com/apis/design/errors)
