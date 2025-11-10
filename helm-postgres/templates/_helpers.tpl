{{- define "db.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "postgresql.host" -}}
{{- printf "%s-postgresql" (include "db.name" .) }}
{{- end }}