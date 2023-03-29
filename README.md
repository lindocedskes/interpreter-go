## writing an interpreter in go中文版-跟着书本敲的代码

目前学习到：+ 2.9 RPPL 语法分析器完成

学完章节ed:
+ 1.3 简单词法分析器

+ 1.4 扩展词法分析器

+ 1.5 编写REPL

+ 2.4 语法分析器第一步：解析let语句

+ 2.5 解析return语句

+ 2.6 解析表达式

  + 2.6.2 普拉特解析法

  + 2.6.4 准备AST ：每类（let,return,exepresion）statement都支持string()，返回对应AST信息供查看
  
  + 2.6.6 标识符
  + 2.6.7 整数字面量
  + 2.6.8 解析前缀运算符!-
  + 2.6.9 解析中缀运算符

+ 2.7 普拉特解析工作方式
+ 2.8 扩展语法分析器
  + 2.8.1 布尔字面量：支持布尔类型解析，在前缀和中缀布尔添加测试
  + 2.8.2 分组表达式：带括号的表达式处理
  + 2.8.3 if表达式：if-else，else可选，包含{}解析函数
  + 2.8.4 函数字面量：fn <parameters> <block statement>
  + 2.8.5 调用表达式： <expression>(<comma separated expressions>) 例：add()
  + 2.8.6 删除TODO ：let 和return语句支持表达式
  + 2.9 RPPL :读取-语法分析-输出 循环 模拟控制台

tag版本解释

+ v2.2 简单词法分析器，支持数字/标识符的表达式语法解析，对应2.7完成
+ v2.1 解析return语句+解析表达式之准备AST
+ v2.0 语法分析器第一步：解析let语句
+ v1.2 REPL （读取-求值-打印 循环）实现用户输入命令，词法分析器求值后，向系统out输出词法单元。并用for实现循环。词法分析器实现完成
+ v1.1 扩展词法分析器：识别更全+支持双字词法单元
+ v1.0 简单词法分析器
