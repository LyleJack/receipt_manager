## To Do:
- Look into setting up OAuth properly and then incorperating that into the user.
- Change DB from SQLite -> Postgres
- Maybe an S3 AWS bucket for receipt images?
- See if it is worthwhile making the AI response for receipt -> JSON into a Go routine - maybe the whole process,
 so that it can handle multiple receipts?
    - however the OCR is quite fast.
- Some form of check to ensure the receipt looks correct (check for certain fields at minimum - possibly another call to genai?)

## To Run:  
> `docker-compose up --build`  
`Ctrl+C` can be used to stop the container.  

## To cURL:  
> curl http://localhost:8000 

### Test example:
> curl http://localhost:8000/test

Should return with the response:
```json
[
    {
        "date":"2024-04-23T19:32:00Z",
        "items":[
            {
                "name":"Lager",
                "quantity":8,
                "total_price":15.75
            },
            {
                "name":"Steak & Ale Pie",
                "quantity":1,
                "total_price":13.50
            },
            {
                "name":"Fish and Chips",
                "quantity":1,
                "total_price":14.95
            },
            {
                "name":"Mixed Grill",
                "quantity":1,
                "total_price":18.25
            },
            {
                "name":"Red Wine",
                "quantity":2,
                "total_price":14.00
            },
            {
                "name":"Dessert",
                "quantity":1,
                "total_price":7.95
            }
        ],
        "store_name":"TasteRadar",
        "total":168.36
    }
]
```

# Issues
