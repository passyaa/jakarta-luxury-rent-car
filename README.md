## Expected Result

### URL Individual Project

1. Local Server : http://localhost:8080/
2. Heroku Server : 

### List API
1. Download Colletion Localhost Postman [Postman Collection]()
2. Download Colletion Heroku Postman [Postman Collection]()

| Method | Endpoint                                  | Description                                  |
|--------|-------------------------------------------|----------------------------------------------|
| POST   | `/register`                               | Register new account                         |
| POST   | `/login`                                  | Login and obtain JWT token                   |
| GET    | `/cars`                                   | Get data luxury cars                         |
| GET    | `/drivers`                                | Get data drivers                             |
| GET    | `/packages`                               | Get data event packages                      |
| POST   | `/users/register-membership`              | Register membership                          |
| GET    | `/users/get-membership`                   | Get data membership                          |
| POST   | `/users/topup`                            | Topup deposit amount                         |
| GET    | `/users/get-deposit`                      | Get data deposit amount                      |
| POST   | `/users/booking`                          | Booking luxury cars                          |
| POST   | `/users/making-payment`                   | Payment                                      |
| POST   | `/users/call-assistance`                  | Call Assistance if you get the trouble       |
| POST   | `/owner/approve-booking`                  | Approve booking from user                    |
| GET    | `/owner/report`                           | Get details report                           |

### Swaggo Doc
1. Access Swagger UI Localhost : Open your browser and navigate to http://localhost:8080/swagger/index.html
2. Access Swagger UI Localhost : Open your browser and navigate to ()
3. Authorize with JWT : When using the Authorize feature, ensure you manually input your token with the "Bearer" prefix. The token should be entered as "Bearer <your_jwt_token>". This is necessary because the "Bearer" prefix must be included manually


## Flow Process untuk Jasa Sewa Mobil Mewah
![Flow Process](https://github.com/passyaa/jakarta-luxury-rent-car/blob/master/assets/image/jakarta-luxury-rent-car.png)


## ERD
Open your browser (https://drawsql.app/teams/kamal-teams/diagrams/jakarta-luxury-rent-car)
![ERD IMAGE](https://github.com/passyaa/jakarta-luxury-rent-car/blob/master/assets/image/ERD.PNG)

### Entities and Relationships
- Setiap User hanya bisa memiliki satu jenis Membership dan satu Membership hanya bisa dimiliki oleh satu User (One-to-One)
- Setiap User dapat memiliki banyak RentalHistory, tapi setiap RentalHistory hanya bisa dikaitkan dengan satu User (One-to-Many)
- Setiap Car dapat disewa berkali-kali, tapi setiap RentalHistory hanya berkaitan dengan satu Car tertentu (One-to-Many)
- Setiap Driver bisa terkait dengan banyak RentalHistory, tapi setiap RentalHistory hanya bisa memiliki satu Driver tertentu (One-to-Many)
- Setiap EventPackage bisa digunakan dalam banyak RentalHistory, tapi setiap RentalHistory hanya bisa memiliki satu EventPackage tertentu (One-to-Many)
- Setiap RentalHistory bisa memiliki beberapa CallAssistance (One-to-Many)
- Setiap User bisa meminta bantuan CallAssistance, tapi setiap CallAssistance hanya terkait dengan satu User (One-to-Many)

## Deployment Notes

### Persiapan DB to SUPABASE
1. Connect to DB PostgreSQL - SUPABASE

### Preparation Environment Variables
1. Sesuaikan .env dengan DB PostgreSQL - SUPABASE
    ```sh
    DB_USER=postgres.csqufkmaujkysluqfqqt
    DB_PASSWORD=4YZzcVnWyRkRXMwfYUeq
    DB_HOST=aws-0-ap-southeast-1.pooler.supabase.com
    DB_PORT=6543
    DB_NAME=postgres
    JWT_SIGNING_KEY=GC2
   ```

### Tambah file DockerFile & heroku.yml
1. Buat file DockerFile
    ```sh
    FROM golang:1.22.1

    WORKDIR /app

    COPY go.mod ./
    COPY go.sum ./
    RUN go mod download

    COPY . ./

    RUN go build -o /main

    EXPOSE 8080

    CMD ["/main"]
    ```
2. Buat file heroku.yml
    ```sh
    build:
    docker:
        web: Dockerfile
    ```

### Inisialisasi Git Repository
1. cd (folder name)
2. git init
3. git add .
4. git commit -m "Deploy to Heroku"

### Buat Aplikasi di Heroku
1. heroku login
2. heroku create api-kamal --buildpack heroku/go
3. git branch // (optional) untuk melihat position file
4. git remote -v // (optional) untuk melihat alamat git
5. git remote add origin https://git.heroku.com/api-kamal-gc2.git // (optional) jika belum ada tambahkan alamat githubnya
6. git push heroku master

### Konfigurasi Environment Variables
1. heroku config:set DB_USER=
2. heroku config:set DB_PASSWORD=
3. heroku config:set DB_HOST=
4. heroku config:set DB_PORT=
5. heroku config:set DB_NAME=
6. heroku config:set JWT_SIGNING_KEY=
7. heroku config:set GO111MODULE=on
8. heroku config:set PORT=8080

### Verifikasi Deployment
1. heroku open

### Monitoring dan Logs
1. heroku logs --tail

### Jika ada Perubahan
1. Setelah melakukan perubahan, commit perubahan dan deploy ulang aplikasi ke Heroku:    
    ```sh
    git add .
    git commit -m "changes messages"
    git branch // (optional) untuk melihat position file
    git remote -v // (optional) untuk melihat alamat git
    git remote add origin https://git.heroku.com/api-kamal-gc2.git // (optional) jika belum ada tambahkan alamat githubnya
    git push heroku master
    heroku restart
   ```


