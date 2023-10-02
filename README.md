# Turso Connection with Go App.

A Simple REST API to perform CRUD Operations on a collection of albums.

Can use Turso DB or local sqlite files for the database. Great for applications which have a low data footprint.

API is Blazingly Fast 🚀 when used with local sqlite3 file.
Below is a result of a stress test with 1000000 requests from 125 concurrent connections.
<img width="1431" alt="Screenshot 2023-10-02 at 4 19 24 PM" src="https://github.com/tapasrm/go-turso-rest-api/assets/24273309/01505c78-d852-4ece-8569-f8a60e76e707">
My M1 Macbook Air was able to handle 30K Requests/Second on average and return the data successfully for all 1M Requests.
