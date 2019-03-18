$url="https://pcm.k636174.net/api/heartbeat"

[String] $MyHostname = [Net.Dns]::GetHostName()
[String] $MyIPs = [system.net.dns]::GetHostAddresses((hostname)) | where {$_.AddressFamily -eq "InterNetwork"} | select -ExpandProperty IPAddressToString

$postParams = @{hostname=$MyHostname; src_lip=$MyIPs}

Invoke-WebRequest $url -Method POST -Body $postParams