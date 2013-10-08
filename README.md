## Meanserver gives you weather data…

It's pretty grumpy and broken, though. 

### Output format

Meanserver will give you weather data in the following json format

	{"Temperature":99,"Conditions":"cloudy"}
	
It's not going to do it consistently, though.

In fact, meanserver breaks in strange ways. 

### Throttling

Meanserver doesn't like to work hard. If you ask it for the weather more than once in a 2 second window, meanserver won't let you talk to it for another 4 seconds. 

If you ask again during those 4 seconds, the 4 second window is reset. 
	
