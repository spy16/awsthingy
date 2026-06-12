package main

const defaultLineTpl = `{{ .InstanceId }} :: {{ tag . "Name" }} {{.PrivateDnsName}}`
const defaultOutTpl = `{{ .InstanceId }}`

const defaultPreviewTpl = `{{.InstanceId}} :: {{ name .}}
---------------------------------------------------
Status: {{.State.Name}}
Architecture: {{.Architecture}}
AMI ID: {{.ImageId}}
Private IP:
IP: 🔐 {{.PrivateIpAddress}} / 🌍 {{.PublicIpAddress}}

Tags:{{range .Tags}}
* {{.Key}} = {{.Value}}{{end}}
`
