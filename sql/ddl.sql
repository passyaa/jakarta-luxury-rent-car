CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    address TEXT NOT NULL,
    deposit_amount NUMERIC(10,2) DEFAULT 0,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
);

CREATE TABLE Cars (
    car_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock_availability INT NOT NULL CHECK (stock_availability >= 0),
    rental_costs NUMERIC(10,2) NOT NULL,
    category VARCHAR(255) NOT NULL,
    make VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    transmission VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    fuel_type VARCHAR(255) NOT NULL,
    class VARCHAR(255) NOT NULL
);

CREATE TABLE Drivers (
    driver_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    license_number VARCHAR(255) UNIQUE NOT NULL,
    experience_years INT NOT NULL,
    rating NUMERIC(3,2) DEFAULT 5.0
);

CREATE TABLE EventPackages (
    package_id SERIAL PRIMARY KEY,
    package_name VARCHAR(255) NOT NULL,
    description TEXT,
    cost NUMERIC(10,2) NOT NULL
);

CREATE TABLE Memberships (
    membership_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    discount_level VARCHAR(10) NOT NULL CHECK (discount_level IN ('Silver', 'Gold', 'Platinum')),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE RentalHistory (
    rental_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    car_id INT NOT NULL,
    driver_id INT,
    rental_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP,
    total_cost NUMERIC(10,2) NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN ('Book', 'Paid', 'Rent', 'Cancel')),
    package_id INT,
    airport_transfer BOOLEAN DEFAULT FALSE,
    pickup_location TEXT,
    dropoff_location TEXT,
    concierge_services BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (car_id) REFERENCES Cars(car_id),
    FOREIGN KEY (driver_id) REFERENCES Drivers(driver_id),
    FOREIGN KEY (package_id) REFERENCES EventPackages(package_id)
);

CREATE TABLE CallAssistance (
    assistance_id SERIAL PRIMARY KEY,
    rental_id INT NOT NULL,
    user_id INT NOT NULL,
    callassistance_date TIMESTAMP NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    FOREIGN KEY (rental_id) REFERENCES RentalHistory(rental_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);