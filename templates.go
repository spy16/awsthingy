package main

const defaultLineTpl = `{{ .InstanceId }} :: {{ tag . "Name" }}`
const defaultOutTpl = `{{ .InstanceId }}`

const defaultPreviewTpl = `{{.InstanceId}} :: {{ name .}}
---------------------------------------------------
Operating System: {{.Platform}}
Architecture: {{.Architecture}}
AMI ID: {{.ImageId}}
Private IP: 
IP: 🔐 {{.PrivateIpAddress}} / 🌍 {{.PublicIpAddress}}

Tags:{{range .Tags}}
* {{.Key}} = {{.Value}}{{end}}
`
