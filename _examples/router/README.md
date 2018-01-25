# oas3 router example

Run the server:

    go run main.go -spec ../petstore.yaml
    
Send requests:

    curl -XPOST -i localhost:3000/v2/pet
    
    curl -XGET  -i localhost:3000/v2/store/inventory
    
Now try a request that does not meet spec parameters requirements:
    
    curl -XGET -i "localhost:3000/v2/user/login?username=johndoe"
    
    # HTTP/1.1 400 Bad Request
    # {"errors":[{"message":"param password is required","field":"password"}]}