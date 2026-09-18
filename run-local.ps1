# run-local.ps1
# Sobe o backend no ambiente local.
#
# O servidor le a conexao do banco de variaveis de ambiente (ver
# internal/database/connection.go). Sem elas ele morre logo no comeco com
# "strconv.Atoi: parsing \"\": invalid syntax", que e o SQL_PORT vazio.
#
# Uso:  .\run-local.ps1

$env:SQL_HOST = "localhost"
$env:SQL_PORT = "5432"
$env:SQL_USER = "postgres"
$env:SQL_DB   = "epi"

# A senha nao fica no arquivo. Ou voce exporta antes
# ($env:SQL_PASSWORD = "..."), ou o script pergunta.
if (-not $env:SQL_PASSWORD) {
    $secure = Read-Host "Senha do postgres" -AsSecureString
    $env:SQL_PASSWORD = [System.Net.NetworkCredential]::new("", $secure).Password
}

# Hardware do armario. Sem o leitor e sem a placa de porta, as telas de cadastro
# e o aviso de retirada funcionam normal; o que nao anda e abrir a porta e
# contar as tags.
if (-not $env:RFID_CONTROLLER) { $env:RFID_CONTROLLER = "" }
if (-not $env:DOOR_PORTS)      { $env:DOOR_PORTS      = "" }
if (-not $env:DOOR_STATES)     { $env:DOOR_STATES     = "" }
if (-not $env:SENSOR_PORTS)    { $env:SENSOR_PORTS    = "" }
if (-not $env:SENSOR_STATES)   { $env:SENSOR_STATES   = "" }

Write-Host "Subindo o backend em http://localhost:8080 ..." -ForegroundColor Cyan
.\epi-backend.exe
