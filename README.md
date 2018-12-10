# xueqiuStock

use:

    docker run  -p127.0.0.1:27017:27017 -d -v /root/xueqiu/data:/data/db mongo:4.1.6
    chmod a+x xueqiuStock
    ./xueqiuStock

default fetch 200 day period

you can run with parameter, use -h parameter for help.

    ./xueqiuStock -h
    Usage of ./xueqiuStock_linux_amd64:
      -GetStockListThreadNum int
        	Get StockList Thread Num (default 10)
      -GetStockThreadNum int
        	Get Stock Thread Num (default 25)
      -MongoDBName string
        	MongoDB database name (default "xueqiuStock")
      -MongoDBUrl string
        	MongoDB connect url (default "127.0.0.1")
      -StockBeginTime string
        	Stock Begin Time (default "1544433839668")
      -StockCount string
        	Stock Count (default "-200")
      -StockPeriod string
        	Stock Period (default "day")
      -StockType string
        	Stock Type (default "before")

example:

    ./xueqiuStock -StockPeriod "week" -StockCount "-500" // fetch 500 week stock data
    ./xueqiuStock -StockPeriod "month" -StockCount "-100" // fetch 100 month stock data

