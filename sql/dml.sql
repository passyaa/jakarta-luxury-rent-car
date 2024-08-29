-- Insert into Drivers
INSERT INTO drivers (name, phone_number, license_number, experience_years, rating) VALUES
('Budi Santoso', '081234567890', 'SIM001', 5, 4.8),
('Andi Wijaya', '082345678901', 'SIM002', 10, 4.9);

-- Insert into Event Packages
INSERT INTO event_packages (package_name, description, cost) VALUES
('Wedding Package', 'Wedding Package', 500.00),
('Corporate Event Package', 'Corporate Event Package', 1000.00);

-- Insert into Users
INSERT INTO Users (email, password, phone_number, address, deposit_amount, role) VALUES
('user1@example.com', 'user1', '+6281211115030', 'Jl. Kebon Jeruk No. 12, Jakarta', 100.00, 'user'),
('user2@example.com', 'user2', '+6285894999562', 'Jl. Kebon Jeruk No. 12, Jakarta', 50.00, 'admin');

-- Insert into Rental History
INSERT INTO rental_histories (user_id, car_id, driver_id, rental_date, return_date, total_cost, status, package_id, airport_transfer, pickup_location, dropoff_location, concierge_services) VALUES
(1, 1, 1, '2024-08-01', '2024-08-05', 500.00, 'Completed', 1, TRUE, 'Jakarta Airport', 'Hotel Indonesia Kempinski', FALSE),
(2, 2, 2, '2024-08-10', NULL, 150.00, 'Ongoing', NULL, FALSE, 'Grand Hyatt Jakarta', 'Plaza Indonesia', TRUE);

-- Insert into Roadside Assistance
INSERT INTO CallAssistance (rental_id, user_id, callassistance_date, description, location) VALUES
(1, 1, '2024-08-03 12:00:00', 'Flat tire assistance needed', 'Main Road, Jakarta'),
(2, 2, '2024-08-12 15:30:00', 'Engine trouble', 'Downtown, Jakarta');

-- Insert into Memberships
INSERT INTO Memberships (user_id, discount_level) VALUES
(1, 'Silver'),
(2, 'Gold');

-- Insert into Cars
INSERT INTO cars (name, stock_availability, rental_costs, category, make, model, transmission, year, fuel_type, class) VALUES
('Toyota Camry', 5, 100.00, 'Sedan', 'Toyota', 'Camry', 'Automatic', 1993, 'Gas', 'Midsize car'),
('Honda Civic', 10, 80.00, 'Sedan', 'Honda', 'Civic', 'Manual', 2005, 'Gas', 'Compact car'),
('Ford Mustang', 2, 150.00, 'Sports', 'Ford', 'Mustang', 'Automatic', 2018, 'Gas', 'Sports car'),
('Chevrolet Impala', 7, 95.50, 'Sedan', 'Chevrolet', 'Impala', 'Automatic', 2016, 'Gas', 'Full-size car'),
('BMW X5', 3, 200.00, 'SUV', 'BMW', 'X5', 'Automatic', 2020, 'Diesel', 'Luxury SUV'),
('Audi A4', 4, 110.00, 'Sedan', 'Audi', 'A4', 'Automatic', 2019, 'Gas', 'Compact executive car'),
('Mercedes-Benz C-Class', 5, 180.00, 'Sedan', 'Mercedes-Benz', 'C-Class', 'Automatic', 2021, 'Gas', 'Luxury car'),
('Nissan Altima', 8, 85.00, 'Sedan', 'Nissan', 'Altima', 'Manual', 2010, 'Gas', 'Midsize car'),
('Hyundai Sonata', 9, 90.00, 'Sedan', 'Hyundai', 'Sonata', 'Automatic', 2015, 'Gas', 'Midsize car'),
('Kia Optima', 6, 75.00, 'Sedan', 'Kia', 'Optima', 'Manual', 2017, 'Gas', 'Midsize car'),
('Mazda CX-5', 3, 130.00, 'SUV', 'Mazda', 'CX-5', 'Automatic', 2019, 'Gas', 'Compact SUV'),
('Subaru Outback', 7, 120.00, 'SUV', 'Subaru', 'Outback', 'Automatic', 2022, 'Gas', 'Crossover SUV'),
('Volkswagen Jetta', 6, 85.00, 'Sedan', 'Volkswagen', 'Jetta', 'Manual', 2018, 'Gas', 'Compact car'),
('Volvo XC90', 2, 210.00, 'SUV', 'Volvo', 'XC90', 'Automatic', 2020, 'Hybrid', 'Luxury SUV'),
('Tesla Model S', 4, 250.00, 'Sedan', 'Tesla', 'Model S', 'Automatic', 2023, 'Electric', 'Luxury electric car'),
('Ford Explorer', 5, 140.00, 'SUV', 'Ford', 'Explorer', 'Automatic', 2017, 'Gas', 'Mid-size SUV'),
('Chevrolet Tahoe', 3, 160.00, 'SUV', 'Chevrolet', 'Tahoe', 'Automatic', 2019, 'Gas', 'Full-size SUV'),
('Toyota RAV4', 10, 110.00, 'SUV', 'Toyota', 'RAV4', 'Automatic', 2016, 'Gas', 'Compact SUV'),
('Honda Accord', 8, 95.00, 'Sedan', 'Honda', 'Accord', 'Automatic', 2014, 'Gas', 'Midsize car'),
('Jeep Wrangler', 4, 180.00, 'SUV', 'Jeep', 'Wrangler', 'Manual', 2021, 'Gas', 'Off-road SUV');