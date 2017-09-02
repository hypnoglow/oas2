# oas3 router example

Run the server:

    go run main.go -spec ../petstore.yaml
    
Send requests:

    curl -XPOST localhost:3000/v2/pet
    
    curl -XGET  localhost:3000/store/inventory