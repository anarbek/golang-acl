@echo off
powershell -Command "& {for ($i=0; $i -lt 10; $i++) {Invoke-WebRequest -Uri http://127.0.0.1:56990/api/v1/users -Headers @{'Authorization' = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1MTQ4MDgwMDAsInVzZXIiOiJKb2huIERvZSIsInVzZXJJZCI6MX0.lAvKklgDjxQOehqbsqBd3btgsFbSQXBS-GShLdmpEYg'}}}"
