## To Do:
- Look into setting up OAuth properly and then incorperating that into the user.
- Store JSON in database
- Maybe an S3 AWS bucket for receipt images?
- See if it is worthwhile making the AI response for receipt -> JSON into a Go routine - maybe the whole process, so that it can handle multiple receipts?
    - however the OCR is quite fast.

## To Run:  
> `docker-compose up --build`  
`Ctrl+C` can be used to stop the container.  

## To cURL:  
> curl http://localhost:8000 

### Test example:
> curl http://localhost:8000/test

Should return with the response:
```json
{
    "date":"04/23/2024",
    "items":[
        {"name":"Lager",
        "price":15.75},
        {"name":"Steak & Ale Pie",
        "price":13.50},
        {"name":"Fish and Chips",
        "price":14.95},
        {"name":"Mixed Grill",
        "price":18.25},
        {"name":"Red Wine",
        "price":14.00},
        {"name":"Dessert",
        "price":7.95}
    ],
    "store_name":"TasteRadar",
    "total":168.36 
}
```

# Issues
