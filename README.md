# HQ-stocks

### 简介

>* 这是一个Go实现的简单策略的股票**图形化**程序；
>* 选择了**聚宽**数据平台的股票数据；
>* 主要通过**2015-2019**的财务数据，选择了平均Roe大于15%的A股股票，
>* 展示了 **ROE 、PE、PB**柱状图，同时将三个数据作为雷达图展示

### 示例图

![ROE](.\src\assets\ROE.png)


### 项目依赖
>* 本工程依赖于 go-echarts 库，请先安装该 https://github.com/go-echarts/go-echarts
>* 最新股票数据，依赖于聚宽数据平台，如果需要强制获取最新的，可以注册账号获取https://www.joinquant.com/ 


### 文件说明
>* ./src/assets/conf.json 基础配置文件
>* ./src/assets/stock.json 默认股票数据文件
>* **bar.html 静态文件，可以不依赖程序直接使用**

### 配置说明
>* ./src/assets/conf.json 基础配置文件
``` json
{
    "forceLoad":false,//是否强制加载最新的股票数据
    "host":"http://127.0.0.1:8080",// 本地服务地址配置
    "minRoe":15.0,// 从聚宽数据平台获取股票数据的最小roe
    "maxRoe":1000.0,// 从聚宽数据平台获取股票数据的最大roe
    "showMinRoe":30.0,// 展示时，显示的最小ROE
    "showMaxRoe":55.0,// 展示时，显示的最大ROE
    "apiAccount":"18200000000", // 聚宽数据平台账号[如果强只获取，需要配置]
    "apiPwd":"ECD928a6bf8eec8ba1"// 聚宽数据平台密码[如果强只获取，需要配置]
}
```

