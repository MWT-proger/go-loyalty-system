# Настройки для проекта в VS Code

## launch.json
```
{
    "version": "0.2.0",
    "configurations": [
      {
        "name": "Launch Go",
        "type": "go",
        "request": "launch",
        "mode": "debug",
        "program": "${workspaceFolder}/cmd/gophermart",
        "env": {
          "RUN_ADDRESS": ":7000",
          "DATABASE_URI": "user=postgres password=postgres host=localhost port=5432 dbname=testDB sslmode=disable",
          "ACCRUAL_SYSTEM_ADDRESS": "http://localhost:6000"
        },
        "args": ["-a=${env.RUN_ADDRESS}", "-d=${env.DATABASE_URI}", "-r=${env.ACCRUAL_SYSTEM_ADDRESS}", "-l=debug"]
      }
    ]
  }
```