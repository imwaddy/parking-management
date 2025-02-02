# parking-spot

#### **1. Create Parking Spot**
* **Create Parking Spot Endpoint:** (/api/parking-spots,"POST")
* **Endpoint/Route:** localhost:8080/api/parking-spots (POST)
* **Method:** POST request
* **Request Body:**
```
{
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "yes"
}
```
* **Response:**
```
{
  "id": 10,
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "yes"
}
```

#### **2. Get All Parking Spot**
* **Get All Parking Spot Endpoint:** (/api/parking-spots/all,"GET")
* **Endpoint/Route:** localhost:8080/api/parking-spots/all (GET)
* **Method:** GET request
* **Response:**
```
[{
  "id": 10,
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "yes"
},{
  "id": 11,
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "yes"
}]
```

#### **3. Get Parking Spot By ID**
* **Get Parking Spot By ID Endpoint:** (/api/parking-spots/{id},"GET")
* **Endpoint/Route:** localhost:8080/api/parking-spots/{id} (GET)
* **Method:** GET request
* **Response:**
```
{
  "id": 10,
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "yes"
}
```

#### **4. Update Parking Spot**
* **Update Parking Spot Endpoint:** (/api/parking-spots/{id},"PUT")
* **Endpoint/Route:** localhost:8080/api/parking-spots/{id} (PUT)
* **Method:** PUT request
* **Request Body:**
```
{
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "no"
}
```
* **Response:**
```
{
  "id": 10,
  "spot_number": "A1",
  "type": "Compact",
  "is_available": "yes"
}
```

#### **4. Delete Parking Spot**
* **Delete Parking Spot Endpoint:** (/api/parking-spots/{id},"DELETE")
* **Endpoint/Route:** localhost:8080/api/parking-spots/{id} (DELETE)
* **Method:** DELETE request
* **Request Body:**
* **Response #1: Successful deletion**
```
{
  "message": "Parking spot has been deleted successfully."
}
```
* **Response #2: Validation if parking spot not provided / not found**
```
Parking spot not found - With status 404
```