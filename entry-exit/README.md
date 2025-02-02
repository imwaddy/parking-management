# entry-exit


#### **1. Create Vehicle Entries**
* **Create Vehicle Entries:** (/api/vehicle-entries,"POST")
* **Endpoint/Route:** localhost:8081/api/vehicle-entries (POST)
* **Method:** POST request
* **Request Body:**
Before entering, ensure that the **spot_number** is available, i.e., the **is_available** field in the **parking_spots** table is set to **"yes"**.
```
{
  "spot_number": "A1",
  "license_type": "KA03122"
}
```
* **Response #1: Successful entry**
```
{
  "id": 10,
  "spot_number": "A1",
  "license_type": "KA03122",
  "entry_time": "31-12-2004 12:22:34"
}
```
* **Response #2: When parking spot not available**
Before entering, ensure that the **spot_number** is not available, i.e., the **is_available** field in the **parking_spots** table is set to **"no"**.
```
Parking spot not found - With status 404
```


#### **2. Create Vehicle Exits**
* **Create Vehicle Exits:** (/api/vehicle-entries,"POST")
* **Endpoint/Route:** localhost:8081/api/vehicle-entries (POST)
* **Method:** POST request
* **Request Body:**
Before exiting, ensure that the vehicle's entry has been recorded.
```
{
  "spot_number": "A1",
  "license_type": "KA03122"
}
```
* **Response**
To format the date and time into the format "dd-mm-yyyy hh:mm:ss" in Go, you can use the time.Format method with the
format string **"02-01-2006 15:04:05"**
```
{
  "id": 10,
  "spot_number": "A1",
  "license_type": "KA03122",
  "entry_time": "31-12-2004 12:22:34"
  "exit_time": "31-12-2004 12:40:34"
}
```


#### **3. Create Vehicle Entries when the spot number is not available or registered**
* **Create Vehicle Entries:** (/api/vehicle-entries/{spot_number},"GET")
* **Endpoint/Route:** localhost:8081/api/vehicle-entries/{spot_number} (GET)
* **Method:** GET request

* **Response #1: Found**
```
[{
  "id": 10,
  "spot_number": "A1",
  "license_type": "KA03122",
  "entry_time": "31-12-2004 12:22:34"
  "exit_time": "31-12-2004 12:40:34"
}]
```

* **Response #2: Error / not found**
```
Parking spot not found - With status 404
```