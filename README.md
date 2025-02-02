# parking-management

In this assignment, participants should work on two separate projects:
1. **Vehicle Entry and Exit:** This project handles the registration of vehicle entries and exits, including the management of
vehicle records such as spot number, license plate, entry time, and exit time.
2. **Parking Lot Management:** This project focuses on managing parking spots, where participants will implement CRUD
operations for parking spots (add, retrieve, update, and delete) and handle the availability and types of parking spots.

The **InitDB** function provided in the stub code connects to a PostgreSQL database using the given credentials and automatically creates two tables: **vehicle_records** and **parking_spots**. The **vehicle_records** table stores information such as **spot_number, license_plate, entry_time, and exit_time,** while the **parking_spots** table includes **spot_number, type, and is_available.** When participants run the application, the database will be initialized, and the tables will be created if they do not already exist.
