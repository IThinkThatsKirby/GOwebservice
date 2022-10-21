# GOwebservice
 roast my line comments and readme PLS :D
 * If this is your first time with the GO language 
	please visit [https://go.dev/dl/](https://go.dev/dl/)
 * Follow the instructions for your system. 
 * Clone this repo to a folder.
 * Open a terminal inside the folder.
 * Enter the command: go run test.go
	
	 
# what it DO:
  1.  This web service provides proccessing for payer transactions with timestamps.
    example: 
```
    POST { "payer": "DANNON", "points": 1000, "timestamp": "2020-11-02T14:00:00Z" }
```

  2.  Spend points from oldest to newest by transaction date regardless of payer.
  example:
 ```
     POST { "points": 5000 }
 
    
   * response from server:
 
    [
      { "payer": "DANNON", "points": -100 },
      { "payer": "UNILEVER", "points": -200 },
      { "payer": "MILLER COORS", "points": -4,700 }
    ]
 ```   
  
  3.  return current balance of all payers.
    example: 
 ``` 
     	{
        "DANNON": 1000,
        ”UNILEVER” : 0,
        "MILLER COORS": 5300
	}
  ```


# HTTP routes for the server:
	POST "/" : Receives transaction call
	POST "/spend" : Receives spend request from user(you), responds with payer points spent.
	GET"/points" : blank request will respond with current payers points.
