{{/* Stable labels; environments are isolated by namespace. */}}
{{- define "0things.labels" -}}
app.kubernetes.io/name: 0things
app.kubernetes.io/component: {{ . }}
app.kubernetes.io/managed-by: Helm
{{- end -}}

{{/* Render image pull secrets only when configured. */}}
{{- define "0things.imagePullSecrets" -}}
{{- if or .Values.registry.create .Values.imagePullSecrets }}
imagePullSecrets:
  {{- if .Values.registry.create }}
  - name: aiot-platform-registry-auth
  {{- end }}
  {{- range .Values.imagePullSecrets }}
  - name: {{ . }}
  {{- end }}
{{- end }}
{{- end -}}

{{/* Optional storage class stanza for volumeClaimTemplates. */}}
{{- define "0things.storageClass" -}}
{{- if .storageClass }}
storageClassName: {{ .storageClass | quote }}
{{- end }}
{{- end -}}
