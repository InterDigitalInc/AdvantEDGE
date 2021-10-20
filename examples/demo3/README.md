# Demo3

## Quick start 

|Operations: | Notes:|
|------------|-------|
|1. Login via the frontend & select Sandbox tab|
|2. Select a network to deploy in the user sandbox | Choose dual-mep network scenario to utilize app mobility service use-case |
|3. Create a unique MEC Application Instance ID | If using dual-mep network and walkthrough app mobility service use-case then both MEC Application Instance IDs must be created with the same MEC Application name|
|4. Apply configuration values by copying the MEC Sandbox endpoint into app_instance1.yaml sandbox mecUrl value  | For example: <br> Endpoints `mep1`:  `https://<the-mec-url>/<the-sandbox-key>/mep1` <br>Endpoints `mep2`: `https://<the-mec-url>/<the-sandbox-key>/mep2` |
|5. Select & copy the MEC Application Instance Id into app_instance1.yaml appid | |
|6. Select & copy your running app ip address into app_instace1.yaml localurl & port | |
|7. Run main.go pass in arguments with path of your config | |
|8. You can choose to launch swagger as client to intereact with api's and find responses in the logs demo 3 app | 

  
## What's inside?

A quick look at the top-level relevant files and directories in demo 3 project.

    .
    ├── api
    ├── server
    ├── util
    ├── main.go
    ├── go.mod
    

- api (swagger documentation)
- server (server code)
- util (configurations)
- main (entry)
